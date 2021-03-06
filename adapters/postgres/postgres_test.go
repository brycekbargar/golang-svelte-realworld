package postgres_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/brycekbargar/realworld-backend/adapters/postgres"
	"github.com/brycekbargar/realworld-backend/adapters/testcases"
	"github.com/brycekbargar/realworld-backend/domain"
	"github.com/jackc/pgx/v4"
)

var uut domain.Repository
var dsn string

func TestMain(m *testing.M) {
	connString := "host=127.0.0.1 user=postgres password=test timezone=universal"
	testDB := fmt.Sprintf("realworld_backend_test_%v", time.Now().UnixNano())
	func() {
		db, err := pgx.Connect(context.Background(), connString)
		if err != nil {
			panic(err)
		}

		defer db.Close(context.Background())
		_, err = db.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", testDB))
		if err != nil {
			panic(err)
		}
	}()

	res := 0
	defer func() {
		if res != 0 {
			// Save the test database when tests fail
			fmt.Print(testDB)
			return
		}

		db, err := pgx.Connect(context.Background(), connString)
		if err != nil {
			panic(err)
		}
		defer db.Close(context.Background())

		_, err = db.Exec(context.Background(), fmt.Sprintf("DROP DATABASE %s", testDB))
		if err != nil {
			fmt.Print(err)
		}
	}()

	dsn = fmt.Sprintf("%s dbname=%s", connString, testDB)
	uut = postgres.
		MustNewInstance(dsn).
		MustMigrate()
	res = m.Run()
	os.Exit(res)
}

func Test_RepositoryMustMigrate(t *testing.T) {
	postgres.MustNewInstance(dsn).MustMigrate()
	postgres.MustNewInstance(dsn).MustMigrate()
	r := postgres.MustNewInstance(dsn)
	r.MustMigrate()
	postgres.MustNewInstance(dsn).MustMigrate()
	r.MustMigrate()
}

func Test_Users(t *testing.T) {
	t.Parallel()

	t.Run("Create and Update User", func(t *testing.T) {
		t.Parallel()
		testcases.Users_CreateUser(t, uut)
	})
	t.Run("Get and Update User By Email", func(t *testing.T) {
		t.Parallel()
		testcases.Users_GetUserByEmail(t, uut)
	})
	t.Run("Get and Update User By Username", func(t *testing.T) {
		t.Parallel()
		testcases.Users_GetUserByUsername(t, uut)
	})
	t.Run("Fanboy Following Users", func(t *testing.T) {
		t.Parallel()
		testcases.Users_UpdateFanboyByEmail_Following(t, uut)
	})
	t.Run("Fanboy Favoriting Articles", func(t *testing.T) {
		t.Parallel()
		testcases.Users_UpdateFanboyByEmail_Favorites(t, uut)
	})
}

func Test_Articles(t *testing.T) {
	t.Parallel()

	t.Run("Create and Update Article", func(t *testing.T) {
		t.Parallel()
		testcases.Articles_CreateArticle(t, uut)
	})
	t.Run("Get and Update Article", func(t *testing.T) {
		t.Parallel()
		testcases.Articles_GetArticleBySlug(t, uut)
	})
	t.Run("Delete Article", func(t *testing.T) {
		t.Parallel()
		testcases.Articles_DeleteArticle(t, uut)
	})
	t.Run("Query Articles", func(t *testing.T) {
		t.Parallel()
		testcases.Articles_LatestArticlesByCriteria(t, uut)
	})
	t.Run("Create and Delete Comments", func(t *testing.T) {
		t.Parallel()
		testcases.Articles_UpdateCommentsBySlug(t, uut)
	})
	t.Run("Query Tags", func(t *testing.T) {
		t.Parallel()
		testcases.Articles_DistinctTags(t, uut)
	})
}
