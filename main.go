package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type Post struct {
	ID    string
	Title string
	Body  string
}

func (p Post) Slug() string {
	slug := fmt.Sprintf("%s-%s", p.Title, p.ID)
	slug = strings.ToLower(slug)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}

var posts = []Post{
	{
		ID:    "1",
		Title: "Hello World",
		Body:  "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam",
	},
	{
		ID:    "2",
		Title: "Foo Bar",
		Body:  "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam",
	},
	{
		ID:    "3",
		Title: "Sticks and Stones",
		Body:  "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam",
	},
}

func GetPostById(id string) (*Post, error) {
	for _, post := range posts {
		if post.ID == id {
			return &post, nil
		}
	}

	return nil, fmt.Errorf("post does not exist")
}

func GetPostId(slug string) string {
	parts := strings.Split(slug, "-")
	return parts[len(parts)-1]
}

func ShowPost(ctx echo.Context) error {
	slug := ctx.Param("slug")

	postId := GetPostId(slug)

	post, err := GetPostById(postId)
	if err != nil {
		return ctx.String(http.StatusNotFound, err.Error())
	}

	if post.Slug() != slug {
		return ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/%s", post.Slug()))
	}

	return ctx.String(http.StatusOK, fmt.Sprintf("%s - %s - %s", post.ID, post.Title, post.Body))
}

func ListPosts(ctx echo.Context) error {
	urls := make([]string, len(posts))

	for _, post := range posts {
		url := fmt.Sprintf("http://localhost:8080/%s", post.Slug())
		urls = append(urls, fmt.Sprintf("<a href='%s'>%s</a>", url, url))
	}

	return ctx.HTML(http.StatusOK, strings.Join(urls, "<br>"))
}

func main() {
	e := echo.New()

	e.GET("/", ListPosts)
	e.GET("/:slug", ShowPost)

	e.Logger.Fatal(e.Start(":8080"))
}
