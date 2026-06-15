// Package client gerencia token e client para reqs no google
package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"lazy-queues/util"
)

const (
	GoogleMonitoringBaseURL string = "https://monitoring.googleapis.com/v3/"
	GooglePubSubBaseURL     string = "https://pubsub.googleapis.com/v1/"
	ProjectName             string = "mercado-libre-prd"
)

type Client struct {
	client *http.Client
	token  string
}

var MainClient = Client{}

func getGoogleAuthToken() (token string, err error) {
	authToken, err := exec.Command("gcloud", "auth", "application-default", "print-access-token").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(authToken)), nil
}

func prepareURL(bareURL string) string {
	return GoogleMonitoringBaseURL + bareURL
}

func FetchData[T any](url string) (T, error) {
	var zero T
	// finalEndpoint := prepareURL(url)
	util.Log.Info("Creating request to endpoint", "url", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		util.Log.Error("Error Creating request", "error", err)
		return zero, err
	}

	req.Header.Set("Authorization", "Bearer "+MainClient.token)
	req.Header.Set("Accept", "application/json")

	response, err := MainClient.client.Do(req)
	if err != nil {
		util.Log.Error("Error fetching request", "error", err)
		return zero, err
	}
	defer response.Body.Close()

	fmt.Println(response.StatusCode)

	var result T

	if response.StatusCode != 200 {
		util.Log.Error("Not successfull response", "status", response.StatusCode)
		return zero, errors.New("Status code different from 200")
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		util.Log.Error("Error decoding response", "error", err)
		return zero, err
	}

	return result, nil
}

func Execute() {
	token, err := getGoogleAuthToken()
	if err != nil {
		fmt.Println("Error getting auth token from gcloud cli", err)
	}

	MainClient = Client{
		client: &http.Client{},
		token:  token,
	}
}
