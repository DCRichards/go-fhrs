package fhrs

import (
	"io"
	"net/http"
	"testing"
)

func TestGetByID(t *testing.T) {
	client, server, router, err := setup()
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

	router.HandleFunc("/Establishments/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, body)
	})

	server.Start()

	est, err := client.Establishments.GetByID("1")
	if err != nil {
		t.Error(err)
	}
}
