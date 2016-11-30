// Some of the code in this file is derived from the examples written by Google
// that are provided under a BSD-style license.
//
// https://github.com/google/google-api-go-client/blob/master/examples/main.go
// https://github.com/google/google-api-go-client/blob/master/LICENSE

package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

// ClientSource is a source of pre-authenticated HTTP clients
type ClientSource interface {
	Client(ctx context.Context, subject string) (*http.Client, error)
}

type serviceSource struct {
	config jwt.Config
}

func (s *serviceSource) Client(ctx context.Context, subject string) (*http.Client, error) {
	cfg := s.config // Copy the config so that we don't modify it
	if subject != "" {
		cfg.Subject = subject
	}

	return cfg.Client(ctx), nil
}

// NewSourceFromKeyfile creates a service-based client source using the
// configuration provided in the specified key file, which must be in JSON
// format.
func NewSourceFromKeyfile(keyfile string, scope ...string) (ClientSource, error) {
	if keyfile == "" {
		return nil, errors.New("No keyfile specified")
	}

	data, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	config, err := google.JWTConfigFromJSON(data, scope...)
	if err != nil {
		return nil, err
	}

	return &serviceSource{
		config: *config,
	}, nil
}

type tokenSource struct {
	config oauth2.Config
	token  *oauth2.Token
}

func (s *tokenSource) Client(ctx context.Context, subject string) (*http.Client, error) {
	return s.config.Client(ctx, s.token), nil
}

// NewSourceFromToken creates a token-based client source using the
// provided client ID and secret.
func NewSourceFromToken(ctx context.Context, clientID, secret string, scope ...string) (ClientSource, error) {
	if clientID == "" {
		return nil, errors.New("No client ID provided")
	}
	if secret == "" {
		return nil, errors.New("No secret provided")
	}

	source := &tokenSource{
		config: oauth2.Config{
			ClientID:     clientID,
			ClientSecret: secret,
			Endpoint:     google.Endpoint,
			Scopes:       scope,
		},
	}

	var err error
	source.token, err = tokenFromWeb(ctx, &source.config)
	if err != nil {
		return nil, err
	}

	return source, nil
}

func tokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	ch := make(chan string)
	randState := fmt.Sprintf("st%d", time.Now().UnixNano())

	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		return nil, fmt.Errorf("Network listener error: %v", err)
	}

	handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/favicon.ico" {
			http.Error(rw, "", 404)
			return
		}
		if req.FormValue("state") != randState {
			log.Printf("State doesn't match: req = %#v", req)
			http.Error(rw, "", 500)
			return
		}
		if code := req.FormValue("code"); code != "" {
			fmt.Fprintf(rw, "<h1>Success</h1>Authorized.")
			rw.(http.Flusher).Flush()
			ch <- code
			return
		}
		log.Printf("no code")
		http.Error(rw, "", 500)
	})

	ts := &httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: handler},
	}

	defer ts.Close()
	ts.Start()

	config.RedirectURL = ts.URL
	authURL := config.AuthCodeURL(randState)
	open.Start(authURL)
	fmt.Printf("Authorize this app at: %s\n", authURL)
	code := <-ch

	//log.Printf("Got code: %s", code)

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("Token exchange error: %v", err)
	}

	fmt.Printf("Authorization successful.\n")

	return token, nil
}
