// +build acceptance

package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type NamedTagList struct {
	Id   string
	Name string
	Tags []string
}

type Build struct {
	Sha1    string
	Version string
}

func getNamedTagLists(baseUrl string, buckets []string) ([]NamedTagList, error) {
	var (
		err      error
		response *http.Response
	)

	if response, err = http.Get(fmt.Sprintf("%s/namedTagLists?%s", baseUrl, queryString(buckets))); err != nil {
		return nil, err
	}

	if err := assertStatusCode(response, 200); err != nil {
		return nil, err
	}

	var namedTagLists []NamedTagList
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&namedTagLists)
	return namedTagLists, err
}

func createNamedTagList(baseUrl string, buckets []string, namedTagList NamedTagList) (*NamedTagList, error) {
	var (
		err         error
		requestBody []byte
		response    *http.Response
	)

	if requestBody, err = json.Marshal(namedTagList); err != nil {
		return nil, err
	}

	if response, err = http.Post(
		fmt.Sprintf("%s/namedTagLists?%s", baseUrl, queryString(buckets)),
		"application/json",
		bytes.NewBuffer(requestBody),
	); err != nil {
		return nil, err
	}

	if err := assertStatusCode(response, 201); err != nil {
		return nil, err
	}

	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&namedTagList)
	return &namedTagList, err
}

func replaceNamedTagList(baseUrl string, id string, namedTagList NamedTagList) error {
	var (
		err         error
		request     *http.Request
		requestBody []byte
		response    *http.Response
	)

	if requestBody, err = json.Marshal(namedTagList); err != nil {
		return err
	}

	if request, err = http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/namedTagLists?id=%s", baseUrl, id),
		bytes.NewReader(requestBody),
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

func deleteNamedTagList(baseUrl string, id string) error {
	var (
		err      error
		request  *http.Request
		response *http.Response
	)

	if request, err = http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/namedTagLists?id=%s", baseUrl, id),
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

func deleteNamedTagLists(baseUrl string) error {
	var (
		err      error
		request  *http.Request
		response *http.Response
	)

	if request, err = http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/namedTagLists", baseUrl),
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

func queryString(buckets []string) string {
	query := []string{}
	for _, bucket := range buckets {
		query = append(query, fmt.Sprintf("bucket=%s", bucket))
	}
	return strings.Join(query, "&")
}

func assertStatusCode(response *http.Response, wantStatusCode int) error {
	gotStatusCode := response.StatusCode
	if gotStatusCode != wantStatusCode {
		var responseBody map[string]string
		defer response.Body.Close()
		if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return fmt.Errorf("got status-code=%d want status-code=%d", gotStatusCode, wantStatusCode)
		} else {
			return fmt.Errorf("got status-code=%d want status-code=%d (response-body=%+v)", gotStatusCode, wantStatusCode, responseBody)
		}
	}
	return nil
}
