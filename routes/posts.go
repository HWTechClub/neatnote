// Neat Note. A notes sharing platform for university students.
// Copyright (C) 2020 Humaid AlQassimi
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package routes

import (
	"bytes"
	"fmt"
	"git.sr.ht/~humaid/neatnote/models"
	"git.sr.ht/~humaid/neatnote/modules/common"
	"git.sr.ht/~humaid/neatnote/modules/namegen"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	macaron "gopkg.in/macaron.v1"
	"html/template"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	// simpleTextExp matches a simple text string.
	simpleTextExp = regexp.MustCompile(`^[a-zA-Z0-9 \-\_\.\#\/\?\,\+\&\:\(\)\[\]]+$`)
	// htmlTagExp roughly matches any HTML tag.
	htmlTagExp = regexp.MustCompile(`\<?(\/)?[a-zA-Z0-9 "=\n:\/\.\@\#\&\;\+\-\?\,\_]+\>`)
)

// CourseHandler response for a course page.
func CourseHandler(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	course, _ := models.GetCourse(ctx.Params("course"))

	course.LoadPosts()
	if ctx.Params("sort") == "top" {
		sort.Sort(models.TopPosts(course.Posts))
	} else if ctx.Params("sort") == "new" {
		sort.Sort(models.NewPosts(course.Posts))
	} else {
		sort.Sort(models.HotPosts(course.Posts))
	}
	ctx.Data["Course"] = course
	if !course.Locked {
		ctx.Data["PostButton"] = 1
	}
	ctx.Data["Title"] = course.Name
	ctx.HTML(200, "course")
}

// ReveaPosterHandler response for revealing the user.
func RevealPosterHandler(ctx *macaron.Context, x csrf.CSRF, sess session.Store, f *session.Flash) {
	p, _ := models.GetPost(ctx.Params("post"))
	poster, err := models.GetUser(p.PosterID)
	if err != nil {
		f.Error("Post no longer exists.")
	} else {
		f.Info(fmt.Sprintf("User: %s (%s)", poster.FullName, poster.Username))
	}

	ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
		ctx.Params("post")))
}

// EditCommentHandler response for a post page.
func EditCommentHandler(ctx *macaron.Context, x csrf.CSRF, sess session.Store, f *session.Flash) {
	comment, err := models.GetComment(ctx.Params("id"))
	if err != nil {
		f.Error("Comment no longer exists.")
		ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
			ctx.Params("post")))
		return
	}

	ctx.Data["Comment"] = comment

	u, _ := models.GetUser(sess.Get("user").(string))
	if !(comment.PosterID == sess.Get("user").(string) || u.IsAdmin) {
		f.Error("You may not edit this comment.")
		ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
			ctx.Params("post")))
		return
	}

	ctx.Data["csrf_token"] = x.GetToken()
	ctx.Data["Title"] = "Edit comment"
	ctx.HTML(200, "edit-comment")
}

// PostEditCommentHandler post response for editing a post.
func PostEditCommentHandler(ctx *macaron.Context, x csrf.CSRF, sess session.Store, f *session.Flash) {
	comment, err := models.GetComment(ctx.Params("id"))
	if err != nil {
		f.Error("Comment no longer exists.")
		ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
			ctx.Params("post")))
		return
	}

	u, _ := models.GetUser(sess.Get("user").(string))
	if !(comment.PosterID == sess.Get("user").(string) || u.IsAdmin) {
		f.Error("You may not edit this comment.")
		ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
			ctx.Params("post")))
		return
	}
	text := ctx.QueryTrim("text")

	if getMarkdownLength(ctx.QueryTrim("text")) < 2 {
		f.Error("The comment is empty or too short!")
		ctx.Redirect(fmt.Sprintf("/c/%s/%s/edit/%s", ctx.Params("course"),
			ctx.Params("post"), ctx.Params("id")))
		return
	}

	err = models.UpdateComment(&models.Comment{
		CommentID: comment.CommentID,
		Text:      text,
	})
	if err != nil {
		panic(err)
	}

	ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
		ctx.Params("post")))
}

// EditPostHandler response for a post page.
func EditPostHandler(ctx *macaron.Context, x csrf.CSRF, sess session.Store, f *session.Flash) {
	// We do not need to check as it is handled by middleware.
	course, _ := models.GetCourse(ctx.Params("course"))
	post, _ := models.GetPost(ctx.Params("post"))

	ctx.Data["Course"] = course
	ctx.Data["Post"] = post

	u, _ := models.GetUser(sess.Get("user").(string))
	if !(post.PosterID == sess.Get("user").(string) || u.IsAdmin) {
		f.Error("You may not edit this post.")
		ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
			ctx.Params("post")))
		return
	}
	if sess.Get("p.title") != nil && len(sess.Get("p.title").(string)) > 0 {
		ctx.Data["ptitle"] = sess.Get("p.title").(string)
		ctx.Data["ptext"] = sess.Get("p.text").(string)
		sess.Delete("p.title")
		sess.Delete("p.text")
	}

	ctx.Data["Title"] = "Edit post"
	ctx.Data["csrf_token"] = x.GetToken()
	ctx.HTML(200, "edit-post")
}

