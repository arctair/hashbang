package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"

	v1 "github.com/arctair/hashbang/v1"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	sha1    string
	version string
)

// StartHTTPServer ...
func StartHTTPServer(wg *sync.WaitGroup) *http.Server {
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	if err = v1.Migrate(pool); err != nil {
		panic(err)
	}

	namedTagListRepository := v1.NewNamedTagListRepository(
		pool,
	)

	server := &http.Server{
		Addr: ":5000",
		Handler: v1.NewRouter(
			v1.NewNamedTagListController(
				v1.NewLogger(),
				namedTagListRepository,
				v1.NewNamedTagListService(
					namedTagListRepository,
					v1.NewUUIDGenerator(),
				),
			),
			v1.NewVersionController(
				v1.NewBuild(sha1, version),
			),
		),
	}

	go func() {
		defer wg.Done()

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	return server
}

func main() {
	serverExit := &sync.WaitGroup{}
	serverExit.Add(1)
	StartHTTPServer(serverExit)
	serverExit.Wait()
}
