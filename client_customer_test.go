package gocardless

import (
	"testing"

	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
)

func TestClientCreateCustomer(t *testing.T) {
	Convey(`Given I have a client`, t, func() {
		client := &Client{}

		Convey(`And I have a server which returns a valid response`, func() {
			var requestMethod string
			var requestPath string

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				requestMethod = req.Method
				requestPath = req.URL.Path
			}))

			client.RemoteURL = srv.URL

			Convey(`And I have a valid customer`, func() {
				customer := &Customer{}

				Convey(`When I call the CreateCustomer method`, func() {
					err := client.CreateCustomer(customer)

					Convey(`Then the request method will be POST`, func() {
						So(requestMethod, ShouldEqual, http.MethodPost)
					})

					Convey(`Then the URL will use the customers endpoint`, func() {
						So(requestPath, ShouldEqual, customerEndpoint)
					})

					Convey(`Then the error will be nil`, func() {
						So(err, ShouldBeNil)
					})

					Convey(`Then the customer ID will be populated`, func() {
						So(customer.ID, ShouldEqual, `CU001`)
					})
				})
			})
		})
	})
}

func TestClientGetCustomer(t *testing.T) {
	Convey(`Given I have a client`, t, func() {
		client := &Client{}

		Convey(`And I have a server which returns a valid response`, func() {
			var requestMethod string
			var requestPath string

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				requestMethod = req.Method
				requestPath = req.URL.Path

				w.Write([]byte(`{
									"customers": {
										"id": "CU123",
										"created_at": "2014-05-08T17:01:06.000Z",
										"email": "user@example.com",
										"given_name": "Frank",
										"family_name": "Osborne",
										"address_line1": "27 Acer Road",
										"address_line2": "Apt 2",
										"address_line3": null,
										"city": "London",
										"region": null,
										"postal_code": "E8 3GX",
										"country_code": "GB",
										"language": "en",
										"swedish_identity_number": null,
										"metadata": {
											"salesforce_id": "ABCD1234"
										}
									}
								}`))
			}))

			client.RemoteURL = srv.URL

			Convey(`And I have a valid Customer ID`, func() {
				customerID := `CU123`

				Convey(`When I call the GetCustomer method`, func() {
					customer, err := client.GetCustomer(customerID)

					Convey(`Then the request method will be GET`, func() {
						So(requestMethod, ShouldEqual, http.MethodGet)
					})

					Convey(`Then the URL will use the customers endpoint and customer ID`, func() {
						So(requestPath, ShouldEqual, fmt.Sprintf("%s/%s", customerEndpoint, customerID))
					})

					Convey(`Then the error will be nil`, func() {
						So(err, ShouldBeNil)
					})

					Convey(`Then the customer ID will be populated`, func() {
						So(customer.ID, ShouldEqual, customerID)
					})
				})
			})
		})
	})
}
