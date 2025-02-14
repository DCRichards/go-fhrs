package fhrs

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
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

	server.Start()
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

	rd, _ := time.Parse("2006-01-02T15:04:05", "2019-08-06T00:00:00")
	ed, _ := time.Parse("2006-01-02T15:04:05", "0001-01-01T00:00:00")

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
		RatingDate:                 Timestamp(rd),
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
			DataSource:  "Lucene",
			ExtractDate: Timestamp(ed),
			ItemCount:   0,
			Returncode:  "OK",
			TotalCount:  1,
			TotalPages:  1,
			PageSize:    1,
			PageNumber:  1,
		},
		Links: []Link{
			{
				Rel:  "self",
				Href: "http://api.ratings.food.gov.uk/establishments/82940",
			},
		},
	}

	idQuery := "1"

	router.GET("/Establishments/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if q := p.ByName("id"); q != idQuery {
			t.Errorf("Expected ID to be %s but got %s", idQuery, q)
		}

		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, body)
	})

	actual, err := client.Establishments.GetByID(idQuery)
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

func TestGetByID_BadRequest(t *testing.T) {
	client, server, router, err := getTestEnv()
	if err != nil {
		t.Error(err)
	}

	server.Start()
	defer server.Close()

	idQuery := "AAAA"
	errorMessage := "The request is invalid"
	body := fmt.Sprintf(`{ "Message": "%s" }`, errorMessage)

	router.GET("/Establishments/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if q := p.ByName("id"); q != idQuery {
			t.Errorf("Expected ID to be %s but got %s", idQuery, q)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, body)
	})

	_, err = client.Establishments.GetByID(idQuery)
	var apiError APIError
	if errors.As(err, &apiError) {
		if apiError.Message != errorMessage {
			t.Errorf("Expected err.Message to be %s but got %s", errorMessage, apiError.Message)
		}
	} else {
		t.Errorf("Expected err to be APIError but type is %T", err)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	client, server, router, err := getTestEnv()
	if err != nil {
		t.Error(err)
	}

	server.Start()
	defer server.Close()

	idQuery := "0"
	body := `{ "Message": "No establishment found with EstablishmentId: 0" }`

	router.GET("/Establishments/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if q := p.ByName("id"); q != idQuery {
			t.Errorf("Expected ID to be %s but got %s", idQuery, q)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, body)
	})

	est, err := client.Establishments.GetByID(idQuery)
	if err != nil {
		t.Error(err)
	}

	if est != nil {
		t.Errorf("Expected response to be nil, but got %v", est)
	}
}

func TestGetByID_Headers(t *testing.T) {
	client, server, router, err := getTestEnv()
	if err != nil {
		t.Error(err)
	}

	server.Start()
	defer server.Close()

	idQuery := "21188"
	// 󠁧󠁢󠁷󠁬󠁳󠁿Sorry Welsh API developers, the message is always in English. I checked.
	body := `{ "Message": "No establishment found with EstablishmentId: 21188" }`

	router.GET("/Establishments/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if ah := r.Header.Get("x-api-version"); ah != "2" {
			t.Errorf("Expected x-api-version to be 2 but got %s", ah)
		}

		if lh := r.Header.Get("Accept-Language"); lh != "cy-GB" {
			t.Errorf("Expected Accept-Language to be cy-GB but got %s", lh)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, body)
	})

	if err := client.SetLanguage(LanguageCymraeg); err != nil {
		t.Error(err)
	}

	_, err = client.Establishments.GetByID(idQuery)
	if err != nil {
		t.Error(err)
	}
}

func TestSearch(t *testing.T) {
	client, server, router, err := getTestEnv()
	if err != nil {
		t.Error(err)
	}

	server.Start()
	defer server.Close()

	body := `{
	  "establishments":[
		{
		  "FHRSID":82940,
		  "LocalAuthorityBusinessID":"2019",
		  "BusinessName":"Ali's",
		  "BusinessType":"Restaurant/Cafe/Canteen",
		  "BusinessTypeID":1,
		  "AddressLine1":"89 Commercial Road",
		  "AddressLine2":"Portsmouth",
		  "AddressLine3":"",
		  "AddressLine4":"",
		  "PostCode":"PO1 1BA",
		  "Phone":"",
		  "RatingValue":"3",
		  "RatingKey":"fhrs_3_en-gb",
		  "RatingDate":"2019-08-06T00:00:00",
		  "LocalAuthorityCode":"876",
		  "LocalAuthorityName":"Portsmouth",
		  "LocalAuthorityWebSite":"http://www.portsmouth.gov.uk",
		  "LocalAuthorityEmailAddress":"public.protection@portsmouthcc.gov.uk",
		  "scores":{
			"Hygiene":null,
			"Structural":null,
			"ConfidenceInManagement":null
		  },
		  "SchemeType":"FHRS",
		  "geocode":{
			"longitude":"-1.09159100055695",
			"latitude":"50.7984199523926"
		  },
		  "RightToReply":"",
		  "Distance":null,
		  "NewRatingPending":false,
		  "meta":{
			"dataSource":"Lucene",
			"extractDate":"0001-01-01T00:00:00",
			"itemCount":0,
			"returncode":"OK",
			"totalCount":1,
			"totalPages":1,
			"pageSize":1,
			"pageNumber":1
		  },
		  "links":[
			{
			  "rel":"self",
			  "href":"http://api.ratings.food.gov.uk/establishments/82940"
			}
		  ]
		}
	  ],
	  "meta":{
		"dataSource":"Lucene",
		"extractDate":"0001-01-01T00:00:00",
		"itemCount":0,
		"returncode":"OK",
		"totalCount":1,
		"totalPages":1,
		"pageSize":1,
		"pageNumber":1
	  },
	  "links":[
		{
		  "rel":"self",
		  "href":"http://api.ratings.food.gov.uk/establishments"
		}
	  ]
	}`

	rd, _ := time.Parse("2006-01-02T15:04:05", "2019-08-06T00:00:00")
	ed, _ := time.Parse("2006-01-02T15:04:05", "0001-01-01T00:00:00")

	expected := &Establishments{
		Establishments: []Establishment{
			{
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
				RatingDate:                 Timestamp(rd),
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
					DataSource:  "Lucene",
					ExtractDate: Timestamp(ed),
					ItemCount:   0,
					Returncode:  "OK",
					TotalCount:  1,
					TotalPages:  1,
					PageSize:    1,
					PageNumber:  1,
				},
				Links: []Link{
					{
						Rel:  "self",
						Href: "http://api.ratings.food.gov.uk/establishments/82940",
					},
				},
			},
		},
		Meta: Meta{
			DataSource:  "Lucene",
			ExtractDate: Timestamp(ed),
			ItemCount:   0,
			Returncode:  "OK",
			TotalCount:  1,
			TotalPages:  1,
			PageSize:    1,
			PageNumber:  1,
		},
		Links: []Link{
			{
				Rel:  "self",
				Href: "http://api.ratings.food.gov.uk/establishments",
			},
		},
	}

	params := &SearchParams{
		Name: "Ali's",
	}

	router.GET("/Establishments", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		q := r.URL.Query()
		if q["name"][0] != params.Name {
			t.Errorf("Expected param 'name' to equal %s but got %s", params.Name, q["name"][0])
		}

		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, body)
	})

	actual, err := client.Establishments.Search(params)
	if err != nil {
		t.Error(err)
	}

	if len(actual.Establishments) == 0 {
		t.Error("Expected response but got nil")
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected:\n%+v\nBut got:\n%+v\n", expected, actual)
	}
}

