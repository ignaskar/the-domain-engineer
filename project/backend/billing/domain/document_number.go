package domain

import (
	"errors"
	"fmt"
)

// TODO: define DocumentSeries and DocumentNumber
// types used by the billing domain.

var DocumentSeriesReceipt = DocumentSeries{"R"}

type DocumentSeries struct {
	series string
}

func NewDocumentSeries(series string) (DocumentSeries, error) {
	if series == "" {
		return DocumentSeries{}, errors.New("document series must not be empty")
	}

	return DocumentSeries{series: series}, nil
}

func (d DocumentSeries) IsZero() bool {
	return d.series == ""
}

func (d DocumentSeries) String() string {
	return d.series
}

type DocumentNumber struct {
	series DocumentSeries
	number int
}

func NewDocumentNumber(series DocumentSeries, number int) (DocumentNumber, error) {
	if series.IsZero() {
		return DocumentNumber{}, errors.New("document series must not be empty")
	}

	if number <= 0 {
		return DocumentNumber{}, errors.New("document number must be greater than zero")
	}

	return DocumentNumber{series: series, number: number}, nil
}

func (d DocumentNumber) IsZero() bool {
	return d.series.IsZero() && d.number == 0
}

func (d DocumentNumber) String() string {
	return fmt.Sprintf("%s-%08d", d.series, d.number)
}
