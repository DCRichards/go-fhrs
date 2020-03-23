package fhrs

import (
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	client, server, router, err := getTestEnv()
	if err != nil {
		t.Error(err)
	}

	server.Start()
	defer server.Close()

	body := `{
	  "ratings": [
		{
		  "ratingId": 12,
		  "ratingName": "5",
		  "ratingKey": "fhrs_5_en-gb",
		  "ratingKeyName": "5",
		  "schemeTypeId": 1,
		  "links": [
			{
			  "rel": "self",
			  "href": "http://api.ratings.food.gov.uk/ratings/12"
			}
		  ]
		},
		{
		  "ratingId": 11,
		  "ratingName": "4",
		  "ratingKey": "fhrs_4_en-gb",
		  "ratingKeyName": "4",
		  "schemeTypeId": 1,
		  "links": [
			{
			  "rel": "self",
			  "href": "http://api.ratings.food.gov.uk/ratings/11"
			}
		  ]
		}
	  ],
	  "meta": {
		"dataSource": "API",
		"extractDate": "2020-02-03T22:32:34.2688747+00:00",
		"itemCount": 11,
		"returncode": "OK",
		"totalCount": 11,
		"totalPages": 1,
		"pageSize": 11,
		"pageNumber": 1
	  },
	  "links": [
		{
		  "rel": "self",
		  "href": "http://api.ratings.food.gov.uk/ratings"
		}
	  ]
	}`

	ed, _ := time.Parse(time.RFC3339Nano, "2020-02-03T22:32:34.2688747+00:00")

	expected := &Ratings{
		Ratings: []Rating{
			{
				RatingID:      12,
				RatingName:    "5",
				RatingKey:     "fhrs_5_en-gb",
				RatingKeyName: "5",
				SchemeTypeID:  1,
				Links: []Link{
					{
						Rel:  "self",
						Href: "http://api.ratings.food.gov.uk/ratings/12",
					},
				},
			},
			{
				RatingID:      11,
				RatingName:    "4",
				RatingKey:     "fhrs_4_en-gb",
				RatingKeyName: "4",
				SchemeTypeID:  1,
				Links: []Link{
					{
						Rel:  "self",
						Href: "http://api.ratings.food.gov.uk/ratings/11",
					},
				},
			},
		},
		Meta: Meta{
			DataSource:  "API",
			ExtractDate: Timestamp(ed),
			ItemCount:   11,
			Returncode:  "OK",
			TotalCount:  11,
			TotalPages:  1,
			PageSize:    11,
			PageNumber:  1,
		},
		Links: []Link{
			{
				Rel:  "self",
				Href: "http://api.ratings.food.gov.uk/ratings",
			},
		},
	}

	router.GET("/Ratings", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, body)
	})

	actual, err := client.Ratings.Get()
	if err != nil {
		t.Error(err)
	}

	if actual == nil {
		t.Error("Expected response but got nil")
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected:\n%+v\nBut got:\n%+v\n", expected, actual)
	}
}

func TestGet_Error(t *testing.T) {
	client, server, router, err := getTestEnv()
	if err != nil {
		t.Error(err)
	}

	server.Start()
	defer server.Close()

	errorMsg := "The service is unavailable."

	router.GET("/Ratings", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusServiceUnavailable)
		io.WriteString(w, errorMsg)
	})

	_, err = client.Ratings.Get()
	if err == nil {
		t.Error("Expected an error to be returned but got nil")
	}

	apiErr, ok := err.(APIError)
	if !ok {
		t.Error("Expected error to be an APIError")
	}

	if apiErr.Message != errorMsg {
		t.Errorf("Expected message to be %s but got %s", errorMsg, apiErr.Message)
	}

	if apiErr.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected status code to be %d but got %d", http.StatusServiceUnavailable, apiErr.StatusCode)
	}
}
