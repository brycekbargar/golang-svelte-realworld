package testcases

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brycekbargar/realworld-backend/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.TODO()

func Articles_CreateArticle(
	t *testing.T,
	r domain.Repository,
) {
	a := testArticle("hospitable")
	_, err := r.CreateArticle(ctx, a)
	assert.ErrorIs(t, err, domain.ErrNoAuthor)

	u := testAuthor("hospitable")
	r.CreateUser(ctx, u)

	now := time.Now().UTC()
	ca, err := r.CreateArticle(ctx, a)
	require.NoError(t, err)

	assert.Equal(t, a.Slug, ca.Slug)
	assert.Equal(t, a.Title, ca.Title)
	assert.Equal(t, a.Body, ca.Body)
	assert.Zero(t, ca.FavoriteCount)

	assert.Equal(t, a.AuthorEmail, ca.AuthorEmail)
	require.NotNil(t, ca.Author)
	assert.Equal(t, a.AuthorEmail, ca.Author.GetEmail())
	assert.Equal(t, u.Bio, ca.Author.GetBio())
	assert.Equal(t, u.Image, ca.Author.GetImage())

	assert.True(t, ca.CreatedAtUTC.After(now))
	assert.Equal(t, ca.CreatedAtUTC, ca.UpdatedAtUTC)

	_, err = r.CreateArticle(ctx, testArticle("hospitable"))
	assert.ErrorIs(t, err, domain.ErrDuplicateArticle)

	now = time.Now().UTC()
	ua := testArticle("enchanting")
	ua.SetTitle("hospitable title")
	ua.AuthorEmail = u.Email
	_, err = r.UpdateArticleBySlug(ctx,
		"hospitable-title",
		func(a *domain.Article) (*domain.Article, error) {
			return ua, nil
		})
	require.NoError(t, err)
	fa, err := r.GetArticleBySlug(ctx, "hospitable-title")
	assert.NoError(t, err)
	assert.Equal(t, ua.Slug, fa.Slug)
	assert.Equal(t, ua.Title, fa.Title)
	assert.Equal(t, ua.Body, fa.Body)

	assert.True(t, fa.CreatedAtUTC.Before(now))
	assert.True(t, fa.UpdatedAtUTC.After(now))

	r.UpdateUserByEmail(ctx,
		"author@hospitable.com",
		func(u *domain.User) (*domain.User, error) {
			u.Email = "author@whole.com"
			return u, nil
		})
	fa, err = r.GetArticleBySlug(ctx, "hospitable-title")
	require.NoError(t, err)
	assert.Equal(t, "author@whole.com", fa.AuthorEmail)
	assert.Equal(t, "author@whole.com", fa.Author.GetEmail())
}