// PostEditPostHandler post response for editing a post.
func PostEditPostHandler(ctx *macaron.Context, x csrf.CSRF, sess session.Store, f *session.Flash) {
	post, _ := models.GetPost(ctx.Params("post"))

	u, _ := models.GetUser(sess.Get("user").(string))
	if !(post.PosterID == sess.Get("user").(string) || u.IsAdmin) {
		f.Error("You may not edit this post.")
		ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
			ctx.Params("post")))
		return
	}
	title := ctx.QueryTrim("title")
	text := ctx.QueryTrim("text")

	if !simpleTextExp.Match([]byte(title)) || len(title) > 32 || len(title) < 1 {
		f.Error("Invalid title.")
		// Pass over the text and title to errored page
		sess.Set("p.title", title)
		sess.Set("p.text", text)
		ctx.Redirect(fmt.Sprintf("/c/%s/%s/edit", ctx.Params("course"),
			ctx.Params("post")))
		return
	} else if getMarkdownLength(text) < 8 {
		f.Error("The post is empty or too short!")
		// Pass over the text and title to errored page
		sess.Set("p.title", title)
		sess.Set("p.text", text)
		ctx.Redirect(fmt.Sprintf("/c/%s/%s/edit", ctx.Params("course"),
			ctx.Params("post")))
		return
	}

	err := models.UpdatePost(&models.Post{
		PostID: post.PostID,
		Title:  title,
		Text:   text,
	})
	if err != nil {
		panic(err)
	}

	ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
		ctx.Params("post")))
}

// LitePostHandler response for a light post page.
func LitePostHandler(ctx *macaron.Context, x csrf.CSRF, sess session.Store, f *session.Flash) {
	post, _ := models.GetPost(ctx.Params("post"))
	ctx.Data["FormattedPost"] = template.HTML(markdownToHTML(post.Text))

	ctx.HTML(200, "post-lite")
}

// PostPageHandler response for a post page.
func PostPageHandler(ctx *macaron.Context, x csrf.CSRF, sess session.Store, f *session.Flash) {
	course, _ := models.GetCourse(ctx.Params("course"))
	post, _ := models.GetPost(ctx.Params("post"))

	ctx.Data["Course"] = course
	post.LoadComments()
	ctx.Data["Post"] = post

	ctx.Data["PosterID"] = strings.Split(post.PosterID, "@")[0]

	ctx.Data["FormattedPost"] = template.HTML(markdownToHTML(post.Text))
	if post.Locked || course.Locked {
		ctx.Data["Locked"] = 1
	}

	if post.Anonymous {
		ctx.Data["Poster"] = post.AnonName
	} else {
		if post.Poster.FullName == "" {
			ctx.Data["Poster"] = strings.Split(post.PosterID, "@")[0]
		} else {
			ctx.Data["Poster"] = post.Poster.FullName
		}
	}

	for i := range post.Comments {
		post.Comments[i].LoadPoster()
		post.Comments[i].FormattedText =
			template.HTML(markdownToHTML(post.Comments[i].Text))
	}
	if sess.Get("c.text") != nil {
		ctx.Data["ctext"] = sess.Get("c.text").(string)
		sess.Delete("c.text")
	}

	if sess.Get("auth") == LoggedIn {
		u, _ := models.GetUser(sess.Get("user").(string))
		if common.ContainsInt64(u.Upvoted, post.PostID) {
			ctx.Data["Upvoted"] = 1
		}
	}

	ctx.Data["csrf_token"] = x.GetToken()
	ctx.Data["Title"] = fmt.Sprintf("%s - %s", course.Code, post.Title)

	ctx.HTML(200, "post")
}

// markdownToHTML converts a string (in Markdown) and outputs (X)HTML.
// The input may also contain HTML, and the output is sanitized.
func markdownToHTML(s string) string {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(s), &buf); err != nil {
		panic(err)
	}
	return string(bluemonday.UGCPolicy().SanitizeBytes(buf.Bytes()))
}

