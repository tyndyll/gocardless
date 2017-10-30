package gocardless

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	customerEndpoint = `/customers`
)

// customerWrapper is a utility struct used to wrap and unwrap the JSON request being passed to the remote API
type customerWrapper struct {
	Customer *Customer `json:"customers"`
}

//
func (c *Client) CreateCustomer(customer *Customer) error {
	custRequest := &customerWrapper{customer}

	data, err := json.Marshal(custRequest)
	if err != nil {
		return err
	}

	req, err := c.newRequest(customerEndpoint, http.MethodPost, data)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("%s WAS NOT CREATED", c.decodeError(respBody).Error()))
	}

	if err := json.Unmarshal(respBody, custRequest); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetCustomer(id string) (*Customer, error) {
	request, err := c.newRequest(fmt.Sprintf(`%s/%s`, customerEndpoint, id), http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.do(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	wrapper := &customerWrapper{}
	if err := json.Unmarshal(body, wrapper); err != nil {
		return nil, err
	}

	return wrapper.Customer, nil
}

func (c *Client) ListCustomer() ([]*Customer, error) {
	req, err := c.newRequest(customerEndpoint, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	customerList := []*Customer{}
	if err := json.Unmarshal(data, customerList); err != nil {
		return nil, err
	}

	return customerList, nil
}

func (c *Client) UpdateCustomer(customer *Customer) error {
	custRequest := &customerWrapper{customer}

	data, err := json.Marshal(custRequest)
	if err != nil {
		return err
	}

	req, err := c.newRequest(fmt.Sprintf("%s/%s", customerEndpoint, customer.ID), http.MethodPut, data)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(respBody, custRequest); err != nil {
		return err
	}
	return nil
}
