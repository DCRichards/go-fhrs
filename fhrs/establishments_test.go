package fhrs

import (
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestGetByID(t *testing.T) {
	client, server, router, err := getTestEnv()
	if err != nil {
		t.Error(err)
	}

	defer server.Close()

	body := `{
	  "FHRSID": 82940,
	  "LocalAuthorityBusinessID": "2019",
	  "BusinessName": "Ali's",
	  "BusinessType": "Restaurant/Cafe/Canteen",
	  "BusinessTypeID": 1,
	  "AddressLine1": "89 Commercial Road",
	  "AddressLine2": "Portsmouth",
	  "AddressLine3": "",
	  "AddressLine4": "",
	  "PostCode": "PO1 1BA",
	  "Phone": "",
	  "RatingValue": "3",
	  "RatingKey": "fhrs_3_en-gb",
	  "RatingDate": "2019-08-06T00:00:00",
	  "LocalAuthorityCode": "876",
	  "LocalAuthorityName": "Portsmouth",
	  "LocalAuthorityWebSite": "http://www.portsmouth.gov.uk",
	  "LocalAuthorityEmailAddress": "public.protection@portsmouthcc.gov.uk",
	  "scores": {
		"Hygiene": null,
		"Structural": null,
		"ConfidenceInManagement": null
	  },
	  "SchemeType": "FHRS",
	  "geocode": {
		"longitude": "-1.09159100055695",
		"latitude": "50.7984199523926"
	  },
	  "RightToReply": "",
	  "Distance": null,
	  "NewRatingPending": false,
	  "meta": {
		"dataSource": "Lucene",
		"extractDate": "0001-01-01T00:00:00",
		"itemCount": 0,
		"returncode": "OK",
		"totalCount": 1,
		"totalPages": 1,
		"pageSize": 1,
		"pageNumber": 1
	  },
	  "links": [
		{
		  "rel": "self",
		  "href": "http://api.ratings.food.gov.uk/establishments/82940"
		}
	  ]
	}`

	expected := &Establishment{
		FHRSID:                     82940,
		LocalAuthorityBusinessID:   "2019",
		BusinessName:               "Ali's",
		BusinessType:               "Restaurant/Cafe/Canteen",
		BusinessTypeID:             1,
		AddressLine1:               "89 Commercial Road",
		AddressLine2:               "Portsmouth",
		PostCode:                   "PO1 1BA",
		RatingValue:                "3",
		RatingKey:                  "fhrs_3_en-gb",
		LocalAuthorityCode:         "876",
		LocalAuthorityName:         "Portsmouth",
		LocalAuthorityWebSite:      "http://www.portsmouth.gov.uk",
		LocalAuthorityEmailAddress: "public.protection@portsmouthcc.gov.uk",
		Scores:                     Scores{},
		SchemeType:                 "FHRS",
		Geocode: Geocode{
			Longitude: "-1.09159100055695",
			Latitude:  "50.7984199523926",
		},
		NewRatingPending: false,
		Meta: Meta{
			DataSource: "Lucene",
			ItemCount:  0,
			Returncode: "OK",
			TotalCount: 1,
			TotalPages: 1,
			PageSize:   1,
			PageNumber: 1,
		},
		Links: []Link{
			{
				Rel:  "self",
				Href: "http://api.ratings.food.gov.uk/establishments/82940",
			},
		},
	}
	expected.RatingDate.Time, _ = time.Parse("2006-01-02T15:04:05", "2019-08-06T00:00:00")
	expected.Meta.ExtractDate.Time, _ = time.Parse("2006-01-02T15:04:05", "0001-01-01T00:00:00")

	router.HandleFunc("/Establishments/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, body)
	})

	server.Start()

	actual, err := client.Establishments.GetByID("1")
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

func TestGetByID_NotFound(t *testing.T) {
	client, server, router, err := getTestEnv()
	if err != nil {
		t.Error(err)
	}

	defer server.Close()

	body := `{
		"Message": "No establishment found with EstablishmentId: 0"
	}`

	router.HandleFunc("/Establishments/0", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, body)
	})

	server.Start()

	est, err := client.Establishments.GetByID("0")
	if err != nil {
		t.Error(err)
	}

	if est != nil {
		t.Errorf("Expected response to be nil, but got %v", est)
	}
}
