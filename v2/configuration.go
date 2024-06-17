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

	"gopkg.in/yaml.v2"
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

	globalconf, err := GetGlobalConfig(path)
	if err != nil {
		return nil, err
	}

	fmt.Println(globalconf)

	config.GlobalConfig = globalconf

	return config, nil
}

func GetGlobalConfig(path string) (GlobalConfig, error) {
	// Read YAML file
	data, err := os.ReadFile(path)
	if err != nil {
		return GlobalConfig{}, fmt.Errorf("error reading YAML file: %v", err)
	}

	// Unmarshal the YAML into the GlobalConfig struct
	var globalConfig struct {
		Global GlobalConfig `yaml:"global"`
	}
	err = yaml.Unmarshal(data, &globalConfig)
	if err != nil {
		return GlobalConfig{}, fmt.Errorf("error unmarshaling YAML: %v", err)
	}

	return globalConfig.Global, nil
}

// loginurl: Input full login URL to application (for example: https://example.com/api/login),
// email: Email for the service,
// password: Password
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
