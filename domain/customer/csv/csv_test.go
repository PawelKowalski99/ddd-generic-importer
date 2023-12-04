package csv

import (
	"context"
	"testing"

	"github.com/PawelKowalski99/customerimporter/domain/customer"

	"github.com/stretchr/testify/require"
)

func TestCsv_GetCustomersByStream(t *testing.T) {
	type testCase struct {
		name            string
		customersAmount int
		expectedErr     error
	}

	testCases := []testCase{
		{
			name:            "Get Customers stream",
			customersAmount: 3005,
			expectedErr:     nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := CsvCustomerRepository{
				FileName: "./testfiles/customers.csv",
			}

			actualWorkerAmount := 0
			ch := make(chan *customer.Customer)
			go func() error {
				err := repo.FindCustomersByStream(context.Background(), ch)
				if err != nil {
					return err
				}
				return nil
			}()

			for cust := range ch {
				if cust != nil {
					actualWorkerAmount++
				}

			}
			require.Equal(t, tc.customersAmount, actualWorkerAmount)
		})
	}
}
