// Generated by ego.
// DO NOT EDIT

//line index.ego:2
package main

import "fmt"
import "html"
import "io"

func MyTmpl(w io.Writer, e struct {
	User              User
	Profile           Profile
	Entries           []Entry
	CommentsForMe     []Comment
	EntriesOfFriends  []Entry
	CommentsOfFriends []CommentWithEntry
	Friends           []Friend
	Footprints        []Footprint
}) {
//line index.ego:14
	_, _ = io.WriteString(w, "<!DOCTYPE html>\n<html>\n<head>\n    <meta http-equiv=\"Content-Type\" content=\"text/html\" charset=\"utf-8\">\n    <link rel=\"stylesheet\" href=\"/css/bootstrap.min.css\">\n    <title>ISUxi</title>\n</head>\n\n<body class=\"container\">\n<h1 class=\"jumbotron\"><a href=\"/\">ISUxiへようこそ!</a></h1>\n<h2>ISUxi index</h2>\n<div class=\"row panel panel-primary\" id=\"prof\">\n  <div class=\"col-md-12 panel-title\" id=\"prof-nickname\">")
//line index.ego:26
	_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(e.User.NickName)))
//line index.ego:26
	_, _ = io.WriteString(w, "</div>\n  <div class=\"col-md-12\"><a href=\"/profile/")
//line index.ego:27
	_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(e.User.AccountName)))
//line index.ego:27
	_, _ = io.WriteString(w, "\">プロフィール</a></div>\n  <div class=\"col-md-4\">\n    <dl>\n      <dt>アカウント名</dt><dd id=\"prof-account-name\">")
//line index.ego:30
	_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(e.User.AccountName)))
//line index.ego:30
	_, _ = io.WriteString(w, "</dd>\n      <dt>メールアドレス</dt><dd id=\"prof-email\">")
//line index.ego:31
	_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(e.User.Email)))
//line index.ego:31
	_, _ = io.WriteString(w, "</dd>\n      <dt>姓</dt><dd id=\"prof-last-name\">")
//line index.ego:32
	if e.Profile.LastName != "" {
//line index.ego:32
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(e.Profile.LastName)))
//line index.ego:32
	} else {
//line index.ego:32
		_, _ = io.WriteString(w, "未入力")
//line index.ego:32
	}
//line index.ego:32
	_, _ = io.WriteString(w, "</dd>\n      <dt>名</dt><dd id=\"prof-first-name\">")
//line index.ego:33
	if e.Profile.FirstName != "" {
//line index.ego:33
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(e.Profile.FirstName)))
//line index.ego:33
	} else {
//line index.ego:33
		_, _ = io.WriteString(w, "未入力")
//line index.ego:33
	}
//line index.ego:33
	_, _ = io.WriteString(w, "</dd>\n      <dt>性別</dt><dd id=\"prof-sex\">")
//line index.ego:34
	if e.Profile.Sex != "" {
//line index.ego:34
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(e.Profile.Sex)))
//line index.ego:34
	} else {
//line index.ego:34
		_, _ = io.WriteString(w, "未入力")
//line index.ego:34
	}
//line index.ego:34
	_, _ = io.WriteString(w, "</dd>\n      <dt>誕生日</dt><dd id=\"prof-birthday\">")
//line index.ego:35
	if e.Profile.Birthday.Valid {
//line index.ego:35
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(e.Profile.Birthday.Time.Format("1月2日"))))
//line index.ego:35
	} else {
//line index.ego:35
		_, _ = io.WriteString(w, "未入力")
//line index.ego:35
	}
//line index.ego:35
	_, _ = io.WriteString(w, "</dd>\n      <dt>住んでいる県</dt><dd id=\"prof-pref\">")
//line index.ego:36
	if e.Profile.Pref != "" {
//line index.ego:36
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(e.Profile.Pref)))
//line index.ego:36
	} else {
//line index.ego:36
		_, _ = io.WriteString(w, "未入力")
//line index.ego:36
	}
//line index.ego:36
	_, _ = io.WriteString(w, "</dd>\n      <dt>友だちの人数</dt><dd id=\"prof-friends\"><a href=\"/friends\">")
