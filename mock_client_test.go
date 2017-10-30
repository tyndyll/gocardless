package gocardless

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMockClientCreateCustomer(t *testing.T) {
	Convey(`Given I have a MockClient`, t, func() {
		client := &MockClient{}

		Convey(`And I have a function to mock CreateCustomer`, func() {
			isCalled := false

			client.CreateCustomerFunc = func(customer *Customer) error {
				isCalled = true
				return nil
			}

			Convey(`When I call CreateCustomer`, func() {
				client.CreateCustomer(&Customer{})

				Convey(`Then the mock function is called`, func() {
					So(isCalled, ShouldBeTrue)
				})
			})
		})
	})
}

func TestMockClientGetCustomer(t *testing.T) {
	Convey(`Given I have a MockClient`, t, func() {
		client := &MockClient{}

		Convey(`And I have a function to mock GetCustomer`, func() {
			isCalled := false

			client.GetCustomerFunc = func(id string) (*Customer, error) {
				isCalled = true
				return nil, nil
			}

			Convey(`When I call GetCustomer`, func() {
				custID := `CU123`
				client.GetCustomer(custID)

				Convey(`Then the mock function is called`, func() {
					So(isCalled, ShouldBeTrue)
				})
			})
		})
	})
}

func TestMockClientListCustomer(t *testing.T) {
	Convey(`Given I have a MockClient`, t, func() {
		client := &MockClient{}

		Convey(`And I have a function to mock ListCustomer`, func() {
			isCalled := false

			client.ListCustomerFunc = func() ([]*Customer, error) {
				isCalled = true
				return nil, nil
			}

			Convey(`When I call ListCustomer`, func() {
				client.ListCustomer()

				Convey(`Then the mock function is called`, func() {
					So(isCalled, ShouldBeTrue)
				})
			})
		})
	})
}
