package client

import (
	"os/exec"
	"fmt"
	"strings"
	"net/http"
	"encoding/json"
)

const GoogleMonitoringBaseURL string = "https://monitoring.googleapis.com/v3/"
const ProjectName = "mercado-libre-prd"

type Client struct {
	client *http.Client
	token string
}

var MainClient = Client{}

func getGoogleAuthToken() (token string, err error){
	authToken, err := exec.Command("gcloud", "auth", "application-default", "print-access-token").Output()
	if err != nil { 
		return "", err
	}

	return strings.TrimSpace(string(authToken)), nil
}

func prepareUrl(bareUrl string) string {
	return GoogleMonitoringBaseURL + bareUrl
}

//get NewRequest
func FetchData[T any](url string) (T, error){
	var zero T
	finalEndpoint := prepareUrl(url)

	req, err := http.NewRequest("GET", finalEndpoint, nil)
	if err != nil {
		return zero, err 
	}

	req.Header.Set("Authorization", "Bearer " + MainClient.token)

	response, err := MainClient.client.Do(req)
	if err != nil {
		return zero, err 
	}
	defer response.Body.Close()

	var result T
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return zero, err

	}

	return result, nil
}

func Execute() {
	token, err := getGoogleAuthToken()
	if(err != nil) {
		fmt.Println("Error getting auth token from gcloud cli", err)
	}

	MainClient = Client {
		client: &http.Client{},
		token: token,
	} 
}
