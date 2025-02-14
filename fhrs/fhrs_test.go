package fhrs

import (
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http/httptest"
	"net/url"
	"testing"
)

func getTestEnv() (*Client, *httptest.Server, *httprouter.Router, error) {
	router := httprouter.New()
	server := httptest.NewUnstartedServer(router)

	// We need to assign our own listener here as server.URL
	// returns nil until it's started. We assign
	// a blank port so that it is randomised, as per
	// https://golang.org/pkg/net/#Listen
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		return nil, nil, nil, err
	}

	server.Listener = listener

	client, err := NewClient()
	if err != nil {
		return nil, nil, nil, err
	}

	u, err := url.Parse("http://" + listener.Addr().String() + "/")
	if err != nil {
		return nil, nil, nil, err
	}

	client.baseURL = u

	return client, server, router, nil
}

func TestNewClient(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Error(err)
	}

	if c.language != LanguageEnglish {
		t.Error("Expected default language to be English")
	}
}

func TestSetLanguage(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Error(err)
	}

	err = c.SetLanguage(LanguageEnglish)
	if err != nil {
		t.Error(err)
	}

	err = c.SetLanguage(LanguageCymraeg)
	if err != nil {
		t.Error(err)
	}

	err = c.SetLanguage(9)
	if err == nil {
		t.Error("Should not be able to set unsupported language")
	}
}

func TestAPILanguageString(t *testing.T) {
	cases := []struct {
		want string
		have string
	}{
		{want: "en-GB", have: LanguageEnglish.String()},
		{want: "cy-GB", have: LanguageCymraeg.String()},
	}

	for _, c := range cases {
		if c.have != c.want {
			t.Errorf("Expected %s but got %s", c.want, c.have)
		}
	}
}
