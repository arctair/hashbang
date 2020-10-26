// +build acceptance

package main_test

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"testing"
	"time"

	assertutil "github.com/arctair/go-assertutil"
	"github.com/cenkalti/backoff/v4"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
)

func TestAcceptance(t *testing.T) {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		baseUrl = "http://localhost:5000"

		assertutil.NotError(t, exec.Command("sh", "build").Run())

		testServer, err := testserver.NewTestServer()
		defer testServer.Stop()
		assertutil.NotError(t, err)

		command := exec.Command("bin/hashbang")
		command.Env = append(command.Env, fmt.Sprintf("DATABASE_URL=%s", testServer.PGURL().String()))
		stderr, err := command.StderrPipe()
		assertutil.NotError(t, err)
		assertutil.NotError(t, command.Start())
		defer dumpPipe("app:", stderr)
		defer command.Process.Kill()

		assertutil.NotError(
			t,
			backoff.Retry(
				func() error {
					_, err := http.Get(baseUrl)
					return err
				},
				NewExponentialBackOff(3*time.Second),
			),
		)
	}

	t.Run("create and get posts", func(t *testing.T) {
		// get posts is empty
		gotPosts, err := getPosts(baseUrl)
		assertutil.NotError(t, err)
		wantPosts := []Post{}

		if !reflect.DeepEqual(gotPosts, wantPosts) {
			t.Errorf("got posts %+v want %+v", gotPosts, wantPosts)
		}

		// create post
		if err = createPost(
			baseUrl,
			Post{
				ImageUri: "https://images.unsplash.com/photo-1603316851229-26637b4bd1b8?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1400&q=80",
				Tags: []string{
					"#windy",
					"#tdd",
				},
			},
		); err != nil {
			t.Fatal(err)
		}

		// get posts is not empty
		gotPosts, err = getPosts(baseUrl)
		assertutil.NotError(t, err)

		wantPosts = []Post{
			{
				ImageUri: "https://images.unsplash.com/photo-1603316851229-26637b4bd1b8?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1400&q=80",
				Tags: []string{
					"#windy",
					"#tdd",
				},
			},
		}

		if !reflect.DeepEqual(gotPosts, wantPosts) {
			t.Errorf("got posts %+v want %+v", gotPosts, wantPosts)
		}

		// delete posts
		err = deletePosts(baseUrl)
		assertutil.NotError(t, err)

		// get posts is empty
		gotPosts, err = getPosts(baseUrl)
		assertutil.NotError(t, err)
		wantPosts = []Post{}

		if !reflect.DeepEqual(gotPosts, wantPosts) {
			t.Errorf("got posts %+v want %+v", gotPosts, wantPosts)
		}
	})

	t.Run("GET /version returns sha1 and version", func(t *testing.T) {
		build, err := getVersion(baseUrl)
		assertutil.NotError(t, err)

		sha1Pattern := regexp.MustCompile("^[0-9a-f]{40}(-dirty)?$")
		versionPattern := regexp.MustCompile("^\\d+\\.\\d+\\.\\d+$")

		if !sha1Pattern.MatchString(build.Sha1) {
			t.Errorf("got sha1 %s want 40 hex digits", build.Sha1)
		}
		if !versionPattern.MatchString(build.Version) && !sha1Pattern.MatchString(build.Version) {
			t.Errorf("got version %s want semver or 40 hex digits", build.Version)
		}
	})
}

func dumpPipe(prefix string, p io.ReadCloser) {
	s := bufio.NewScanner(p)
	for s.Scan() {
		log.Printf("%s: %s", prefix, s.Text())
	}
	if err := s.Err(); err != nil {
		log.Printf("Failed to dump pipe: %s", err)
	}
}

func NewExponentialBackOff(timeout time.Duration) *backoff.ExponentialBackOff {
	b := &backoff.ExponentialBackOff{
		InitialInterval:     backoff.DefaultInitialInterval,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         backoff.DefaultMaxInterval,
		MaxElapsedTime:      timeout,
		Stop:                backoff.Stop,
		Clock:               backoff.SystemClock,
	}
	b.Reset()
	return b
}
