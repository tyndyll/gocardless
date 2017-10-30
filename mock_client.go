package gocardless

// MockClient is provided to assist with your testing. It implements the API interface.
//
// Each field in the struct corresponds to the associated method, enabling you to implement whatever functionality is
// necessary in your test e.g. Implementing a customer not found using GetCustomer
//
//     mock := &MockClient{}
//     mock.GetCustomerFunc = func(_ string) (*Customer, error) {
//         return nil, errors.New("No customer")
//     }
//     // Calls the GetCustomerFunc value. Note that this will panic if this value has not been set
//     mock.GetCustomer(`random-id`)
//
type MockClient struct {
	CreateCustomerFunc func(*Customer) error
	GetCustomerFunc    func(string) (*Customer, error)
	ListCustomerFunc   func() ([]*Customer, error)
}

func (mock *MockClient) CreateCustomer(c *Customer) error {
	return mock.CreateCustomerFunc(c)
}

func (mock *MockClient) GetCustomer(id string) (*Customer, error) {
	return mock.GetCustomerFunc(id)
}

func (mock *MockClient) ListCustomer() ([]*Customer, error) {
	return mock.ListCustomerFunc()
}
