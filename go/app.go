package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	redis "gopkg.in/redis.v3"
)

var (
	db    *sql.DB
	store *sessions.CookieStore
	rd    *redis.Client
)

type User struct {
	ID          int
	AccountName string
	NickName    string
	Email       string
}

type Profile struct {
	UserID    int
	FirstName string
	LastName  string
	Sex       string
	Birthday  mysql.NullTime
	Pref      string
	UpdatedAt time.Time
}

type Entry struct {
	ID        int
	UserID    int
	Private   bool
	Title     string
	Content   string
	CreatedAt time.Time
}

type Comment struct {
	ID        int
	EntryID   int
	UserID    int
	Comment   string
	CreatedAt time.Time
}

type CommentWithEntry struct {
	ID        int
	EntryID   int
	UserID    int
	Comment   string
	CreatedAt time.Time
	Entry     Entry
}

type Friend struct {
	ID        int
	CreatedAt time.Time
}

type Footprint struct {
	UserID    int
	OwnerID   int
	CreatedAt time.Time
	Updated   time.Time
}

var prefs = []string{"未入力",
	"北海道", "青森県", "岩手県", "宮城県", "秋田県", "山形県", "福島県", "茨城県", "栃木県", "群馬県", "埼玉県", "千葉県", "東京都", "神奈川県", "新潟県", "富山県",
	"石川県", "福井県", "山梨県", "長野県", "岐阜県", "静岡県", "愛知県", "三重県", "滋賀県", "京都府", "大阪府", "兵庫県", "奈良県", "和歌山県", "鳥取県", "島根県",
	"岡山県", "広島県", "山口県", "徳島県", "香川県", "愛媛県", "高知県", "福岡県", "佐賀県", "長崎県", "熊本県", "大分県", "宮崎県", "鹿児島県", "沖縄県"}

var (
	ErrAuthentication   = errors.New("Authentication error.")
	ErrPermissionDenied = errors.New("Permission denied.")
	ErrContentNotFound  = errors.New("Content not found.")
)

var tport = flag.Uint("port", 0, "port to listen")
var mysqlsock = flag.Bool("mysqlsock", true, "connect to mysql via unix domain socket")

func authenticate(w http.ResponseWriter, r *http.Request, email, passwd string) {
	query := `SELECT u.id AS id, u.account_name AS account_name, u.nick_name AS nick_name, u.email AS email
FROM users u
JOIN salts s ON u.id = s.user_id
WHERE u.email = ? AND u.passhash = SHA2(CONCAT(?, s.salt), 512)`
	row := db.QueryRow(query, email, passwd)
	user := User{}
	err := row.Scan(&user.ID, &user.AccountName, &user.NickName, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			checkErr(ErrAuthentication)
		}
		checkErr(err)
	}
	session := getSession(w, r)
	session.Values["user_id"] = user.ID
	session.Save(r, w)
}

func getCurrentUser(w http.ResponseWriter, r *http.Request) *User {
	u := context.Get(r, "user")
	if u != nil {
		user := u.(User)
		return &user
	}
	session := getSession(w, r)
	userID, ok := session.Values["user_id"]
	if !ok || userID == nil {
		return nil
	}
	userIDInt, ok := userID.(int)
	if !ok {
		return nil
	}
	return getUser(w, userIDInt)
}

func authenticated(w http.ResponseWriter, r *http.Request) bool {
	user := getCurrentUser(w, r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return false
	}
	return true
}

var allUsers = make(map[int]User, 50000)

func getUser(w http.ResponseWriter, userID int) *User {
	user, ok := allUsers[userID]
	if !ok {
		checkErr(ErrContentNotFound)
	}
	return &user
}

func initAllUsers() {
	allUsers = make(map[int]User, 50000)
	fmt.Println("start initAllUsers")
	rows, err := db.Query(`SELECT * FROM users`)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.AccountName, &user.NickName, &user.Email, new(string))
		if err != nil {
			panic(err)
		}
		allUsers[user.ID] = user
	}
	rows.Close()
	fmt.Println("finish initAllUsers")
}

func getUserFromAccount(w http.ResponseWriter, name string) *User {
	row := db.QueryRow(`SELECT * FROM users WHERE account_name = ?`, name)
	user := User{}
	err := row.Scan(&user.ID, &user.AccountName, &user.NickName, &user.Email, new(string))
	if err == sql.ErrNoRows {
		checkErr(ErrContentNotFound)
	}
	checkErr(err)
	return &user
}