func Articles_GetArticleBySlug(
	t *testing.T,
	r domain.Repository,
) {
	u := testAuthor("observant")
	r.CreateUser(ctx, u)

	now := time.Now().UTC()
	a := testArticle("observant")
	_, err := r.CreateArticle(ctx, a)
	require.NoError(t, err)

	_, err = r.GetArticleBySlug(ctx, "befitting-title")
	assert.ErrorIs(t, err, domain.ErrArticleNotFound)

	fa, err := r.GetArticleBySlug(ctx, a.Slug)
	assert.NoError(t, err)
	assert.Equal(t, a.Slug, fa.Slug)
	assert.Equal(t, a.Title, fa.Title)
	assert.Equal(t, a.Body, fa.Body)

	assert.Equal(t, a.AuthorEmail, fa.AuthorEmail)
	require.NotNil(t, fa.Author)
	assert.Equal(t, a.AuthorEmail, fa.Author.GetEmail())
	assert.Equal(t, u.Bio, fa.Author.GetBio())
	assert.Equal(t, u.Image, fa.Author.GetImage())

	assert.True(t, fa.CreatedAtUTC.After(now))
	assert.Equal(t, fa.CreatedAtUTC, fa.UpdatedAtUTC)

	r.UpdateArticleBySlug(ctx,
		"observant-title",
		func(a *domain.Article) (*domain.Article, error) {
			a.SetTitle("silent title")
			return a, nil
		})
	_, err = r.GetArticleBySlug(ctx, "observant-title")
	assert.ErrorIs(t, err, domain.ErrArticleNotFound)
	_, err = r.GetArticleBySlug(ctx, "silent-title")
	assert.NoError(t, err)

	r.CreateUser(ctx, testAuthor("modern"))
	r.CreateArticle(ctx, testArticle("modern"))
	_, err = r.UpdateArticleBySlug(ctx,
		"silent-title",
		func(a *domain.Article) (*domain.Article, error) {
			a.SetTitle("modern title")
			return a, nil
		})
	assert.ErrorIs(t, err, domain.ErrDuplicateArticle)
	fa, err = r.GetArticleBySlug(ctx, "silent-title")
	assert.NoError(t, err)
}
func Articles_DeleteArticle(
	t *testing.T,
	r domain.Repository,
) {
	r.CreateUser(ctx, testAuthor("deranged"))
	r.CreateArticle(ctx, testArticle("deranged"))

	a, err := r.GetArticleBySlug(ctx, "deranged-title")
	require.NoError(t, err)
	err = r.DeleteArticle(ctx, &a.Article)
	assert.NoError(t, err)
	a, err = r.GetArticleBySlug(ctx, "deranged-title")
	assert.ErrorIs(t, err, domain.ErrArticleNotFound)

	err = r.DeleteArticle(ctx, nil)
	assert.NoError(t, err)
}

