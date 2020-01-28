package fhrs

import (
	"fmt"
)

// EstablishmentsService encapsulates the Establishments methods of the API.
//
// https://api.ratings.food.gov.uk/help#Establishments
type EstablishmentsService service

type Establishments struct {
	Establishments []Establishment `json:"establishments"`
	Meta           Meta            `json:"meta"`
	Links          []Link          `json:"links"`
}

type Scores struct {
	Hygiene                *int `json:"Hygiene"`
	Structural             *int `json:"Structural"`
	ConfidenceInManagement *int `json:"ConfidenceInManagement"`
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

type Link struct {
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
	Distance                   *float64  `json:"Distance"`
	NewRatingPending           bool      `json:"NewRatingPending"`
	Meta                       Meta      `json:"meta"`
	Links                      []Link    `json:"links"`
}

// GetByID returns an establishment with the given FHRSID.
//
// https://api.ratings.food.gov.uk/Help/Api/GET-Establishments-id
func (s *EstablishmentsService) GetByID(id string) (*Establishment, error) {
	var establishment *Establishment
	if err := s.client.get(fmt.Sprintf("Establishments/%s", id), &establishment); err != nil {
		return nil, err
	}

	return establishment, nil
}