func TestSearch_NoResults(t *testing.T) {
	client, server, router, err := getTestEnv()
	if err != nil {
		t.Error(err)
	}

	server.Start()
	defer server.Close()

	body := `{
	  "establishments": [],
	  "meta":{
		"dataSource":"Lucene",
		"extractDate":"0001-01-01T00:00:00",
		"itemCount":0,
		"returncode":"OK",
		"totalCount":1,
		"totalPages":1,
		"pageSize":1,
		"pageNumber":1
	  },
	  "links":[
		{
		  "rel":"self",
		  "href":"http://api.ratings.food.gov.uk/establishments"
		}
	  ]
	}`

	ed, _ := time.Parse("2006-01-02T15:04:05", "0001-01-01T00:00:00")

	expected := &Establishments{
		Establishments: []Establishment{},
		Meta: Meta{
			DataSource:  "Lucene",
			ExtractDate: Timestamp(ed),
			ItemCount:   0,
			Returncode:  "OK",
			TotalCount:  1,
			TotalPages:  1,
			PageSize:    1,
			PageNumber:  1,
		},
		Links: []Link{
			{
				Rel:  "self",
				Href: "http://api.ratings.food.gov.uk/establishments",
			},
		},
	}

	params := &SearchParams{
		Address: "Obscurenameplace",
	}

	router.GET("/Establishments", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		q := r.URL.Query()
		if q["address"][0] != params.Address {
			t.Errorf("Expected param 'address' to equal %s but got %s", params.Address, q["address"][0])
		}

		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, body)
	})

	actual, err := client.Establishments.Search(params)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected:\n%+v\nBut got:\n%+v\n", expected, actual)
	}
}

func TestSearch_Params(t *testing.T) {
	client, server, router, err := getTestEnv()
	if err != nil {
		t.Error(err)
	}

	server.Start()
	defer server.Close()

	lat := 50.8261
	lon := -0.8231
	maxDistanceLimit := 10
	pageNumber := 1
	pageSize := 20

	params := &SearchParams{
		Name:              "Bob's Burgers",
		Address:           "'Merica",
		Longitude:         &lat,
		Latitude:          &lon,
		MaxDistanceLimit:  &maxDistanceLimit,
		BusinessTypeID:    "10",
		SchemeTypeKey:     "1",
		RatingKey:         "AwatingInspection",
		RatingOperatorKey: "LessThanOrEqual",
		LocalAuthorityID:  "HCC",
		CountryID:         "GB",
		SortOptionKey:     "rating",
		PageNumber:        &pageNumber,
		PageSize:          &pageSize,
	}

	paramError := func(name string, expected interface{}, actual string) string {
		return fmt.Sprintf("Expected param '%s' to be %v but got %s", name, expected, actual)
	}

	testValues := map[string]interface{}{
		"name":              params.Name,
		"address":           params.Address,
		"longitude":         "50.8261",
		"latitude":          "-0.8231",
		"maxDistanceLimit":  "10",
		"businessTypeId":    params.BusinessTypeID,
		"schemeTypeKey":     params.SchemeTypeKey,
		"ratingKey":         params.RatingKey,
		"ratingOperatorKey": params.RatingOperatorKey,
		"localAuthorityId":  params.LocalAuthorityID,
		"countryId":         params.CountryID,
		"sortOptionKey":     params.SortOptionKey,
		"pageNumber":        "1",
		"pageSize":          "20",
	}

	router.GET("/Establissments", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		q := r.URL.Query()

		for p, expected := range testValues {
			actual := q[p][0]
			if actual != expected {
				t.Errorf(paramError(p, expected, actual))
			}
		}
	})

	_, err = client.Establishments.Search(params)
	if err != nil {
		t.Error(err)
	}
}