func isFriend(w http.ResponseWriter, r *http.Request, anotherID int) bool {
	session := getSession(w, r)
	id := session.Values["user_id"]
	row := db.QueryRow(`SELECT COUNT(1) AS cnt FROM relations WHERE (one = ? AND another = ?) OR (one = ? AND another = ?)`, id, anotherID, anotherID, id)
	cnt := new(int)
	err := row.Scan(cnt)
	checkErr(err)
	return *cnt > 0
}

func isFriendAccount(w http.ResponseWriter, r *http.Request, name string) bool {
	user := getUserFromAccount(w, name)
	if user == nil {
		return false
	}
	return isFriend(w, r, user.ID)
}

func permitted(w http.ResponseWriter, r *http.Request, anotherID int) bool {
	user := getCurrentUser(w, r)
	if anotherID == user.ID {
		return true
	}
	return isFriend(w, r, anotherID)
}

func footprintKey(userID int) string {
	return fmt.Sprintf("footprint-%d", userID)
}

func markFootprint(w http.ResponseWriter, r *http.Request, id int) {
	user := getCurrentUser(w, r)
	if user.ID != id {
		rd.ZAdd(footprintKey(id /*見られてる人*/), redis.Z{float64(time.Now().Unix()), user.ID /*見てる人*/})
		//_, err := db.Exec(`INSERT INTO footprints (user_id,owner_id) VALUES (?,?)`, id, user.ID)
		//checkErr(err)
	}
}

func getFootprintsFromRedis(userID int, num int64) ([]Footprint, error) {
	vals, err := rd.ZRevRangeWithScores(footprintKey(userID), 0, num).Result()
	if err != nil {
		return nil, err
	}
	footprints := make([]Footprint, 0, num)
	for _, val := range vals {
		t := time.Unix(int64(val.Score), 0)
		viewerIDStr, _ := val.Member.(string)
		viewerID, _ := strconv.Atoi(viewerIDStr)
		footprints = append(footprints, Footprint{userID, viewerID, t, t})
	}
	return footprints, nil
}

func myHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rcv := recover()
			if rcv != nil {
				switch {
				case rcv == ErrAuthentication:
					session := getSession(w, r)
					delete(session.Values, "user_id")
					session.Save(r, w)
					render(w, r, http.StatusUnauthorized, "login.html", struct{ Message string }{"ログインに失敗しました"})
					return
				case rcv == ErrPermissionDenied:
					render(w, r, http.StatusForbidden, "error.html", struct{ Message string }{"友人のみしかアクセスできません"})
					return
				case rcv == ErrContentNotFound:
					render(w, r, http.StatusNotFound, "error.html", struct{ Message string }{"要求されたコンテンツは存在しません"})
					return
				default:
					var msg string
					if e, ok := rcv.(runtime.Error); ok {
						msg = e.Error()
					}
					if s, ok := rcv.(string); ok {
						msg = s
					}
					msg = rcv.(error).Error()
					http.Error(w, msg, http.StatusInternalServerError)
				}
			}
		}()
		fn(w, r)
	}
}

func getSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "isucon5q-go.session")
	return session
}

func getTemplatePath(file string) string {
	return path.Join("templates", file)
}

func render(w http.ResponseWriter, r *http.Request, status int, file string, data interface{}) {
	fmap := template.FuncMap{
		"getUser": func(id int) *User {
			return getUser(w, id)
		},
		"getCurrentUser": func() *User {
			return getCurrentUser(w, r)
		},
		"isFriend": func(id int) bool {
			return isFriend(w, r, id)
		},
		"prefectures": func() []string {
			return prefs
		},
		"substring": func(s string, l int) string {
			if len(s) > l {
				return s[:l]
			}
			return s
		},
		"split": strings.Split,
		"getEntry": func(id int) Entry {
			row := db.QueryRow(`SELECT entries.id, entries.user_id, private, SUBSTRING_INDEX(body, '\n', 1) as subject, created_at FROM entries WHERE id=?`, id)
			var entryID, userID, private int
			var subject string
			var createdAt time.Time
			checkErr(row.Scan(&entryID, &userID, &private, &subject, &createdAt))
			return Entry{id, userID, private == 1, subject, "", createdAt}
		},
		"numComments": func(id int) int {
			row := db.QueryRow(`SELECT COUNT(*) AS c FROM comments WHERE entry_id = ?`, id)
			var n int
			checkErr(row.Scan(&n))
			return n
		},
	}
	tpl := template.Must(template.New(file).Funcs(fmap).ParseFiles(getTemplatePath(file), getTemplatePath("header.html")))
	w.WriteHeader(status)
	checkErr(tpl.Execute(w, data))
}

