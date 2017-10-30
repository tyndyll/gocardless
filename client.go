package gocardless

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	APIVersion     = `2015-07-06`
	BaseLiveURL    = `https://api.gocardless.com`
	BaseSandboxURL = `https://api-sandbox.gocardless.com`

	jsonMimeType = `application/json`
)

// API defines the interface to interact with the GoCardless API. An instance of the Client is returned by calling
// NewClient. A mock type called MockClient is also provided
type API interface {
	CreateCustomer(*Customer) error
	GetCustomer(string) (*Customer, error)
	ListCustomer() ([]*Customer, error)
	UpdateCustomer(*Customer) error
}

// Client is an implementation of the GoCardless API interface.
type Client struct {
	// AccessToken is the bearer token used to authenticate requests to the GoCardless API
	AccessToken string
	// RemoteURL is the address of the GoCardless API
	RemoteURL string
}

func (c *Client) do(req *http.Request) (*Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		return nil, errors.New("StatusTooManyRequests")
	}

	return &Response{resp}, nil
}

func (c *Client) newRequest(path, method string, data []byte) (*http.Request, error) {
	if strings.ToUpper(method) == http.MethodPatch {
		return nil, errors.New(InvalidMethodError)
	}

	endpoint := fmt.Sprintf("%s%s", c.RemoteURL, path)

	body := ioutil.NopCloser(bytes.NewBuffer(data))
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add(`Authorization`, fmt.Sprintf(`Bearer %s`, c.AccessToken))
	req.Header.Add(`GoCardless-Version`, APIVersion)
	req.Header.Add(`Accept`, jsonMimeType)
	req.Header.Add(`Content-Type`, jsonMimeType)

	return req, nil
}

func (c *Client) decodeError(error []byte) *Error {
	newErr := &errorContainer{}
	if err := json.Unmarshal(error, newErr); err != nil {
		return nil
	}
	return newErr.Error
}

// NewClient returns a populated API Client which has been configured for the supplied environment, using the Access
// Token for authenticating all requests. An error will be returned in cases where the environment is not recognised
func NewClient(accessToken string, environment Environment) (API, error) {
	c := &Client{
		AccessToken: accessToken,
	}

	switch environment {
	case SandboxEnvironment:
		c.RemoteURL = BaseSandboxURL
	case LiveEnvironment:
		c.RemoteURL = BaseLiveURL
	default:
		return nil, errors.New(fmt.Sprintf("%s is not a valid environment", environment))
	}

	return c, nil
}
