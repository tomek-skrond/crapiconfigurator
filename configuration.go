package crapiconfigurator

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func GetConfig(path string) (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config *Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config, nil
}

func GetJWTToken(loginurl, email, password string) string {
	client := CustomHttpClient()

	jsonData, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	if err != nil {
		log.Fatalln(err)
	}

	req, _ := http.NewRequest("POST", loginurl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Send a POST request to the login endpoint
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer resp.Body.Close()
	// resp, err := http.Post(loginurl, "application/json", bytes.NewBuffer(jsonData))

	// Decode the response body into TokenResponse
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	return tokenResp.Token
}

func ConfigureRequest(req *http.Request, token string) *http.Request {
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "asdfsd")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	return req
}

func CustomHttpClient() *http.Client {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	return client
}

func ReadBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return body, nil
}