//line index.ego:37
	_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(len(e.Friends))))
//line index.ego:37
	_, _ = io.WriteString(w, "人</a></dd>\n    </dl>\n  </div>\n\n  <div class=\"col-md-4\">\n    <div id=\"entries-title\"><a href=\"/diary/entries/")
//line index.ego:42
	_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(e.User.AccountName)))
//line index.ego:42
	_, _ = io.WriteString(w, "\">あなたの日記エントリ</a></div>\n    <div id=\"entries\">\n      <ul class=\"list-group\">\n        ")
//line index.ego:45
	for _, entry := range e.Entries {
//line index.ego:46
		_, _ = io.WriteString(w, "\n        <li class=\"list-group-item entries-entry\"><a href=\"/diary/entry/")
//line index.ego:46
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(entry.ID)))
//line index.ego:46
		_, _ = io.WriteString(w, "\">")
//line index.ego:46
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(entry.Title)))
//line index.ego:46
		_, _ = io.WriteString(w, "</a></li>\n        ")
//line index.ego:47
	}
//line index.ego:48
	_, _ = io.WriteString(w, "\n      </ul>\n    </div>\n  </div>\n\n  <div class=\"col-md-4\">\n    <div><a href=\"/footprints\">あなたのページへの足あと</a></div>\n    <div id=\"footprints\">\n      <ul class=\"list-group\">\n        ")
//line index.ego:56
	for _, f := range e.Footprints {
//line index.ego:57
		_, _ = io.WriteString(w, "\n        ")
//line index.ego:57
		owner := getUser(f.OwnerID)
//line index.ego:58
		_, _ = io.WriteString(w, "\n        <li class=\"list-group-item footprints-footprint\">")
//line index.ego:58
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(f.CreatedAt.Format("2006-01-02 15:04:05"))))
//line index.ego:58
		_, _ = io.WriteString(w, ": <a href=\"/profile/")
//line index.ego:58
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(owner.AccountName)))
//line index.ego:58
		_, _ = io.WriteString(w, "\">")
//line index.ego:58
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(owner.NickName)))
//line index.ego:58
		_, _ = io.WriteString(w, "さん</a></li>\n        ")
//line index.ego:59
	}
//line index.ego:60
	_, _ = io.WriteString(w, "\n      </ul>\n    </div>\n  </div>\n</div>\n\n<div class=\"row panel panel-primary\">\n  <div class=\"col-md-4\">\n    <div>あなたへのコメント</div>\n    <div id=\"comments\">\n      ")
//line index.ego:69
	for _, c := range e.CommentsForMe {
//line index.ego:70
		_, _ = io.WriteString(w, "\n      <div class=\"comments-comment\">\n        <ul class=\"list-group\">\n          ")
//line index.ego:72
		commentUser := getUser(c.UserID)
//line index.ego:73
		_, _ = io.WriteString(w, "\n          <li class=\"list-group-item comment-owner\"><a href=\"/profile/")
//line index.ego:73
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(commentUser.AccountName)))
//line index.ego:73
		_, _ = io.WriteString(w, "\">")
//line index.ego:73
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(commentUser.NickName)))
//line index.ego:73
		_, _ = io.WriteString(w, "さん</a>:</li>\n          <li class=\"list-group-item comment-comment\">")
//line index.ego:74
		if len(c.Comment) >= 30 {
//line index.ego:74
			_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(c.Comment[:27])))
//line index.ego:74
			_, _ = io.WriteString(w, "...")
//line index.ego:74
		} else {
//line index.ego:74
			_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(c.Comment)))
//line index.ego:74
		}
//line index.ego:74
		_, _ = io.WriteString(w, "</li>\n          <li class=\"list-group-item comment-created-at\">投稿時刻:")
//line index.ego:75
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(c.CreatedAt.Format("2006-01-02 15:04:05"))))
//line index.ego:75
		_, _ = io.WriteString(w, "</li>\n        </ul>\n      </div>\n      ")
//line index.ego:78
	}