func GetLogin(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusOK, "login.html", struct{ Message string }{"高負荷に耐えられるSNSコミュニティサイトへようこそ!"})
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	passwd := r.FormValue("password")
	authenticate(w, r, email, passwd)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetLogout(w http.ResponseWriter, r *http.Request) {
	session := getSession(w, r)
	delete(session.Values, "user_id")
	session.Options = &sessions.Options{MaxAge: -1}
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func isFriend2(friendsMap map[int]time.Time, userID int) bool {
	_, ok := friendsMap[userID]
	return ok
}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	if !authenticated(w, r) {
		return
	}

	user := getCurrentUser(w, r)

	prof := Profile{}
	row := db.QueryRow(`SELECT * FROM profiles WHERE user_id = ?`, user.ID)
	err := row.Scan(&prof.UserID, &prof.FirstName, &prof.LastName, &prof.Sex, &prof.Birthday, &prof.Pref, &prof.UpdatedAt)
	if err != sql.ErrNoRows {
		checkErr(err)
	}

	rows, err := db.Query(`SELECT * FROM entries WHERE user_id = ? ORDER BY created_at LIMIT 5`, user.ID)
	if err != sql.ErrNoRows {
		checkErr(err)
	}
	entries := make([]Entry, 0, 5)
	for rows.Next() {
		var id, userID, private int
		var body string
		var createdAt time.Time
		checkErr(rows.Scan(&id, &userID, &private, &body, &createdAt))
		entries = append(entries, Entry{id, userID, private == 1, strings.SplitN(body, "\n", 2)[0], strings.SplitN(body, "\n", 2)[1], createdAt})
	}
	rows.Close()

	commentsForMe := getCommentsForMe(user.ID, 10)

	// フレンド一覧を取得
	rows, err = db.Query(`SELECT * FROM relations WHERE one = ? OR another = ? ORDER BY created_at DESC`, user.ID, user.ID)
	if err != sql.ErrNoRows {
		checkErr(err)
	}
	friendsMap := make(map[int]time.Time)
	friendIds := make([]string, 0)
	for rows.Next() {
		var id, one, another int
		var createdAt time.Time
		checkErr(rows.Scan(&id, &one, &another, &createdAt))
		var friendID int
		if one == user.ID {
			friendID = another
		} else {
			friendID = one
		}
		if _, ok := friendsMap[friendID]; !ok {
			friendsMap[friendID] = createdAt
			friendIds = append(friendIds, strconv.Itoa(friendID))
		}
	}
	friends := make([]Friend, 0, len(friendsMap))
	for key, val := range friendsMap {
		friends = append(friends, Friend{key, val})
	}
	rows.Close()

	entryIndex := ""
	if len(friendIds) > 20 {
		entryIndex = "created_at"
	} else {
		entryIndex = "user_id"
	}
	// fmt.Println(entryIndex)
	rows, err = db.Query(fmt.Sprintf(`SELECT id, user_id, private, SUBSTRING_INDEX(body, '\n', 1) as subject, created_at FROM entries FORCE INDEX(%s) WHERE user_id IN (%s) ORDER BY created_at DESC LIMIT 10`, entryIndex, strings.Join(friendIds, ",")))
	if err != sql.ErrNoRows {
		checkErr(err)
	}
	entriesOfFriends := make([]Entry, 0, 10)
	for rows.Next() {
		var id, userID, private int
		var subject string
		var createdAt time.Time
		checkErr(rows.Scan(&id, &userID, &private, &subject, &createdAt))
		if !isFriend2(friendsMap, userID) {
			continue
		}
		entriesOfFriends = append(entriesOfFriends, Entry{id, userID, private == 1, subject, "", createdAt})
		if len(entriesOfFriends) >= 10 {
			break
		}
	}
	rows.Close()

	// コメントを1000件取得し、privateであれば自分と相手がフレンドかどうかをチェックする
	commentIndex := ""
	if len(friendIds) > 20 {
		commentIndex = "created_at"
	} else {
		commentIndex = "user_id"
	}
	query := fmt.Sprintf(`SELECT c.id, c.entry_id, c.user_id, SUBSTRING_INDEX(comment, '\n', 1) as comment, c.created_at, e.id, e.user_id, e.private, SUBSTRING_INDEX(e.body, '\n', 1), e.created_at FROM comments c FORCE INDEX (%s) JOIN entries e ON c.entry_id = e.id WHERE c.user_id IN (%s) ORDER BY c.created_at DESC LIMIT 10`, commentIndex, strings.Join(friendIds, ","))
	// fmt.Println(query)
	rows, err = db.Query(query)
	if err != sql.ErrNoRows {
		checkErr(err)
	}
	commentsOfFriends := make([]CommentWithEntry, 0, 10)
	for rows.Next() {
		c := CommentWithEntry{}
		var id, userID, private int
		var subject string
		var createdAt time.Time
		checkErr(rows.Scan(&c.ID, &c.EntryID, &c.UserID, &c.Comment, &c.CreatedAt, &id, &userID, &private, &subject, &createdAt))
		if !isFriend2(friendsMap, c.UserID) {
			continue
		}
		entry := Entry{id, userID, private == 1, subject, "", createdAt}
		c.Entry = entry
		if entry.Private {
			if !permitted(w, r, entry.UserID) {
				continue
			}
		}
		commentsOfFriends = append(commentsOfFriends, c)
		if len(commentsOfFriends) >= 10 {
			break
		}
	}
	rows.Close()

	footprints, err := getFootprintsFromRedis(user.ID, 10)
	if err != nil {
		checkErr(err)
	}

	render(w, r, http.StatusOK, "index.html", struct {
		User              User
		Profile           Profile
		Entries           []Entry
		CommentsForMe     []Comment
		EntriesOfFriends  []Entry
		CommentsOfFriends []CommentWithEntry
		Friends           []Friend
		Footprints        []Footprint
	}{
		*user, prof, entries, commentsForMe, entriesOfFriends, commentsOfFriends, friends, footprints,
	})
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	if !authenticated(w, r) {
		return
	}

	account := mux.Vars(r)["account_name"]
	owner := getUserFromAccount(w, account)
	row := db.QueryRow(`SELECT * FROM profiles WHERE user_id = ?`, owner.ID)
	prof := Profile{}
	err := row.Scan(&prof.UserID, &prof.FirstName, &prof.LastName, &prof.Sex, &prof.Birthday, &prof.Pref, &prof.UpdatedAt)
	if err != sql.ErrNoRows {
		checkErr(err)
	}
	var query string
	if permitted(w, r, owner.ID) {
		query = `SELECT * FROM entries WHERE user_id = ? ORDER BY created_at LIMIT 5`
	} else {
		query = `SELECT * FROM entries WHERE user_id = ? AND private=0 ORDER BY created_at LIMIT 5`
	}
	rows, err := db.Query(query, owner.ID)
	if err != sql.ErrNoRows {
		checkErr(err)
	}
	entries := make([]Entry, 0, 5)
	for rows.Next() {
		var id, userID, private int
		var body string
		var createdAt time.Time
		checkErr(rows.Scan(&id, &userID, &private, &body, &createdAt))
		entry := Entry{id, userID, private == 1, strings.SplitN(body, "\n", 2)[0], strings.SplitN(body, "\n", 2)[1], createdAt}
		entries = append(entries, entry)
	}
	rows.Close()

	markFootprint(w, r, owner.ID)

	render(w, r, http.StatusOK, "profile.html", struct {
		Owner   User
		Profile Profile
		Entries []Entry
		Private bool
	}{
		*owner, prof, entries, permitted(w, r, owner.ID),
	})
}

