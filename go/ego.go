// Generated by ego.
// DO NOT EDIT

package main
import (
"fmt"
"html"
"io"
)
var _ = fmt.Sprint("") // just so that we can keep the fmt import for now
//line templates/index.ego:1
 func MyTmpl(w io.Writer, e struct {
	User              User
	Profile           Profile
	Entries           []Entry
	CommentsForMe     []Comment
	EntriesOfFriends  []Entry
	CommentsOfFriends []CommentWithEntry
	Friends           []Friend
	Footprints        []Footprint
}) error  {
//line templates/index.ego:10
_, _ = io.WriteString(w, "<!DOCTYPE html>\n<html>\n<head>\n    <meta http-equiv=\"Content-Type\" content=\"text/html\" charset=\"utf-8\">\n    <link rel=\"stylesheet\" href=\"/css/bootstrap.min.css\">\n    <title>ISUxi</title>\n</head>\n\n<body class=\"container\">\n<h1 class=\"jumbotron\"><a href=\"/\">ISUxiへようこそ!</a></h1>\n<h2>ISUxi index</h2>\n<div class=\"row panel panel-primary\" id=\"prof\">\n  <div class=\"col-md-12 panel-title\" id=\"prof-nickname\">")
//line templates/index.ego:22
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  e.User.NickName )))
//line templates/index.ego:22
_, _ = io.WriteString(w, "</div>\n  <div class=\"col-md-12\"><a href=\"/profile/")
//line templates/index.ego:23
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  e.User.AccountName )))
//line templates/index.ego:23
_, _ = io.WriteString(w, "\">プロフィール</a></div>\n  <div class=\"col-md-4\">\n    <dl>\n      <dt>アカウント名</dt><dd id=\"prof-account-name\">")
//line templates/index.ego:26
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  e.User.AccountName )))
//line templates/index.ego:26
_, _ = io.WriteString(w, "</dd>\n      <dt>メールアドレス</dt><dd id=\"prof-email\">")
//line templates/index.ego:27
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  e.User.Email )))
//line templates/index.ego:27
_, _ = io.WriteString(w, "</dd>\n      <dt>姓</dt><dd id=\"prof-last-name\">")
//line templates/index.ego:28
 if e.Profile.LastName != "" { 
//line templates/index.ego:28
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  e.Profile.LastName )))
//line templates/index.ego:28
 } else { 
//line templates/index.ego:28
_, _ = io.WriteString(w, "未入力")
//line templates/index.ego:28
 } 
//line templates/index.ego:28
_, _ = io.WriteString(w, "</dd>\n      <dt>名</dt><dd id=\"prof-first-name\">")
//line templates/index.ego:29
 if e.Profile.FirstName != "" { 
//line templates/index.ego:29
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  e.Profile.FirstName )))
//line templates/index.ego:29
 } else { 
//line templates/index.ego:29
_, _ = io.WriteString(w, "未入力")
//line templates/index.ego:29
 } 
//line templates/index.ego:29
_, _ = io.WriteString(w, "</dd>\n      <dt>性別</dt><dd id=\"prof-sex\">")
//line templates/index.ego:30
 if e.Profile.Sex != "" { 
//line templates/index.ego:30
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  e.Profile.Sex )))
//line templates/index.ego:30
 } else { 
//line templates/index.ego:30
_, _ = io.WriteString(w, "未入力")
//line templates/index.ego:30
 } 
//line templates/index.ego:30
_, _ = io.WriteString(w, "</dd>\n      <dt>誕生日</dt><dd id=\"prof-birthday\">")
//line templates/index.ego:31
 if e.Profile.Birthday.Valid { 
//line templates/index.ego:31
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  e.Profile.Birthday.Time.Format("1月2日") )))
//line templates/index.ego:31
 } else { 
//line templates/index.ego:31
_, _ = io.WriteString(w, "未入力")
//line templates/index.ego:31
 } 
//line templates/index.ego:31
_, _ = io.WriteString(w, "</dd>\n      <dt>住んでいる県</dt><dd id=\"prof-pref\">")
//line templates/index.ego:32
 if e.Profile.Pref != "" { 
//line templates/index.ego:32
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  e.Profile.Pref )))
//line templates/index.ego:32
 } else { 
//line templates/index.ego:32
_, _ = io.WriteString(w, "未入力")
//line templates/index.ego:32
 } 
