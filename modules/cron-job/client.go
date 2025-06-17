package cronjob

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type APIClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(config *ProviderConfig) *APIClient {
	return &APIClient{
		BaseURL:    config.APIUrl,
		APIKey:     config.APIKey,
		HTTPClient: http.DefaultClient,
	}
}

func (c *APIClient) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		j, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(j)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, respBody)
	}
	return resp, nil
}