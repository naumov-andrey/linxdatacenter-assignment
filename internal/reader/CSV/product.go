package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/naumov-andrey/linxdatacenter-assignment/internal/model"
	"github.com/naumov-andrey/linxdatacenter-assignment/internal/reader"
)

type ProductCSVReader struct {
	reader *csv.Reader
}

func NewProductCSVReader(file *os.File) reader.ProductReader {
	reader := csv.NewReader(file)
	reader.Comma = ';'
	// skip header
	reader.Read()
	return &ProductCSVReader{reader}
}

func (r *ProductCSVReader) Read() (model.Product, error) {
	record, err := r.reader.Read()
	if err != nil {
		return model.Product{}, err
	}

	if len(record) != 3 {
		return model.Product{}, fmt.Errorf("record must contain 3 fields, got :%d", len(record))
	}

	price, err := strconv.Atoi(record[1])
	if err != nil {
		return model.Product{}, err
	}

	rating, err := strconv.Atoi(record[2])
	if err != nil {
		return model.Product{}, err
	}

	return model.Product{
		Name:   record[0],
		Price:  price,
		Rating: rating,
	}, nil
}
