package importer

import (
	"encoding/json"

	"github.com/PawelKowalski99/customerimporter/services/counter"
	"github.com/PawelKowalski99/customerimporter/services/worker"
	"github.com/PawelKowalski99/customerimporter/services/writer"
)

type ImporterService struct {
	Counter    *counter.CounterService
	Writer     *writer.WriterService
	WorkerPool *worker.WorkerPoolService
}

type ImporterServiceConfiguration func(*ImporterService) error

func New(cfgs ...ImporterServiceConfiguration) (*ImporterService, error) {
	is := &ImporterService{}

	cr, err := counter.New()
	if err != nil {
		return nil, err
	}

	is.Counter = cr

	wr, err := writer.New()
	if err != nil {
		return nil, err
	}
	is.Writer = wr

	for _, cfg := range cfgs {
		err := cfg(is)
		if err != nil {
			return nil, err
		}
	}
	return is, nil
}

func WithCounterService(cs *counter.CounterService) func(*ImporterService) error {
	return func(is *ImporterService) error {
		is.Counter = cs
		return nil
	}
}

func WithWriterService(ws *writer.WriterService) func(*ImporterService) error {
	return func(is *ImporterService) error {
		is.Writer = ws
		return nil
	}
}

func (i *ImporterService) ImportUserDomainCount() error {

	err := i.Counter.CountDomainCustomers()
	if err != nil {
		return err
	}

	usc := i.Counter.GetSortedUserDomainCount()

	// Well... We could always include different encoders json,yaml,csv etc. For now calling directly json library.

	parsedUsc, err := json.Marshal(&usc)
	if err != nil {
		return err
	}

	err = i.Writer.Save(parsedUsc)
	if err != nil {
		return err
	}

	return nil
}
