// +build acceptance

package main_test

import (
	"bufio"
	"bytes"
	"encoding/json"
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
)

type Post struct {
	ImageUri string
	Tags     []string
}

func getPosts(baseUrl string) ([]Post, error) {
	var (
		err      error
		response *http.Response
	)

	if response, err = http.Get(fmt.Sprintf("%s/posts", baseUrl)); err != nil {
		return nil, err
	}

	gotStatusCode := response.StatusCode
	wantStatusCode := 200

	if gotStatusCode != wantStatusCode {
		return nil, fmt.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
	}

	var posts []Post
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&posts)
	return posts, err
}

func createPost(baseUrl string, post Post) error {
	var (
		err         error
		requestBody []byte
		response    *http.Response
	)

	if requestBody, err = json.Marshal(post); err != nil {
		return err
	}

	if response, err = http.Post(
		fmt.Sprintf("%s/posts", baseUrl),
		"application/json",
		bytes.NewBuffer(requestBody),
	); err != nil {
		return err
	}

	gotStatusCode := response.StatusCode
	wantStatusCode := 201

	if gotStatusCode != wantStatusCode {
		return fmt.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
	}
	return nil
}

func deletePosts(baseUrl string) error {
	var (
		err      error
		request  *http.Request
		response *http.Response
	)

	if request, err = http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/posts", baseUrl),
		nil,
	); err != nil {
		return err
	}

	if response, err = http.DefaultClient.Do(request); err != nil {
		return err
	}

	gotStatusCode := response.StatusCode
	wantStatusCode := 204

	if gotStatusCode != wantStatusCode {
		return fmt.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
	}
	return nil
}

func TestAcceptance(t *testing.T) {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		baseUrl = "http://localhost:5000"

		assertutil.NotError(t, exec.Command("sh", "build").Run())

		command := exec.Command("bin/hashbang")
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
		response, err := http.Get(fmt.Sprintf("%s/version", baseUrl))
		assertutil.NotError(t, err)

		gotStatusCode := response.StatusCode
		wantStatusCode := 200

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		var gotBody map[string]string
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&gotBody)
		assertutil.NotError(t, err)

		sha1Pattern := regexp.MustCompile("^[0-9a-f]{40}(-dirty)?$")
		versionPattern := regexp.MustCompile("^\\d+\\.\\d+\\.\\d+$")

		if !sha1Pattern.MatchString(gotBody["sha1"]) {
			t.Errorf("got sha1 %s want 40 hex digits", gotBody["sha1"])
		}
		if !versionPattern.MatchString(gotBody["version"]) && !sha1Pattern.MatchString(gotBody["version"]) {
			t.Errorf("got version %s want semver or 40 hex digits", gotBody["version"])
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