//line index.ego:79
	_, _ = io.WriteString(w, "\n    </div>\n  </div>\n\n  <div class=\"col-md-4\">\n    <div>あなたの友だちの日記エントリ</div>\n    <div id=\"friend-entries\">\n      ")
//line index.ego:85
	for _, ef := range e.EntriesOfFriends {
//line index.ego:86
		_, _ = io.WriteString(w, "\n      <div class=\"friend-entry\">\n        <ul class=\"list-group\">\n          ")
//line index.ego:88
		entryOwner := getUser(ef.UserID)
//line index.ego:89
		_, _ = io.WriteString(w, "\n          <li class=\"list-group-item entry-owner\"><a href=\"/diary/entries/")
//line index.ego:89
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(entryOwner.AccountName)))
//line index.ego:89
		_, _ = io.WriteString(w, "\">")
//line index.ego:89
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(entryOwner.NickName)))
//line index.ego:89
		_, _ = io.WriteString(w, "さん</a>:</li>\n          <li class=\"list-group-item entry-title\"><a href=\"/diary/entry/")
//line index.ego:90
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(ef.ID)))
//line index.ego:90
		_, _ = io.WriteString(w, "\">")
//line index.ego:90
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(ef.Title)))
//line index.ego:90
		_, _ = io.WriteString(w, "</a></li>\n          <li class=\"list-group-item entry-created-at\">投稿時刻:")
//line index.ego:91
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(ef.CreatedAt.Format("2006-01-02 15:04:05"))))
//line index.ego:91
		_, _ = io.WriteString(w, "</li>\n        </ul>\n      </div>\n      ")
//line index.ego:94
	}
//line index.ego:95
	_, _ = io.WriteString(w, "\n    </div>\n  </div>\n\n  <div class=\"col-md-4\">\n    <div>あなたの友だちのコメント</div>\n    <div id=\"friend-comments\">\n      ")
//line index.ego:101
	for _, c := range e.CommentsOfFriends {
//line index.ego:102
		_, _ = io.WriteString(w, "\n      <div class=\"friend-comment\">\n        <ul class=\"list-group\">\n          ")
//line index.ego:104
		commentOwner := getUser(c.UserID)
//line index.ego:105
		_, _ = io.WriteString(w, "\n          ")
//line index.ego:105
		entryOwner := getUser(c.Entry.UserID)
//line index.ego:106
		_, _ = io.WriteString(w, "\n          <li class=\"list-group-item comment-from-to\"><a href=\"/profile/")
//line index.ego:106
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(commentOwner.AccountName)))
//line index.ego:106
		_, _ = io.WriteString(w, "\">")
//line index.ego:106
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(commentOwner.NickName)))
//line index.ego:106
		_, _ = io.WriteString(w, "さん</a>から<a href=\"/profile/")
//line index.ego:106
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(entryOwner.AccountName)))
//line index.ego:106
		_, _ = io.WriteString(w, "\">")
//line index.ego:106
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(entryOwner.NickName)))
//line index.ego:106
		_, _ = io.WriteString(w, "さん</a>へのコメント:</li>\n          <li class=\"list-group-item comment-comment\">")
//line index.ego:107
		if len(c.Comment) >= 30 {
//line index.ego:107
			_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(c.Comment[:27])))
//line index.ego:107
			_, _ = io.WriteString(w, "...")
//line index.ego:107
		} else {
//line index.ego:107
			_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(c.Comment)))
//line index.ego:107
		}
//line index.ego:107
		_, _ = io.WriteString(w, "</li>\n          <li class=\"list-group-item comment-created-at\">投稿時刻:")
//line index.ego:108
		_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(c.CreatedAt.Format("2006-01-02 15:04:05"))))
//line index.ego:108
		_, _ = io.WriteString(w, "</li>\n        </ul>\n      </div>\n      ")
//line index.ego:111
	}
//line index.ego:112
	_, _ = io.WriteString(w, "\n    </div>\n  </div>\n</div>\n\n</body>\n</html>\n")
//line index.ego:118
}

var _ fmt.Stringer
var _ io.Reader
var _ = html.EscapeString