func PostProfile(w http.ResponseWriter, r *http.Request) {
	if !authenticated(w, r) {
		return
	}
	user := getCurrentUser(w, r)
	account := mux.Vars(r)["account_name"]
	if account != user.AccountName {
		checkErr(ErrPermissionDenied)
	}
	query := `UPDATE profiles
SET first_name=?, last_name=?, sex=?, birthday=?, pref=?, updated_at=CURRENT_TIMESTAMP()
WHERE user_id = ?`
	birth := r.FormValue("birthday")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	sex := r.FormValue("sex")
	pref := r.FormValue("pref")
	_, err := db.Exec(query, firstName, lastName, sex, birth, pref, user.ID)
	checkErr(err)
	// TODO should escape the account name?
	http.Redirect(w, r, "/profile/"+account, http.StatusSeeOther)
}

func ListEntries(w http.ResponseWriter, r *http.Request) {
	if !authenticated(w, r) {
		return
	}

	account := mux.Vars(r)["account_name"]
	owner := getUserFromAccount(w, account)
	var query string
	if permitted(w, r, owner.ID) {
		query = `SELECT * FROM entries WHERE user_id = ? ORDER BY created_at DESC LIMIT 20`
	} else {
		query = `SELECT * FROM entries WHERE user_id = ? AND private=0 ORDER BY created_at DESC LIMIT 20`
	}
	rows, err := db.Query(query, owner.ID)
	if err != sql.ErrNoRows {
		checkErr(err)
	}
	entries := make([]Entry, 0, 20)
	for rows.Next() {
		var id, userID, private int
		var body string
		var createdAt time.Time
		checkErr(rows.Scan(&id, &userID, &private, &body, &createdAt))
		entry := Entry{id, userID, private == 1, strings.SplitN(body, "\n", 2)[0], strings.SplitN(body, "\n", 2)[1], createdAt}
		entries = append(entries, entry)
	}
	rows.Close()

	markFootprint(w, r, owner.ID)

	render(w, r, http.StatusOK, "entries.html", struct {
		Owner   *User
		Entries []Entry
		Myself  bool
	}{owner, entries, getCurrentUser(w, r).ID == owner.ID})
}

