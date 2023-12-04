package counter

import (
	"context"
	"io"
	"log/slog"
	"sort"
	"sync"

	"github.com/PawelKowalski99/customerimporter/domain/customer"

	"github.com/PawelKowalski99/customerimporter/services/worker"
)

type CounterService struct {
	customer   customer.CustomerRepository
	writer     io.Writer
	workerPool *worker.WorkerPoolService
	CountMap   *SafeMap
}

type SafeMap struct {
	Map map[string]map[string]bool
	mu  sync.Mutex
}

type CounterConfiguration func(*CounterService) error

func New(cfgs ...CounterConfiguration) (*CounterService, error) {
	cs := &CounterService{}

	// Default amount of workers in counterService worker pool is 1 -> basically no concurenncy
	cs.workerPool = worker.New(1)

	cs.CountMap = &SafeMap{Map: map[string]map[string]bool{}, mu: sync.Mutex{}}

	for _, cfg := range cfgs {
		err := cfg(cs)
		if err != nil {
			return nil, err
		}
	}

	// Start workers to wait for tasks
	cs.workerPool.Start()

	return cs, nil
}

func WithCustomerRepository(custRep customer.CustomerRepository) func(cs *CounterService) error {
	return func(cs *CounterService) error {
		cs.customer = custRep
		return nil
	}
}

func WithExternalWorkerPool(wp *worker.WorkerPoolService) func(cs *CounterService) error {
	return func(cs *CounterService) error {
		cs.workerPool = wp
		return nil
	}
}

func (cs *CounterService) CountDomainCustomers() error {

	// Find customers via stream

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	custChan := make(chan *customer.Customer)
	go func() error {
		err := cs.customer.FindCustomersByStream(ctx, custChan)
		if err != nil {
			return err
		}

		return nil
	}()
	for cust := range custChan {
		t, err := NewCounterTask(cust,
			WithInternalTaskVerifier(),
			WithExternalResultCounterTask(cs.CountMap),
		)
		if err != nil {
			return err
		}
		// t.Task()
		cs.workerPool.TaskQueue <- t
	}

	return nil

}

func (cs *CounterService) GetSortedUserDomainCount() []UserDomainCount {

	m := []UserDomainCount{}
	for domain, sliceOfEmails := range cs.CountMap.Map {
		m = append(m, UserDomainCount{Domain: domain, Amount: len(sliceOfEmails)})
	}

	sort.SliceStable(m, func(i, j int) bool {
		return m[i].Amount > m[j].Amount
	})

	return m
}

type CounterTask struct {
	customer   *customer.Customer
	verifier   Verifier
	resultChan *SafeMap
}
type CounterTaskConfiguration func(*CounterTask) error

func WithExternalTaskVerifier() func(cs *CounterTask) error {
	return func(cs *CounterTask) error {
		cs.verifier = &externalVerifier{}
		return nil
	}
}
func WithInternalTaskVerifier() func(cs *CounterTask) error {
	return func(cs *CounterTask) error {
		cs.verifier = &internalVerifier{}
		return nil
	}
}

func WithExternalResultCounterTask(m *SafeMap) func(cs *CounterTask) error {
	return func(cs *CounterTask) error {
		cs.resultChan = m
		return nil
	}
}

func NewCounterTask(cust *customer.Customer, cfgs ...CounterTaskConfiguration) (*CounterTask, error) {
	ct := &CounterTask{
		resultChan: &SafeMap{Map: map[string]map[string]bool{}, mu: sync.Mutex{}},
		verifier:   &internalVerifier{},
		customer:   cust,
	}

	for _, cfg := range cfgs {
		err := cfg(ct)
		if err != nil {
			return nil, err
		}

	}

	return ct, nil

}

func (t *CounterTask) Task() {
	res, err := t.verifier.Verify(t.customer.Email)
	if err != nil {
		slog.Error("Could not verify", slog.Any("error", err))
		return
	}

	t.resultChan.mu.Lock()
	defer t.resultChan.mu.Unlock()

	// Check if key was already populated
	if _, valid := t.resultChan.Map[res.Domain]; !valid {
		t.resultChan.Map[res.Domain] = map[string]bool{}
	}
	t.resultChan.Map[res.Domain][res.Email] = true

	// TODO: Check for type and make error assertion
}

type UserDomainCount struct {
	Domain string
	Amount int
}
