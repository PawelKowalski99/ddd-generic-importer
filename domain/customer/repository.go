package customer

import "context"

// first_name,last_name,email,gender,ip_address
type Customer struct {
	FirstName string
	LastName  string
	Email     string
	Gender    string
	IpAddress string
}

type CustomerRepository interface {
	FindCustomersByStream(context.Context, chan *Customer) error
}
