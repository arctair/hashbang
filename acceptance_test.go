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
		stdout, err := command.StdoutPipe()
		assertutil.NotError(t, err)
		stderr, err := command.StderrPipe()
		assertutil.NotError(t, err)
		assertutil.NotError(t, command.Start())
		defer dumpPipe("appout", stdout)
		defer dumpPipe("apperr", stderr)
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

	t.Run("life of a named tag list", func(t *testing.T) {
		// get named tag lists is empty
		gotNamedTagLists, err := getNamedTagLists(baseUrl)
		assertutil.NotError(t, err)
		wantNamedTagLists := []NamedTagList{}

		if !reflect.DeepEqual(gotNamedTagLists, wantNamedTagLists) {
			t.Errorf("got named tag lists %+v want %+v", gotNamedTagLists, wantNamedTagLists)
		}

		// create named tag list
		if err = createNamedTagList(
			baseUrl,
			NamedTagList{
				Name: "named tag list",
				Tags: []string{
					"#windy",
					"#tdd",
				},
			},
		); err != nil {
			t.Fatal(err)
		}

		// get named tag lists is not empty
		gotNamedTagLists, err = getNamedTagLists(baseUrl)
		assertutil.NotError(t, err)

		wantNamedTagLists = []NamedTagList{
			{
				Name: "named tag list",
				Tags: []string{
					"#windy",
					"#tdd",
				},
			},
		}

		if !reflect.DeepEqual(gotNamedTagLists, wantNamedTagLists) {
			t.Errorf("got named tag lists %+v want %+v", gotNamedTagLists, wantNamedTagLists)
		}

		// delete named tag lists
		err = deleteNamedTagLists(baseUrl)
		assertutil.NotError(t, err)

		// get named tag lists is empty
		gotNamedTagLists, err = getNamedTagLists(baseUrl)
		assertutil.NotError(t, err)
		wantNamedTagLists = []NamedTagList{}

		if !reflect.DeepEqual(gotNamedTagLists, wantNamedTagLists) {
			t.Errorf("got named tag lists %+v want %+v", gotNamedTagLists, wantNamedTagLists)
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
