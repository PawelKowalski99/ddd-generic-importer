// Package github.com/PawelKowalski99/customerimporter reads from a CSV file and returns a sorted (data
// structure of your choice) of email domains along with the number of customers
// with e-mail addresses for each domain. This should be able to be ran from the
// CLI and output the sorted domains to the terminal or to a file. Any errors
// should be logged (or handled). Performance matters (this is only ~3k lines,
// but could be 1m lines or run on a small machine).
package main

import (
	"log/slog"
	"os"
	"time"

	"net/http"
	_ "net/http/pprof"

	importer "github.com/PawelKowalski99/customerimporter"
	customer_csv "github.com/PawelKowalski99/customerimporter/domain/customer/csv"
	"github.com/PawelKowalski99/customerimporter/services/counter"
	"github.com/PawelKowalski99/customerimporter/services/worker"
	"github.com/PawelKowalski99/customerimporter/services/writer"
)

// Assume this could be set up somewhere in env files etc.
var (

	// 12 * 2^15 * 2kB == 768MB of goroutines -> on my machine more slowed down program due to jumping.
	// workersCount = int(math.Exp2(1))
	workersCount = 300

	csvFile string
)

const (
	defaultLogLevel slog.Level = slog.LevelInfo
	defaultFile     string     = "../customers2.csv"
)

var (
	globalWorkerPool *worker.WorkerPoolService
)

func init() {
	go func() {
		http.ListenAndServe(":1234", nil)
	}()
}
func main() {

	timer := time.Now()

	// TODO: Config interface
	// GetLogLevel -> default/env
	// GetWorkersAmount -> default/env

	// Setup logger
	// Assume we are setting up some logger configs and then
	// slog.SetDefault(customLogger)

	// Setup global customers file parameter
	csvFile = getCSVFileName()

	globalWorkerPool = worker.New(workersCount)

	customerRepository := customer_csv.New(csvFile)

	counterService, err := counter.New(
		counter.WithCustomerRepository(customerRepository),
		counter.WithExternalWorkerPool(globalWorkerPool),
	)
	if err != nil {
		slog.Error("error ocurred when init counter", slog.Any("err", err))
		return
	}

	writerService, err := writer.New(
		writer.WithStdOutWriter(),
		writer.WithFileWriter("sorted_output.json"),
	)
	if err != nil {
		slog.Error("error ocurred when init writer", slog.Any("err", err))
		return
	}

	imp, err := importer.New(
		importer.WithCounterService(counterService),
		importer.WithWriterService(writerService),
	)
	if err != nil {
		slog.Error("error occured when init importer", slog.Any("err", err))
		return
	}

	imp.ImportUserDomainCount()

	slog.Info("program ended", slog.Any("finishTime", time.Now().Sub(timer)))
}

func getCSVFileName() string {

	if len(os.Args) == 2 && os.Args[1] != "" {
		return os.Args[1]
	}

	return defaultFile

}

func setupGlobalLogger() {
	// logger := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: getLogLevel()})
}
