package fhrs

import (
	"encoding/json"
	"fmt"
	"strings"
	"net/http"
	"strconv"
	"time"
)

type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
    if s == "null" || s == "" || s == "undefined" {
		*t = Timestamp{time.Time{}}
		return nil
    }

	parsed, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err	
	}

	*t = Timestamp{parsed}
	return nil
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(t.Time.Format("2006-01-02T15:04:05")), nil
}

type EstablishmentsService interface {
	GetByID(id string) (*Establishment, error)
}

type EstablishmentsInstance service

type ErrorResponse struct {
	Message string `json:"Message"`
}

type Establishments struct {
	Establishments []Establishment `json:"establishments"`
	Meta           Meta            `json:"meta"`
	Links          []Links         `json:"links"`
}

type Scores struct {
	Hygiene                int `json:"Hygiene"`
	Structural             int `json:"Structural"`
	ConfidenceInManagement int `json:"ConfidenceInManagement"`
}

type Geocode struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

type Meta struct {
	DataSource  string    `json:"dataSource"`
	ExtractDate Timestamp `json:"extractDate"`
	ItemCount   int       `json:"itemCount"`
	Returncode  string    `json:"returncode"`
	TotalCount  int       `json:"totalCount"`
	TotalPages  int       `json:"totalPages"`
	PageSize    int       `json:"pageSize"`
	PageNumber  int       `json:"pageNumber"`
}

type Links struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type Establishment struct {
	FHRSID                     int       `json:"FHRSID"`
	LocalAuthorityBusinessID   string    `json:"LocalAuthorityBusinessID"`
	BusinessName               string    `json:"BusinessName"`
	BusinessType               string    `json:"BusinessType"`
	BusinessTypeID             int       `json:"BusinessTypeID"`
	AddressLine1               string    `json:"AddressLine1"`
	AddressLine2               string    `json:"AddressLine2"`
	AddressLine3               string    `json:"AddressLine3"`
	AddressLine4               string    `json:"AddressLine4"`
	PostCode                   string    `json:"PostCode"`
	Phone                      string    `json:"Phone"`
	RatingValue                string    `json:"RatingValue"`
	RatingKey                  string    `json:"RatingKey"`
	RatingDate                 Timestamp `json:"RatingDate"`
	LocalAuthorityCode         string    `json:"LocalAuthorityCode"`
	LocalAuthorityName         string    `json:"LocalAuthorityName"`
	LocalAuthorityWebSite      string    `json:"LocalAuthorityWebSite"`
	LocalAuthorityEmailAddress string    `json:"LocalAuthorityEmailAddress"`
	Scores                     Scores    `json:"scores"`
	SchemeType                 string    `json:"SchemeType"`
	Geocode                    Geocode   `json:"geocode"`
	RightToReply               string    `json:"RightToReply"`
	Distance                   float64   `json:"Distance"`
	NewRatingPending           bool      `json:"NewRatingPending"`
	Meta                       Meta      `json:"meta"`
	Links                      []Links   `json:"links"`
}

// GetByID returns an establishment with the given FHRSID.
//
// https://api.ratings.food.gov.uk/Help/Api/GET-Establishments-id
func (s *EstablishmentsInstance) GetByID(id string) (*Establishment, error) {
	url, err := s.client.baseURL.Parse(fmt.Sprintf("Establishments/%s", id))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-version", strconv.Itoa(s.client.version))
	req.Header.Set("Accept-Language", s.client.language.String())

	res, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var establishment Establishment
	if err := json.NewDecoder(res.Body).Decode(&establishment); err != nil {
		return nil, err
	}

	return &establishment, nil
}