//line templates/index.ego:32
_, _ = io.WriteString(w, "</dd>\n      <dt>友だちの人数</dt><dd id=\"prof-friends\"><a href=\"/friends\">")
//line templates/index.ego:33
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  len(e.Friends) )))
//line templates/index.ego:33
_, _ = io.WriteString(w, "人</a></dd>\n    </dl>\n  </div>\n\n  <div class=\"col-md-4\">\n    <div id=\"entries-title\"><a href=\"/diary/entries/")
//line templates/index.ego:38
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  e.User.AccountName )))
//line templates/index.ego:38
_, _ = io.WriteString(w, "\">あなたの日記エントリ</a></div>\n    <div id=\"entries\">\n      <ul class=\"list-group\">\n        ")
//line templates/index.ego:41
 for _, entry := range e.Entries { 
//line templates/index.ego:42
_, _ = io.WriteString(w, "\n        <li class=\"list-group-item entries-entry\"><a href=\"/diary/entry/")
//line templates/index.ego:42
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  entry.ID )))
//line templates/index.ego:42
_, _ = io.WriteString(w, "\">")
//line templates/index.ego:42
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  entry.Title )))
//line templates/index.ego:42
_, _ = io.WriteString(w, "</a></li>\n        ")
//line templates/index.ego:43
 } 
//line templates/index.ego:44
_, _ = io.WriteString(w, "\n      </ul>\n    </div>\n  </div>\n\n  <div class=\"col-md-4\">\n    <div><a href=\"/footprints\">あなたのページへの足あと</a></div>\n    <div id=\"footprints\">\n      <ul class=\"list-group\">\n        ")
//line templates/index.ego:52
 for _, f := range e.Footprints { 
//line templates/index.ego:53
_, _ = io.WriteString(w, "\n        ")
//line templates/index.ego:53
 owner := getUser(f.OwnerID) 
//line templates/index.ego:54
_, _ = io.WriteString(w, "\n        <li class=\"list-group-item footprints-footprint\">")
//line templates/index.ego:54
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  f.CreatedAt.Format("2006-01-02 15:04:05") )))
//line templates/index.ego:54
_, _ = io.WriteString(w, ": <a href=\"/profile/")
//line templates/index.ego:54
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  owner.AccountName )))
//line templates/index.ego:54
_, _ = io.WriteString(w, "\">")
//line templates/index.ego:54
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  owner.NickName )))
//line templates/index.ego:54
_, _ = io.WriteString(w, "さん</a></li>\n        ")
//line templates/index.ego:55
 } 
//line templates/index.ego:56
_, _ = io.WriteString(w, "\n      </ul>\n    </div>\n  </div>\n</div>\n\n<div class=\"row panel panel-primary\">\n  <div class=\"col-md-4\">\n    <div>あなたへのコメント</div>\n    <div id=\"comments\">\n      ")
//line templates/index.ego:65
 for _, c := range e.CommentsForMe { 
//line templates/index.ego:66
_, _ = io.WriteString(w, "\n      <div class=\"comments-comment\">\n        <ul class=\"list-group\">\n          ")
//line templates/index.ego:68
 commentUser := getUser(c.UserID) 
//line templates/index.ego:69
_, _ = io.WriteString(w, "\n          <li class=\"list-group-item comment-owner\"><a href=\"/profile/")
//line templates/index.ego:69
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  commentUser.AccountName)))
//line templates/index.ego:69
_, _ = io.WriteString(w, "\">")
//line templates/index.ego:69
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  commentUser.NickName )))
//line templates/index.ego:69
_, _ = io.WriteString(w, "さん</a>:</li>\n          <li class=\"list-group-item comment-comment\">")
//line templates/index.ego:70
 if len(c.Comment) >= 30 { 
//line templates/index.ego:70
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  c.Comment[:27] )))
//line templates/index.ego:70
_, _ = io.WriteString(w, "...")
//line templates/index.ego:70
 } else { 
//line templates/index.ego:70
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  c.Comment )))
//line templates/index.ego:70
 } 
//line templates/index.ego:70
_, _ = io.WriteString(w, "</li>\n          <li class=\"list-group-item comment-created-at\">投稿時刻:")
//line templates/index.ego:71
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  c.CreatedAt.Format("2006-01-02 15:04:05") )))
//line templates/index.ego:71
_, _ = io.WriteString(w, "</li>\n        </ul>\n      </div>\n      ")
//line templates/index.ego:74
 } 
