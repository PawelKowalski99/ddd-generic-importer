package csv

import (
	"context"
	enc_csv "encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/PawelKowalski99/customerimporter/domain/customer"
)

type CsvCustomerRepository struct {
	FileName string
}

func New(fileName string) *CsvCustomerRepository {
	return &CsvCustomerRepository{
		FileName: fileName,
	}
}

func (c *CsvCustomerRepository) FindCustomersByStream(ctx context.Context, cst chan *customer.Customer) error {
	// Open customers persistance
	f, err := os.Open(c.FileName)
	if err != nil {
		return err
	}

	lineCount := 0
	reader := enc_csv.NewReader(f)
	// Skip first line
	// first_name,last_name,email,gender,ip_address
	_, err = reader.Read()
	for {
		line, err := reader.Read()
		if err == io.EOF {
			fmt.Println("closed cst")
			close(cst)

			break
		}

		cust := &customer.Customer{
			Email:     line[2],
			FirstName: line[0],
			LastName:  line[1],
			Gender:    line[3],
			IpAddress: line[4],
		}
		select {
		case <-ctx.Done():

		case cst <- cust:
		}
		lineCount++
	}
	return nil
}
