package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Product struct {
	Name   string `json:"prodcut" csv:"name"`
	Price  int    `json:"price" csv:"price"`
	Rating int    `json:"rating" csv:"rating"`
}

func (p Product) String() string {
	return fmt.Sprintf("product: %s price: %d rating: %d", p.Name, p.Price, p.Rating)
}

type ProductReader interface {
	Read() (Product, error)
}

type ProductCSVReader struct {
	reader *csv.Reader
}

func NewProductCSVReader(file *os.File) ProductReader {
	reader := csv.NewReader(file)
	reader.Comma = ';'
	// skiping header
	reader.Read()
	return &ProductCSVReader{reader}
}

func (r *ProductCSVReader) Read() (Product, error) {
	record, err := r.reader.Read()
	if err != nil {
		return Product{}, err
	}

	if len(record) != 3 {
		return Product{}, fmt.Errorf("record must contain 3 fields, got :%d", len(record))
	}

	price, err := strconv.Atoi(record[1])
	if err != nil {
		return Product{}, err
	}

	rating, err := strconv.Atoi(record[2])
	if err != nil {
		return Product{}, err
	}

	return Product{record[0], price, rating}, nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("there must be 1 command-line argument: path of file to read")
	}

	filePath := os.Args[1]
	fileExt := filepath.Ext(filePath)
	var readerConstructore func(file *os.File) ProductReader
	switch fileExt {
	case ".csv":
		readerConstructore = NewProductCSVReader
	case ".json":
		// TODO: implement json reader
		readerConstructore = NewProductCSVReader
	default:
		log.Fatalf("file extension must be '.csv' or '.json', got '%v'", fileExt)
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot open file: %w", err))
	}

	records := make(chan Product)
	reader := readerConstructore(file)

	go func(records chan<- Product) {
		defer close(records)

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(fmt.Errorf("error while reading: %w", err))
			}

			records <- record
		}
	}(records)

	var maxPriceProduct, maxRatingProduct Product
	isFirstComparing := true

	for record := range records {
		if isFirstComparing {
			maxPriceProduct = record
			maxRatingProduct = record
			isFirstComparing = false
		}

		if record.Price > maxPriceProduct.Price {
			maxPriceProduct = record
		}
		if record.Rating > maxRatingProduct.Rating {
			maxRatingProduct = record
		}
	}

	fmt.Printf(
		"Product with max price: %v \nProduct with max rating: %v",
		maxPriceProduct,
		maxRatingProduct,
	)
}
