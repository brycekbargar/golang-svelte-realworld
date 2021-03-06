package domain

import (
	"context"
	"errors"
)

// ErrUserNotFound indicates the requested user was not found.
var ErrUserNotFound = errors.New("user not found")

// ErrDuplicateUser indicates the requested user could not be created because they already exist.
var ErrDuplicateUser = errors.New("user has a duplicate username or email address")

// ErrNoAuthor indicates when the author of an Article can't be found.
var ErrNoAuthor = errors.New("author not found")

// ErrArticleNotFound indicates the requested article was not found.
var ErrArticleNotFound = errors.New("article not found")

// ErrDuplicateArticle indicates the requested article could not be created because another article has the same slug.
var ErrDuplicateArticle = errors.New("article has a duplicate slug")

// ListCriteria is the set of optional parameters to page/filter the Articles.
type ListCriteria struct {
	Tag                  string
	AuthorEmails         []string
	FavoritedByUserEmail string
	Limit                int
	Offset               int
}

// Repository allows performing abstracted I/O operations on users.
type Repository interface {
	// CreateUser creates a new user.
	CreateUser(context.Context, *User) (*User, error)
	// GetUserByEmail finds a single user based on their email address.
	GetUserByEmail(context.Context, string) (*Fanboy, error)
	// GetAuthorByEmail finds a single author based on their email address or nil if they don't exist.
	GetAuthorByEmail(context.Context, string) Author
	// GetUserByUsername finds a single user based on their username.
	GetUserByUsername(context.Context, string) (*User, error)
	// UpdateUserByEmail finds a single user based on their email address,
	// then applies the provide mutations.
	UpdateUserByEmail(context.Context, string, func(*User) (*User, error)) (*User, error)
	// UpdateFanboyByEmail finds a single user based on their email address,
	// then applies the provide mutations (probably to the follower list).
	UpdateFanboyByEmail(context.Context, string, func(*Fanboy) (*Fanboy, error)) error

	// CreateArticle creates a new article.
	CreateArticle(context.Context, *Article) (*AuthoredArticle, error)
	// LatestArticlesByCriteria lists articles paged/filtered by the given criteria.
	LatestArticlesByCriteria(context.Context, ListCriteria) ([]AuthoredArticle, error)
	// GetArticleBySlug gets a single article with the given slug.
	GetArticleBySlug(context.Context, string) (*AuthoredArticle, error)
	// GetCommentsBySlug gets a single article and its comments with the given slug.
	GetCommentsBySlug(context.Context, string) (*CommentedArticle, error)
	// UpdateArticleBySlug finds a single article based on its slug
	// then applies the provide mutations.
	UpdateArticleBySlug(context.Context, string, func(*Article) (*Article, error)) (*AuthoredArticle, error)
	// UpdateCommentsBySlug finds a single article based on its slug
	// then applies the provide mutations to its comments.
	UpdateCommentsBySlug(context.Context, string, func(*CommentedArticle) (*CommentedArticle, error)) (*Comment, error)
	// DeleteArticle deletes the article if it exists.
	DeleteArticle(context.Context, *Article) error
	// DistinctTags returns a distinct list of tags on all articles
	DistinctTags(context.Context) ([]string, error)
}
