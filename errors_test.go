package gocardless

import (
	"testing"

	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
)

func TestErrorUnmarshal(t *testing.T) {
	Convey(`Given I have a GoCardless error message`, t, func() {
		msg := []byte(`{
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
				}`)

		Convey(`And I have a Error instance`, func() {
			err := &Error{}

			Convey(`When I call Unmarshal`, func() {
				if unmarshalErr := json.Unmarshal(msg, err); unmarshalErr != nil {
					panic(err)
				}

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
	})
}

func TestErrorDetailUnmarshal(t *testing.T) {
	Convey(`Given I have a GoCardless error detail message`, t, func() {
		msg := []byte(`{
					"message": "must be a number",
					"field": "branch_code",
					"request_pointer": "/customer_bank_accounts/branch_code"
				}`)

		Convey(`And I have a Error instance`, func() {
			err := &ErrorDetail{}

			Convey(`When I call Unmarshal`, func() {
				if unmarshalErr := json.Unmarshal(msg, err); unmarshalErr != nil {
					panic(err)
				}

				Convey(`Then the Message field will be populated correctly`, func() {
					So(err.Message, ShouldEqual, `must be a number`)
				})

				Convey(`Then the Field field will be populated correctly`, func() {
					So(err.Field, ShouldEqual, `branch_code`)
				})

				Convey(`Then the RequestPointer field will be populated correctly`, func() {
					So(err.RequestPointer, ShouldEqual, `/customer_bank_accounts/branch_code`)
				})
			})
		})
	})
}
