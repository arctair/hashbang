// +build acceptance

package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Post struct {
	ImageUri string
	Tags     []string
}

type Build struct {
	Sha1    string
	Version string
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

func getVersion(baseUrl string) (*Build, error) {
	var (
		response *http.Response
		err      error
	)

	if response, err = http.Get(fmt.Sprintf("%s/version", baseUrl)); err != nil {
		return nil, err
	}

	gotStatusCode := response.StatusCode
	wantStatusCode := 200

	if gotStatusCode != wantStatusCode {
		return nil, fmt.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
	}

	var build Build
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&build)
	return &build, err
}
