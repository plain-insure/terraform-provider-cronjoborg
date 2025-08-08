// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// APIError represents an error response from the API.
type APIError struct {
	StatusCode int
	Message    string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: http.DefaultClient,
	}
}

// Job represents a cron job from the API.
type Job struct {
	JobID           int         `json:"jobId"`
	Enabled         bool        `json:"enabled"`
	Title           string      `json:"title"`
	SaveResponses   bool        `json:"saveResponses"`
	URL             string      `json:"url"`
	LastStatus      int         `json:"lastStatus"`
	LastDuration    int         `json:"lastDuration"`
	LastExecution   int         `json:"lastExecution"`
	NextExecution   *int        `json:"nextExecution"` // Nullable
	Type            int         `json:"type"`
	RequestTimeout  int         `json:"requestTimeout"`
	RedirectSuccess bool        `json:"redirectSuccess"`
	FolderID        int         `json:"folderId"`
	Schedule        JobSchedule `json:"schedule"`
	RequestMethod   int         `json:"requestMethod"`
}

// DetailedJob represents a cron job with detailed settings from the API.
type DetailedJob struct {
	Job
	Auth         JobAuth                 `json:"auth"`
	Notification JobNotificationSettings `json:"notification"`
	ExtendedData JobExtendedData         `json:"extendedData"`
}

// JobAuth represents HTTP authentication settings for a job.
type JobAuth struct {
	Enable   bool   `json:"enable"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// JobNotificationSettings specifies notification settings for a job.
type JobNotificationSettings struct {
	OnFailure bool `json:"onFailure"`
	OnSuccess bool `json:"onSuccess"`
	OnDisable bool `json:"onDisable"`
}

// JobExtendedData holds extended request data for a job.
type JobExtendedData struct {
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// JobSchedule represents the execution schedule of a job.
type JobSchedule struct {
	Timezone  string `json:"timezone"`
	ExpiresAt int    `json:"expiresAt"`
	Hours     []int  `json:"hours"`
	MDays     []int  `json:"mdays"`
	Minutes   []int  `json:"minutes"`
	Months    []int  `json:"months"`
	WDays     []int  `json:"wdays"`
}

// JobHistory represents job execution history.
type JobHistory struct {
	JobLogID    int                 `json:"jobLogId"`
	JobID       int                 `json:"jobId"`
	Identifier  string              `json:"identifier"`
	Date        int                 `json:"date"`
	DatePlanned int                 `json:"datePlanned"`
	Jitter      int                 `json:"jitter"`
	URL         string              `json:"url"`
	Duration    int                 `json:"duration"`
	Status      int                 `json:"status"`
	StatusText  string              `json:"statusText"`
	HttpStatus  int                 `json:"httpStatus"`
	Headers     *string             `json:"headers"` // Nullable
	Body        *string             `json:"body"`    // Nullable
	Stats       JobHistoryItemStats `json:"stats"`
}

// JobHistoryItemStats contains additional timing information for a job execution history item.
type JobHistoryItemStats struct {
	NameLookup    int `json:"nameLookup"`
	Connect       int `json:"connect"`
	AppConnect    int `json:"appConnect"`
	PreTransfer   int `json:"preTransfer"`
	StartTransfer int `json:"startTransfer"`
	Total         int `json:"total"`
}

// JobsResponse represents the API response for listing jobs.
type JobsResponse struct {
	Jobs       []Job `json:"jobs"`
	SomeFailed bool  `json:"someFailed"`
}

// JobHistoryResponse represents the API response for job history.
type JobHistoryResponse struct {
	History     []JobHistory `json:"history"`
	Predictions []int        `json:"predictions"`
}

// JobDetailsResponse represents the API response for job details.
type JobDetailsResponse struct {
	JobDetails DetailedJob `json:"jobDetails"`
}

// doRequest performs an HTTP request to the cron-job.org API.
func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var bodyReader *bytes.Reader
	if body != nil {
		j, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(j)
	} else {
		bodyReader = bytes.NewReader([]byte{})
	}

	req, err := http.NewRequestWithContext(context.Background(), method, c.BaseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	// Check for API errors
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)

		// Try to parse error message from response
		var errorResp struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}

		errorMessage := "Unknown error"
		if json.Unmarshal(bodyBytes, &errorResp) == nil {
			if errorResp.Error != "" {
				errorMessage = errorResp.Error
			} else if errorResp.Message != "" {
				errorMessage = errorResp.Message
			}
		}

		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    errorMessage,
			Body:       string(bodyBytes),
		}
	}

	return resp, nil
}

// GetJobs retrieves all jobs from the cron-job.org API.
func (c *Client) GetJobs() ([]Job, error) {
	resp, err := c.doRequest("GET", "/jobs", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jobsResp JobsResponse
	if err := json.NewDecoder(resp.Body).Decode(&jobsResp); err != nil {
		return nil, fmt.Errorf("failed to decode jobs response: %w", err)
	}

	return jobsResp.Jobs, nil
}

// GetJob retrieves a specific job by ID from the cron-job.org API.
func (c *Client) GetJob(jobID string) (*Job, error) {
	// Use GetJobDetails and convert to Job
	detailedJob, err := c.GetJobDetails(jobID)
	if err != nil {
		return nil, err
	}

	// Convert DetailedJob to Job by extracting the embedded Job struct
	return &detailedJob.Job, nil
}

// GetJobHistory retrieves the execution history for a specific job.
func (c *Client) GetJobHistory(jobID string) ([]JobHistory, []int, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/jobs/%s/history", jobID), nil)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var historyResp JobHistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&historyResp); err != nil {
		return nil, nil, fmt.Errorf("failed to decode job history response: %w", err)
	}

	return historyResp.History, historyResp.Predictions, nil
}

// CreateJob creates a new cron job.
func (c *Client) CreateJob(job map[string]interface{}) (int, error) {
	reqBody := map[string]interface{}{
		"job": job,
	}

	resp, err := c.doRequest("PUT", "/jobs", reqBody)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		JobId int `json:"jobId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode create job response: %w", err)
	}

	return result.JobId, nil
}

// UpdateJob updates an existing cron job.
func (c *Client) UpdateJob(jobID string, job map[string]interface{}) error {
	reqBody := map[string]interface{}{
		"job": job,
	}

	_, err := c.doRequest("PATCH", fmt.Sprintf("/jobs/%s", jobID), reqBody)
	if err != nil {
		return err
	}

	return nil
}

// DeleteJob deletes a cron job.
func (c *Client) DeleteJob(jobID string) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/jobs/%s", jobID), nil)
	if err != nil {
		return err
	}

	return nil
}

// GetJobDetails retrieves detailed information for a specific job.
func (c *Client) GetJobDetails(jobID string) (*DetailedJob, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/jobs/%s", jobID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result JobDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode job details response: %w", err)
	}

	return &result.JobDetails, nil
}