// PostCommentPostHandler post response for posting a comment to a post.
func PostCommentPostHandler(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	postID, _ := strconv.ParseInt(ctx.Params("post"), 10, 64)

	if getMarkdownLength(ctx.QueryTrim("text")) < 2 {
		f.Error("The post is empty or too short!")
		// Pass over the text and title to errored page
		sess.Set("c.text", ctx.QueryTrim("text"))
		ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"), ctx.Params("post")))
		return
	}

	com := &models.Comment{
		PosterID: sess.Get("user").(string),
		PostID:   postID,
		Text:     ctx.QueryTrim("text"),
	}

	err := models.AddComment(com)

	if err != nil {
		panic(err)
	}

	ctx.Redirect(fmt.Sprintf("/c/%s/%s#c-%d", ctx.Params("course"),
		ctx.Params("post"), com.CommentID))
}

// CreatePostHandler response for creating a new post.
func CreatePostHandler(ctx *macaron.Context, x csrf.CSRF, sess session.Store, f *session.Flash) {
	if sess.Get("p.title") != nil && len(sess.Get("p.title").(string)) > 0 {
		ctx.Data["ptitle"] = sess.Get("p.title").(string)
		ctx.Data["ptext"] = sess.Get("p.text").(string)
		sess.Delete("p.title")
		sess.Delete("p.text")
	}

	if i, err := models.GetUserPostCount(sess.Get("user").(string)); err != nil {
		panic(err)
	} else if i > 2 {
		ctx.Data["HideNotice"] = 1 // Hide 'read the guidelines' notice
	}

	ctx.Data["csrf_token"] = x.GetToken()
	ctx.Data["Title"] = "Create post"
	ctx.HTML(200, "create-post")
}

// PostCreatePostHandler post response for creating a new post.
func PostCreatePostHandler(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	title := ctx.QueryTrim("title")
	text := ctx.Query("text")

	if !simpleTextExp.Match([]byte(title)) || len(title) > 32 || len(title) < 1 {
		f.Error("Invalid title.")
		// Pass over the text and title to errored page
		sess.Set("p.title", title)
		sess.Set("p.text", text)
		ctx.Redirect(fmt.Sprintf("/c/%s/post", ctx.Params("course")))
		return
	} else if getMarkdownLength(text) < 8 {
		f.Error("The post is empty or too short!")
		// Pass over the text and title to errored page
		sess.Set("p.title", title)
		sess.Set("p.text", text)
		ctx.Redirect(fmt.Sprintf("/c/%s/post", ctx.Params("course")))
		return
	}

	post := &models.Post{
		CourseCode: ctx.Params("course"),
		PosterID:   sess.Get("user").(string),
		Locked:     false,
		Title:      ctx.QueryTrim("title"),
		Text:       ctx.Query("text"),
		Anonymous:  false,
	}

	if ctx.Query("anon") == "on" {
		post.Anonymous = true
		post.AnonName = namegen.GetName()
	}

	err := models.AddPost(post)
	if err != nil {
		panic(err)
	}

	ctx.Redirect(fmt.Sprintf("/c/%s/%d", ctx.Params("course"), post.PostID))
}

// getMarkdownLenngth renders the string provided, removes HTML tags and
// returns the length of the trimmed final string. This is useful to determine
// the post actual length.
func getMarkdownLength(s string) int {
	markdownHTMLWithoutTags := htmlTagExp.ReplaceAll([]byte(markdownToHTML(s)), []byte(""))
	return len(strings.TrimSpace(string(markdownHTMLWithoutTags)))
}

// UpvotePostHandler post response for posting a comment to a post.
func UpvotePostHandler(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	postID, _ := strconv.ParseInt(ctx.Params("post"), 10, 64)

	u, _ := models.GetUser(sess.Get("user").(string))
	if common.ContainsInt64(u.Upvoted, postID) {
		err := models.UnvotePost(sess.Get("user").(string), postID)
		if err != nil {
			f.Error(fmt.Sprintf("%s.", err))
			ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
				ctx.Params("post")))
			return
		}
	} else {
		err := models.UpvotePost(sess.Get("user").(string), postID)
		if err != nil {
			f.Error(fmt.Sprintf("%s.", err))
			ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
				ctx.Params("post")))
			return
		}
	}

	ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
		ctx.Params("post")))
}

// DeleteCommentHandler response for deleting a comment.
func DeleteCommentHandler(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	_, err := models.GetComment(ctx.Params("id"))
	if err != nil {
		f.Error("Comment does not exist.")
		ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
			ctx.Params("post")))
		return
	}

	err = models.DeleteComment(ctx.Params("id"))
	if err != nil {
		f.Error("Failed to remove comment.")
		ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
			ctx.Params("post")))
		return
	}

	ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
		ctx.Params("post")))
}

// DeletePostHandler response for deleting a post.
func DeletePostHandler(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	err := models.DeletePost(ctx.Params("post"))
	if err != nil {
		f.Error("Failed to remove post.")
		ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
			ctx.Params("post")))
		return
	}

	ctx.Redirect(fmt.Sprintf("/c/%s", ctx.Params("course")))
}
