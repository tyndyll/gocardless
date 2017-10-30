package gocardless

import (
	"testing"

	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
)

func TestNewClient(t *testing.T) {
	Convey(`Given I have an access token`, t, func() {
		accessToken := `abcdef`

		Convey(`And I have a sandbox environment`, func() {
			environment := SandboxEnvironment

			Convey(`When I call NewClient`, func() {
				client, err := NewClient(accessToken, environment)

				Convey(`Then the AccessToken field in the underlying type will be set`, func() {
					So(client.(*Client).AccessToken, ShouldEqual, accessToken)
				})

				Convey(`Then the BaseURL field in the underlying type will be set`, func() {
					So(client.(*Client).RemoteURL, ShouldEqual, BaseSandboxURL)
				})

				Convey(`Then the returned error will be nil`, func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey(`And I have a live environment`, func() {
			environment := LiveEnvironment

			Convey(`When I call NewClient`, func() {
				client, _ := NewClient(accessToken, environment)

				Convey(`Then the BaseURL field in the underlying type will be set`, func() {
					So(client.(*Client).RemoteURL, ShouldEqual, BaseLiveURL)
				})
			})
		})

		Convey(`And I have an invalid environment`, func() {
			environment := Environment(`Undead. Not live. It's funny'`)

			Convey(`When I call NewClient`, func() {
				client, err := NewClient(accessToken, environment)

				Convey(`Then the client will be nil`, func() {
					So(client, ShouldBeNil)
				})

				Convey(`Then the error will be not nil`, func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}

func TestClientDecodeError(t *testing.T) {
	Convey(`Given I have a client`, t, func() {
		client := &Client{}

		Convey(`And I have a valid GoCardless error message`, func() {
			msg := []byte(`{
				"error": {
					"documentation_url": "https://developer.gocardless.com/#validation_failed",
					"message": "Validation failed",
					"type": "validation_failed",
					"code": 422,
					"request_id": "dd50eaaf-8213-48fe-90d6-5466872efbc4",
					"errors": [
						{
							"message": "must be a number",
							"field": "branch_code",
							"request_pointer": "/customer_bank_accounts/branch_code"
						}, {
							"message": "is the wrong length (should be 8 characters)",
							"field": "branch_code",
							"request_pointer": "/customer_bank_accounts/branch_code"
						}
					]
				}
			}`)

			Convey(`When I call decodeError`, func() {
				err := client.decodeError(msg)

				Convey(`Then the Documentation URL will be populated`, func() {
					So(err.DocumentationURL, ShouldEqual, `https://developer.gocardless.com/#validation_failed`)
				})

				Convey(`Then the Message will be populated`, func() {
					So(err.Message, ShouldEqual, `Validation failed`)
				})

				Convey(`Then the Type will be populated`, func() {
					So(err.Type, ShouldEqual, `validation_failed`)
				})

				Convey(`Then the Code will be populated`, func() {
					So(err.Code, ShouldEqual, 422)
				})

				Convey(`Then the RequestID will be populated`, func() {
					So(err.RequestID, ShouldEqual, `dd50eaaf-8213-48fe-90d6-5466872efbc4`)
				})

				Convey(`Then the Details will contain 2 items`, func() {
					So(len(err.Details), ShouldEqual, 2)
				})
			})
		})

		Convey(`And I have a invalid GoCardless error message`, func() {
			msg := []byte(`totally invalid JSON`)

			Convey(`When I call decodeError`, func() {
				err := client.decodeError(msg)

				Convey(`Then the err will be nil`, func() {
					// TODO: this doesn't feel right
					So(err, ShouldBeNil)
				})
			})
		})
	})
}

func TestNewRequest(t *testing.T) {
	Convey(`Given I have a client`, t, func() {
		client := &Client{
			AccessToken: `test-demo-access-token`,
			RemoteURL:   BaseSandboxURL,
		}

		Convey(`And I have a path`, func() {
			path := `/demo-path`

			Convey(`And I am making a POST request`, func() {
				method := http.MethodPost

				Convey(`When I call newRequest`, func() {
					req, err := client.newRequest(path, method, nil)
					if err != nil {
						panic(err)
					}

					Convey(`Then the Authorization header will be set`, func() {
						So(req.Header.Get(`Authorization`), ShouldEqual, fmt.Sprintf(`Bearer %s`, client.AccessToken))
					})

					Convey(`Then the GoCardless-Version header will be set`, func() {
						So(req.Header.Get(`GoCardless-Version`), ShouldEqual, APIVersion)
					})

					Convey(`Then the Accept header will be set`, func() {
						So(req.Header.Get(`Accept`), ShouldEqual, jsonMimeType)
					})

					Convey(`Then the Content-Type header will be set`, func() {
						So(req.Header.Get(`Content-Type`), ShouldEqual, jsonMimeType)
					})

					// TODO: Check body
				})
			})

			Convey(`And I am making a PATCH request`, func() {
				method := http.MethodPatch

				Convey(`When I call newRequest`, func() {
					req, err := client.newRequest(path, method, nil)

					Convey(`Then the request will be nil`, func() {
						So(req, ShouldBeNil)
					})

					Convey(`Then the error will be the InvalidMethodError`, func() {
						So(err.Error(), ShouldEqual, InvalidMethodError)
					})
				})
			})
		})

	})
}

func TestClientDo(t *testing.T) {
	Convey(`Given I have a client`, t, func() {
		client := &Client{}

		Convey(`And a server that returns a StatusOK response`, func() {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				return
			}))
			client.RemoteURL = srv.URL

			Convey(`And I have a request`, func() {
				req, err := client.newRequest(`/`, http.MethodGet, nil)
				if err != nil {
					panic(err)
				}

				Convey(`When I call the do method`, func() {
					resp, err := client.do(req)

					Convey(`Then the response will not be nil`, func() {
						// TODO: Expand
						So(resp, ShouldNotBeNil)
					})

					Convey(`Then the error will be nil`, func() {
						So(err, ShouldBeNil)
					})
				})
			})

		})

		Convey(`And a server that returns a TooManyRequests response`, func() {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(http.StatusTooManyRequests)
			}))
			client.RemoteURL = srv.URL

			Convey(`And I have a request`, func() {
				req, err := client.newRequest(`/`, http.MethodGet, nil)
				if err != nil {
					panic(err)
				}

				Convey(`When I call the do method`, func() {
					resp, err := client.do(req)

					Convey(`Then the response will be nil`, func() {
						So(resp, ShouldBeNil)
					})

					Convey(`Then the Error will not be nil`, func() {
						// TODO: Expand
						So(err, ShouldNotBeNil)
					})
				})
			})

		})
	})
}