func GetEntry(w http.ResponseWriter, r *http.Request) {
	if !authenticated(w, r) {
		return
	}
	entryID := mux.Vars(r)["entry_id"]
	row := db.QueryRow(`SELECT * FROM entries WHERE id = ?`, entryID)
	var id, userID, private int
	var body string
	var createdAt time.Time
	err := row.Scan(&id, &userID, &private, &body, &createdAt)
	if err == sql.ErrNoRows {
		checkErr(ErrContentNotFound)
	}
	checkErr(err)
	entry := Entry{id, userID, private == 1, strings.SplitN(body, "\n", 2)[0], strings.SplitN(body, "\n", 2)[1], createdAt}
	owner := getUser(w, entry.UserID)
	if entry.Private {
		if !permitted(w, r, owner.ID) {
			checkErr(ErrPermissionDenied)
		}
	}
	rows, err := db.Query(`SELECT * FROM comments WHERE entry_id = ?`, entry.ID)
	if err != sql.ErrNoRows {
		checkErr(err)
	}
	comments := make([]Comment, 0, 10)
	for rows.Next() {
		c := Comment{}
		checkErr(rows.Scan(&c.ID, &c.EntryID, &c.UserID, &c.Comment, &c.CreatedAt))
		comments = append(comments, c)
	}
	rows.Close()

	markFootprint(w, r, owner.ID)

	render(w, r, http.StatusOK, "entry.html", struct {
		Owner    *User
		Entry    Entry
		Comments []Comment
	}{owner, entry, comments})
}

func PostEntry(w http.ResponseWriter, r *http.Request) {
	if !authenticated(w, r) {
		return
	}

	user := getCurrentUser(w, r)
	title := r.FormValue("title")
	if title == "" {
		title = "タイトルなし"
	}
	content := r.FormValue("content")
	var private int
	if r.FormValue("private") == "" {
		private = 0
	} else {
		private = 1
	}
	_, err := db.Exec(`INSERT INTO entries (user_id, private, body) VALUES (?,?,?)`, user.ID, private, title+"\n"+content)
	checkErr(err)
	http.Redirect(w, r, "/diary/entries/"+user.AccountName, http.StatusSeeOther)
}

func PostComment(w http.ResponseWriter, r *http.Request) {
	if !authenticated(w, r) {
		return
	}

	entryID := mux.Vars(r)["entry_id"]
	row := db.QueryRow(`SELECT * FROM entries WHERE id = ?`, entryID)
	var id, userID, private int
	var body string
	var createdAt time.Time
	err := row.Scan(&id, &userID, &private, &body, &createdAt)
	if err == sql.ErrNoRows {
		checkErr(ErrContentNotFound)
	}
	checkErr(err)

	entry := Entry{id, userID, private == 1, strings.SplitN(body, "\n", 2)[0], strings.SplitN(body, "\n", 2)[1], createdAt}
	owner := getUser(w, entry.UserID)
	if entry.Private {
		if !permitted(w, r, owner.ID) {
			checkErr(ErrPermissionDenied)
		}
	}
	user := getCurrentUser(w, r)

	res, err := db.Exec(`INSERT INTO comments (entry_id, user_id, comment) VALUES (?,?,?)`, entry.ID, user.ID, r.FormValue("comment"))
	checkErr(err)
	lastID, err := res.LastInsertId()
	checkErr(err)
	commentCh <- Comment{int(lastID), entry.ID, user.ID, r.FormValue("comment"), time.Now()}

	http.Redirect(w, r, "/diary/entry/"+strconv.Itoa(entry.ID), http.StatusSeeOther)
}