func Articles_LatestArticlesByCriteria(
	t *testing.T,
	r domain.Repository,
) {
	tt := "Articles_LatestArticlesByCriteria"

	authors := []domain.User{
		*testAuthor("frail"),
		*testAuthor("simple"),
		*testAuthor("reminiscent"),
	}
	for _, a := range authors {
		_, err := r.CreateUser(ctx, &a)
		require.NoError(t, err)
	}

	source := make([]domain.AuthoredArticle, 0, 13)
	for i, adj := range []string{
		"bright",
		"colorful",
		"sour",
		"fantastic",
		"mellow",
		"splendid",
		"gruesome",
		"madly",
		"kind",
		"organic",
		"public",
		"flagrant",
		"waggish",
	} {
		a := testArticle(adj)
		a.AuthorEmail = authors[i%3].Email
		a.TagList = append(a.TagList, tt)

		_, err := r.CreateArticle(ctx, a)
		require.NoError(t, err)

		fa, err := r.GetArticleBySlug(ctx, a.Slug)
		require.NoError(t, err)
		source = append(source, *fa)
	}

	cases := []struct {
		Name  string
		Query domain.ListCriteria
		Try   func(*testing.T) func([]domain.AuthoredArticle, error)
	}{
		{
			"All Tagged",
			domain.ListCriteria{
				Tag:    tt,
				Offset: 0,
				Limit:  13,
			},
			func(t *testing.T) func([]domain.AuthoredArticle, error) {
				return func(all []domain.AuthoredArticle, err error) {
					require.NoError(t, err)

					assert.ElementsMatch(t, all, source)
				}
			},
		},
		{
			"Offset Limit Sort",
			domain.ListCriteria{
				Tag:    tt,
				Offset: 2,
				Limit:  5,
			},
			func(t *testing.T) func([]domain.AuthoredArticle, error) {
				return func(some []domain.AuthoredArticle, err error) {
					require.NoError(t, err)

					require.Len(t, some, 5)
					assert.Equal(t, "public-title", some[0].Slug)
					assert.Equal(t, "gruesome-title", some[4].Slug)
				}
			},
		},
		{
			"Authored By",
			domain.ListCriteria{
				Tag:          tt,
				Offset:       0,
				Limit:        13,
				AuthorEmails: []string{"author@frail.com", "author@simple.com"},
			},
			func(t *testing.T) func([]domain.AuthoredArticle, error) {
				return func(authored []domain.AuthoredArticle, err error) {
					require.NoError(t, err)

					assert.NotEmpty(t, authored)
					for _, a := range authored {
						assert.NotEqual(t, "author@reminiscent.com", a.AuthorEmail)
					}
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			tc.Try(t)(r.LatestArticlesByCriteria(ctx, tc.Query))
		})
	}
}

func Articles_UpdateCommentsBySlug(
	t *testing.T,
	r domain.Repository,
) {
	r.CreateUser(ctx, testUser("simplistic"))
	r.CreateUser(ctx, testAuthor("envious"))

	now := time.Now().UTC()
	for _, adj := range []string{
		"tedious",
		"polite",
		"divergent",
	} {
		a := testArticle(adj)
		a.AuthorEmail = "author@envious.com"

		_, err := r.CreateArticle(ctx, a)
		require.NoError(t, err)
	}

	c, err := r.UpdateCommentsBySlug(ctx,
		"tedious-title",
		func(a *domain.CommentedArticle) (*domain.CommentedArticle, error) {
			err := a.AddComment("enchanting body", "user@simplistic.com")
			if err != nil {
				return nil, err
			}
			return a, nil
		})
	require.NoError(t, err)

	assert.Positive(t, c.ID)
	assert.Equal(t, "enchanting body", c.Body)
	assert.Equal(t, "user@simplistic.com", c.AuthorEmail)
	assert.True(t, now.Before(c.CreatedAtUTC))

	_, err = r.UpdateCommentsBySlug(ctx,
		"tedious-title",
		func(a *domain.CommentedArticle) (*domain.CommentedArticle, error) {
			err := a.AddComment("quirky body", "user@simplistic.com")
			if err != nil {
				return nil, err
			}
			return a, nil
		})
	require.NoError(t, err)

	a, err := r.GetCommentsBySlug(ctx, "tedious-title")
	require.NoError(t, err)

	assert.Len(t, a.Comments, 2)

	_, err = r.UpdateCommentsBySlug(ctx,
		"tedious-title",
		func(a *domain.CommentedArticle) (*domain.CommentedArticle, error) {
			a.RemoveComment(c.ID)
			return a, nil
		})
	require.NoError(t, err)

	a, err = r.GetCommentsBySlug(ctx, "tedious-title")
	require.NoError(t, err)

	require.Len(t, a.Comments, 1)
	assert.Positive(t, a.Comments[0].ID)
	assert.Equal(t, "quirky body", a.Comments[0].Body)
	assert.Equal(t, "user@simplistic.com", a.Comments[0].AuthorEmail)
	assert.True(t, now.Before(a.Comments[0].CreatedAtUTC))
}

func Articles_DistinctTags(
	t *testing.T,
	r domain.Repository,
) {
	tt := "Articles_DistinctTags"

	r.CreateUser(ctx, testAuthor("threatening"))
	for _, adj := range []string{
		"stimulating",
		"exultant",
		"helpless",
	} {
		a := testArticle(adj)
		a.AuthorEmail = "author@threatening.com"
		a.TagList = append(a.TagList, tt)

		_, err := r.CreateArticle(ctx, a)
		require.NoError(t, err)
	}

	tags, err := r.DistinctTags(ctx)
	require.NoError(t, err)

	assert.Contains(t, tags, "stimulating one", "stimulating two")
	assert.Contains(t, tags, "stimulating two")
	assert.Contains(t, tags, "helpless three")
	htt := false
	for i := range tags {
		if tags[i] == tt {
			if htt {
				assert.Fail(t, "the testing tag was not distinct")
			}
			htt = true
		}
	}
}

func testAuthor(adj string) *domain.User {
	a := testUser(adj)
	a.Email = fmt.Sprintf("author@%v.com", adj)
	return a
}

func testArticle(adj string) (a *domain.Article) {
	a, _ = domain.NewArticle(
		fmt.Sprintf("%v title", adj),
		fmt.Sprintf("%v description", adj),
		fmt.Sprintf("%v body", adj),
		fmt.Sprintf("author@%v.com", adj),
		fmt.Sprintf("%v one", adj),
		fmt.Sprintf("%v two", adj),
		fmt.Sprintf("%v three", adj),
	)
	return
}
