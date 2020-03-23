package fhrs

import (
	"fmt"
	"net/url"
	"strconv"
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

// SearchParams are the parameters available for searching for establishments.
type SearchParams struct {
	Name              string
	Address           string
	Longitude         *float64
	Latitude          *float64
	MaxDistanceLimit  *int
	BusinessTypeID    string
	SchemeTypeKey     string
	RatingKey         string
	RatingOperatorKey string
	LocalAuthorityID  string
	CountryID         string
	SortOptionKey     string
	PageNumber        *int
	PageSize          *int
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

// Search returns establishments matching the given set of parameters.
//
// https://api.ratings.food.gov.uk/Help/Api/GET-Establishments_name_address_longitude_latitude_maxDistanceLimit
// _businessTypeId_schemeTypeKey_ratingKey_ratingOperatorKey_localAuthorityId_countryId_sortOptionKey_pageNumber_pageSize
func (s *EstablishmentsService) Search(params *SearchParams) (*Establishments, error) {
	var establishments *Establishments
	u := url.URL{Path: "Establishments"}
	q := u.Query()

	// We could have used struct tags and reflection here but as they're all
	// different types here this is boring but more clear.
	if params != nil {
		if params.Name != "" {
			q.Set("name", params.Name)
		}
		if params.Address != "" {
			q.Set("address", params.Address)
		}
		if params.Longitude != nil {
			q.Set("longitude", strconv.FormatFloat(*params.Longitude, 'f', -1, 64))
		}
		if params.Latitude != nil {
			q.Set("latitude", strconv.FormatFloat(*params.Latitude, 'f', -1, 64))
		}
		if params.MaxDistanceLimit != nil {
			q.Set("maxDistanceLimit", strconv.Itoa(*params.MaxDistanceLimit))
		}
		if params.BusinessTypeID != "" {
			q.Set("businessTypeId", params.BusinessTypeID)
		}
		if params.SchemeTypeKey != "" {
			q.Set("schemeTypeKey", params.SchemeTypeKey)
		}
		if params.RatingKey != "" {
			q.Set("ratingKey", params.RatingKey)
		}
		if params.RatingOperatorKey != "" {
			q.Set("ratingOperatorKey", params.RatingOperatorKey)
		}
		if params.LocalAuthorityID != "" {
			q.Set("localAuthorityId", params.LocalAuthorityID)
		}
		if params.CountryID != "" {
			q.Set("countryId", params.CountryID)
		}
		if params.SortOptionKey != "" {
			q.Set("sortOptionKey", params.SortOptionKey)
		}
		if params.PageNumber != nil {
			q.Set("pageNumber", strconv.Itoa(*params.PageNumber))
		}
		if params.PageSize != nil {
			q.Set("pageSize", strconv.Itoa(*params.PageSize))
		}
	}

	u.RawQuery = q.Encode()
	if err := s.client.get(u.String(), &establishments); err != nil {
		return nil, err
	}

	return establishments, nil
}