func GetFootprints(w http.ResponseWriter, r *http.Request) {
	if !authenticated(w, r) {
		return
	}

	user := getCurrentUser(w, r)
	footprints, err := getFootprintsFromRedis(user.ID, 50)
	if err != nil {
		checkErr(err)
	}
	render(w, r, http.StatusOK, "footprints.html", struct{ Footprints []Footprint }{footprints})
}
func GetFriends(w http.ResponseWriter, r *http.Request) {
	if !authenticated(w, r) {
		return
	}

	user := getCurrentUser(w, r)
	rows, err := db.Query(`SELECT * FROM relations WHERE one = ? OR another = ? ORDER BY created_at DESC`, user.ID, user.ID)
	if err != sql.ErrNoRows {
		checkErr(err)
	}
	friendsMap := make(map[int]time.Time)
	for rows.Next() {
		var id, one, another int
		var createdAt time.Time
		checkErr(rows.Scan(&id, &one, &another, &createdAt))
		var friendID int
		if one == user.ID {
			friendID = another
		} else {
			friendID = one
		}
		if _, ok := friendsMap[friendID]; !ok {
			friendsMap[friendID] = createdAt
		}
	}
	rows.Close()
	friends := make([]Friend, 0, len(friendsMap))
	for key, val := range friendsMap {
		friends = append(friends, Friend{key, val})
	}
	render(w, r, http.StatusOK, "friends.html", struct{ Friends []Friend }{friends})
}

func PostFriends(w http.ResponseWriter, r *http.Request) {
	if !authenticated(w, r) {
		return
	}

	user := getCurrentUser(w, r)
	anotherAccount := mux.Vars(r)["account_name"]
	if !isFriendAccount(w, r, anotherAccount) {
		another := getUserFromAccount(w, anotherAccount)
		_, err := db.Exec(`INSERT INTO relations (one, another) VALUES (?,?), (?,?)`, user.ID, another.ID, another.ID, user.ID)
		checkErr(err)
		http.Redirect(w, r, "/friends", http.StatusSeeOther)
	}
}

func GetInitialize(w http.ResponseWriter, r *http.Request) {
	db.Exec("DELETE FROM relations WHERE id > 500000")
	db.Exec("DELETE FROM footprints WHERE id > 500000")
	db.Exec("DELETE FROM entries WHERE id > 500000")
	db.Exec("DELETE FROM comments WHERE id > 1500000")
	initCommentsForMe()
	initAllUsers()
}

func init() {
	flag.Parse()

}

var commentCh = make(chan Comment, 10000)
var commentsForMeTable = make(map[int][]Comment)
var commentsForMeMutex = sync.Mutex{}

func commentsForMeWorker() {
	for {
		comment := <-commentCh
		commentsForMeMutex.Lock()
		_, ok := commentsForMeTable[comment.UserID]
		if !ok {
			commentsForMeTable[comment.UserID] = make([]Comment, 0, 100)
		}
		commentsForMeTable[comment.UserID] = append(commentsForMeTable[comment.UserID], comment)
		commentsForMeMutex.Unlock()
	}
}

func getCommentsForMe(userID int, num int) []Comment {
	commentsForMeMutex.Lock()
	defer commentsForMeMutex.Unlock()

	comments, ok := commentsForMeTable[userID]
	if !ok {
		return make([]Comment, 0)
	}
	coms := make([]Comment, num)
	for i, c := range comments[len(comments)-num:] {
		coms[num-i-1] = c
	}
	return coms
}