//line templates/index.ego:75
_, _ = io.WriteString(w, "\n    </div>\n  </div>\n\n  <div class=\"col-md-4\">\n    <div>あなたの友だちの日記エントリ</div>\n    <div id=\"friend-entries\">\n      ")
//line templates/index.ego:81
 for _, ef := range e.EntriesOfFriends { 
//line templates/index.ego:82
_, _ = io.WriteString(w, "\n      <div class=\"friend-entry\">\n        <ul class=\"list-group\">\n          ")
//line templates/index.ego:84
 entryOwner := getUser(ef.UserID) 
//line templates/index.ego:85
_, _ = io.WriteString(w, "\n          <li class=\"list-group-item entry-owner\"><a href=\"/diary/entries/")
//line templates/index.ego:85
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  entryOwner.AccountName )))
//line templates/index.ego:85
_, _ = io.WriteString(w, "\">")
//line templates/index.ego:85
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  entryOwner.NickName )))
//line templates/index.ego:85
_, _ = io.WriteString(w, "さん</a>:</li>\n          <li class=\"list-group-item entry-title\"><a href=\"/diary/entry/")
//line templates/index.ego:86
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  ef.ID )))
//line templates/index.ego:86
_, _ = io.WriteString(w, "\">")
//line templates/index.ego:86
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  ef.Title )))
//line templates/index.ego:86
_, _ = io.WriteString(w, "</a></li>\n          <li class=\"list-group-item entry-created-at\">投稿時刻:")
//line templates/index.ego:87
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  ef.CreatedAt.Format("2006-01-02 15:04:05") )))
//line templates/index.ego:87
_, _ = io.WriteString(w, "</li>\n        </ul>\n      </div>\n      ")
//line templates/index.ego:90
 } 
//line templates/index.ego:91
_, _ = io.WriteString(w, "\n    </div>\n  </div>\n\n  <div class=\"col-md-4\">\n    <div>あなたの友だちのコメント</div>\n    <div id=\"friend-comments\">\n      ")
//line templates/index.ego:97
 for _, c := range e.CommentsOfFriends { 
//line templates/index.ego:98
_, _ = io.WriteString(w, "\n      <div class=\"friend-comment\">\n        <ul class=\"list-group\">\n          ")
//line templates/index.ego:100
 commentOwner := getUser(c.UserID) 
//line templates/index.ego:101
_, _ = io.WriteString(w, "\n          ")
//line templates/index.ego:101
 entryOwner := getUser(c.Entry.UserID) 
//line templates/index.ego:102
_, _ = io.WriteString(w, "\n          <li class=\"list-group-item comment-from-to\"><a href=\"/profile/")
//line templates/index.ego:102
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  commentOwner.AccountName )))
//line templates/index.ego:102
_, _ = io.WriteString(w, "\">")
//line templates/index.ego:102
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  commentOwner.NickName )))
//line templates/index.ego:102
_, _ = io.WriteString(w, "さん</a>から<a href=\"/profile/")
//line templates/index.ego:102
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  entryOwner.AccountName )))
//line templates/index.ego:102
_, _ = io.WriteString(w, "\">")
//line templates/index.ego:102
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  entryOwner.NickName )))
//line templates/index.ego:102
_, _ = io.WriteString(w, "さん</a>へのコメント:</li>\n          <li class=\"list-group-item comment-comment\">")
//line templates/index.ego:103
 if len(c.Comment) >= 30 { 
//line templates/index.ego:103
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  c.Comment[:27] )))
//line templates/index.ego:103
_, _ = io.WriteString(w, "...")
//line templates/index.ego:103
 } else { 
//line templates/index.ego:103
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  c.Comment )))
//line templates/index.ego:103
 } 
//line templates/index.ego:103
_, _ = io.WriteString(w, "</li>\n          <li class=\"list-group-item comment-created-at\">投稿時刻:")
//line templates/index.ego:104
_, _ = io.WriteString(w, html.EscapeString(fmt.Sprintf("%v",  c.CreatedAt.Format("2006-01-02 15:04:05") )))
//line templates/index.ego:104
_, _ = io.WriteString(w, "</li>\n        </ul>\n      </div>\n      ")
//line templates/index.ego:107
 } 
//line templates/index.ego:108
_, _ = io.WriteString(w, "\n    </div>\n  </div>\n</div>\n\n</body>\n</html>\n")
return nil
}
