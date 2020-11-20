// +build acceptance

package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

func getNamedTagLists(baseUrl string) ([]NamedTagList, error) {
	var (
		err      error
		response *http.Response
	)

	if response, err = http.Get(fmt.Sprintf("%s/namedTagLists", baseUrl)); err != nil {
		return nil, err
	}

	gotStatusCode := response.StatusCode
	wantStatusCode := 200

	if gotStatusCode != wantStatusCode {
		return nil, fmt.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
	}

	var namedTagLists []NamedTagList
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&namedTagLists)
	return namedTagLists, err
}

func createNamedTagList(baseUrl string, bucket string, namedTagList NamedTagList) (*NamedTagList, error) {
	var (
		err         error
		url         string
		requestBody []byte
		response    *http.Response
	)

	if requestBody, err = json.Marshal(namedTagList); err != nil {
		return nil, err
	}

	if len(bucket) > 0 {
		url = fmt.Sprintf("%s/namedTagLists?bucket=%s", baseUrl, bucket)
	} else {
		url = fmt.Sprintf("%s/namedTagLists", baseUrl)
	}

	if response, err = http.Post(
		url,
		"application/json",
		bytes.NewBuffer(requestBody),
	); err != nil {
		return nil, err
	}

	gotStatusCode := response.StatusCode
	wantStatusCode := 201

	if gotStatusCode != wantStatusCode {
		var responseBody map[string]string
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&responseBody)
		if err != nil {
			return nil, fmt.Errorf("got status-code=%d want status-code=%d", gotStatusCode, wantStatusCode)
		} else {
			return nil, fmt.Errorf("got status-code=%d response-body=%+v want status-code=%d", gotStatusCode, responseBody, wantStatusCode)
		}
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
