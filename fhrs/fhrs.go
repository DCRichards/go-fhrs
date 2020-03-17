package fhrs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	endpoint = "https://api.ratings.food.gov.uk/"
	version  = 2
)

// APILanguage represents the language API responses will be returned in.
type APILanguage int

const (
	English APILanguage = iota // English (en-GB)
	Cymraeg                    // Welsh (cy-GB)
)

func (l APILanguage) String() string {
	return []string{"en-GB", "cy-GB"}[l]
}

// APIError encapsulated a general error coming from an API request. This is for
// the cases which do not have specific errors.
type APIError struct {
	Method     string
	URL        string
	StatusCode int
	Message    string
}

func (e APIError) Error() string {
	return fmt.Sprintf(
		"API Error: %s %s returned status %d. %s",
		e.Method, e.URL, e.StatusCode, e.Message,
	)
}

// Timestamp is a representation of the date/time format used throughout the API.
//
// It is mostly RFC3339 with the timezone omitted, but this is not consistent, so
// both variations can be parsed as JSON.
type Timestamp time.Time

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" || s == "undefined" {
		*t = Timestamp(time.Time{})
		return nil
	}

	parsed, err := time.Parse("2006-01-02T15:04:05", s)
	if err == nil {
		*t = Timestamp(parsed)
		return nil
	}

	// If we can't pase as the above, try RFC3339.
	parsed, err = time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}

	*t = Timestamp(parsed)
	return nil
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

func (t Timestamp) String() string {
	ts := time.Time(t)
	return ts.String()
}

type ErrorResponse struct {
	Message string `json:"Message"`
}

// Client provides the entry point to all of the available services.
type Client struct {
	httpClient *http.Client
	language   APILanguage
	baseURL    *url.URL
	version    int
	common     service // Reuse this for all services.

	Establishments *EstablishmentsService
}

type service struct {
	client *Client
}

// NewClient creates a new FHRS Client.
func NewClient() (*Client, error) {
	httpClient := &http.Client{Timeout: 15 * time.Second}

	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	client := &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		version:    version,
	}

	client.common.client = client
	client.Establishments = (*EstablishmentsService)(&client.common)

	return client, nil
}

// SetLanguage sets the response language.
func (c *Client) SetLanguage(l APILanguage) error {
	switch l {
	case English, Cymraeg:
		c.language = l
		return nil
	}

	return errors.New("Language not supported")
}

func (c *Client) get(url string, responseBody interface{}) error {
	u, err := c.baseURL.Parse(url)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("x-api-version", strconv.Itoa(c.version))
	req.Header.Set("Accept-Language", c.language.String())

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	switch {
	// 404: Simply return nil to denote nothing to return.
	case res.StatusCode == http.StatusNotFound:
		return nil
	// Otherwise parse and return general API error.
	case res.StatusCode < 200 || res.StatusCode >= 300:
		var errorResponse ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&errorResponse); err != nil {
			if err != io.EOF {
				return err
			}
		}

		return APIError{
			Method:     req.Method,
			URL:        req.URL.String(),
			StatusCode: res.StatusCode,
			Message:    errorResponse.Message,
		}
	}

	if err := json.NewDecoder(res.Body).Decode(responseBody); err != nil {
		if err == io.EOF {
			return nil
		}

		return err
	}

	return nil
}