func initCommentsForMe() {
	commentsForMeMutex.Lock()
	commentsForMeTable = make(map[int][]Comment)
	commentsForMeMutex.Unlock()

	fmt.Println("start initCommentsForMe")

	rows, err := db.Query(`SELECT * FROM comments ORDER BY created_at ASC`)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	for rows.Next() {
		c := Comment{}
		err = rows.Scan(&c.ID, &c.EntryID, &c.UserID, &c.Comment, &c.CreatedAt)
		if err != nil {
			panic(err)
		}
		commentCh <- c
	}
	rows.Close()

	fmt.Println("finish initCommentsForMe")
}

func main() {
	host := os.Getenv("ISUCON5_DB_HOST")
	if host == "" {
		host = "localhost"
	}
	portstr := os.Getenv("ISUCON5_DB_PORT")
	if portstr == "" {
		portstr = "3306"
	}
	port, err := strconv.Atoi(portstr)
	if err != nil {
		log.Fatalf("Failed to read DB port number from an environment variable ISUCON5_DB_PORT.\nError: %s", err.Error())
	}
	user := os.Getenv("ISUCON5_DB_USER")
	if user == "" {
		user = "root"
	}
	password := os.Getenv("ISUCON5_DB_PASSWORD")
	dbname := os.Getenv("ISUCON5_DB_NAME")
	if dbname == "" {
		dbname = "isucon5q"
	}
	ssecret := os.Getenv("ISUCON5_SESSION_SECRET")
	if ssecret == "" {
		ssecret = "beermoris"
	}

	var dsn string
	if *mysqlsock {
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
			user,
			password,
			"instance-1",
			port,
			dbname,
		)

	} else {
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
			user,
			password,
			"localhost",
			port,
			dbname,
		)
	}
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %s.", err.Error())
	}
	defer db.Close()

	if *mysqlsock {
		rd = redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr:    "instance-1",
			DB:      0,
		})
	} else {
		rd = redis.NewClient(&redis.Options{
			Network: "unix",
			Addr:    "/tmp/redis.sock",
			DB:      0,
		})
	}
	store = sessions.NewCookieStore([]byte(ssecret))

	r := mux.NewRouter()

	l := r.Path("/login").Subrouter()
	l.Methods("GET").HandlerFunc(myHandler(GetLogin))
	l.Methods("POST").HandlerFunc(myHandler(PostLogin))
	r.Path("/logout").Methods("GET").HandlerFunc(myHandler(GetLogout))

	p := r.Path("/profile/{account_name}").Subrouter()
	p.Methods("GET").HandlerFunc(myHandler(GetProfile))
	p.Methods("POST").HandlerFunc(myHandler(PostProfile))

	d := r.PathPrefix("/diary").Subrouter()
	d.HandleFunc("/entries/{account_name}", myHandler(ListEntries)).Methods("GET")
	d.HandleFunc("/entry", myHandler(PostEntry)).Methods("POST")
	d.HandleFunc("/entry/{entry_id}", myHandler(GetEntry)).Methods("GET")

	d.HandleFunc("/comment/{entry_id}", myHandler(PostComment)).Methods("POST")

	r.HandleFunc("/footprints", myHandler(GetFootprints)).Methods("GET")

	r.HandleFunc("/friends", myHandler(GetFriends)).Methods("GET")
	r.HandleFunc("/friends/{account_name}", myHandler(PostFriends)).Methods("POST")

	r.HandleFunc("/initialize", myHandler(GetInitialize))
	r.HandleFunc("/", myHandler(GetIndex))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../static")))

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, syscall.SIGTERM)
	signal.Notify(sigchan, syscall.SIGINT)

	go commentsForMeWorker()
	initCommentsForMe() // edvakf
	initAllUsers()

	var li net.Listener
	// sock := "/dev/shm/server.sock"
	if *tport == 0 {
		// ferr := os.Remove(sock)
		// if ferr != nil {
		// 	if !os.IsNotExist(ferr) {
		// 		panic(ferr.Error())
		// 	}
		// }
		// li, err = net.Listen("unix", sock)
		// cerr := os.Chmod(sock, 0666)
		// if cerr != nil {
		// 	panic(cerr.Error())
		// }
		li, err = net.ListenTCP("tcp", &net.TCPAddr{Port: 8080})
	} else {
		li, err = net.ListenTCP("tcp", &net.TCPAddr{Port: int(*tport)})
	}
	if err != nil {
		panic(err.Error())
	}
	go func() {
		// func Serve(l net.Listener, handler Handler) error
		log.Println(http.Serve(li, r))
	}()

	<-sigchan

	// log.Fatal(http.ListenAndServe(":8080", r))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
