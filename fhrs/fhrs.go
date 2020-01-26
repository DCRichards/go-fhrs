package fhrs

import (
	"errors"
	"net/http"
	"net/url"
	"time"
)

const (
	endpoint = "https://api.ratings.food.gov.uk/"
	version  = 2
)

type APILanguage int

const (
	English APILanguage = iota
	Cymraeg
)

func (l APILanguage) String() string {
	return []string{"en-GB", "cy-GB"}[l]
}

// Client provides the entry point to all of the available services.
type Client struct {
	httpClient *http.Client
	language   APILanguage
	baseURL    *url.URL
	version    int
	common     service // Reuse this for all services.

	Establishments *EstablishmentsInstance
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
	client.Establishments = (*EstablishmentsInstance)(&client.common)

	return client, nil
}

// SetLanguage sets the response language.
func (c *Client) SetLanguage(l APILanguage) error {
	languages := []APILanguage{English, Cymraeg}
	for _, lang := range languages {
		if l == lang {
			c.language = l
			return nil
		}
	}

	return errors.New("Language not supported")
}
