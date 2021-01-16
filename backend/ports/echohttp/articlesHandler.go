package echohttp

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/brycekbargar/realworld-backend/ports"
)

type articlesHandler struct {
	authed      echo.MiddlewareFunc
	maybeAuthed echo.MiddlewareFunc
	jc          ports.JWTConfig
}

func newArticlesHandler(
	authed echo.MiddlewareFunc,
	maybeAuthed echo.MiddlewareFunc,
	jc ports.JWTConfig) *articlesHandler {
	return &articlesHandler{
		authed,
		maybeAuthed,
		jc,
	}
}

func (r *articlesHandler) mapRoutes(g *echo.Group) {
	g.GET("/articles", r.list, r.maybeAuthed)
	g.GET("/articles/feed", r.feed, r.authed)
	g.GET("/articles/:slug", r.article, r.maybeAuthed)
}

type author struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

type articleArticle struct {
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	TagList        []string  `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	Author         author    `json:"author"`
}

type list struct {
	Articles      []articleArticle `json:"articles"`
	ArticlesCount int              `json:"articlesCount"`
}

func (r *articlesHandler) list(c echo.Context) error {
	em, _, _ := c.(*userContext).identity()

	// get all articles
	if len(em) > 0 {
		// set the following/favorited logic?
	}

	return c.JSON(http.StatusOK, list{
		make([]articleArticle, 0),
		0,
	})
}

func (r *articlesHandler) feed(c echo.Context) error {
	_, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	// Get the feed articles

	return c.JSON(http.StatusOK, list{
		make([]articleArticle, 0),
		0,
	})
}

type article struct {
	Article articleArticle `json:"article"`
}

func (r *articlesHandler) article(c echo.Context) error {
	em, _, _ := c.(*userContext).identity()

	// get the article
	if len(em) > 0 {
		// set the following/favorited logic?
	}

	return c.JSON(http.StatusOK, article{})
}