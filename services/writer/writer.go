package writer

import (
	"io"
	"os"
)

type WriterService struct {
	writer io.Writer
}

type WriterServiceConfiguration func(*WriterService) error

func New(cfgs ...WriterServiceConfiguration) (*WriterService, error) {
	ws := &WriterService{}

	ws.writer = io.Discard

	for _, cfg := range cfgs {
		err := cfg(ws)
		if err != nil {
			return nil, err
		}
	}
	return ws, nil
}

func WithStdOutWriter() func(*WriterService) error {
	return func(ws *WriterService) error {
		ws.writer = io.MultiWriter(ws.writer, os.Stdout)
		return nil
	}
}

func WithFileWriter(fileName string) func(*WriterService) error {
	return func(ws *WriterService) error {
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			return err
		}

		ws.writer = io.MultiWriter(ws.writer, file)
		return nil
	}
}

func (w *WriterService) Save(a []byte) error {
	_, err := w.writer.Write(a)
	if err != nil {
		return err
	}

	return nil
}

// Why doesn't the csv Writer implement writer...
// func WithCsvFileWriter(csvFileWriter io.Writer) func(*WriterService) error {
// 	csvWr := csv.NewWriter(csvFileWriter)
// 	return func(ws *WriterService) error {
// 		ws.writer = io.MultiWriter(ws.writer, csvWr)
// 		return nil
// 	}
// }
