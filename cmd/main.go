package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/naumov-andrey/linxdatacenter-assignment/internal/model"
	"github.com/naumov-andrey/linxdatacenter-assignment/internal/reader"
	"github.com/naumov-andrey/linxdatacenter-assignment/internal/reader/csv"
	"github.com/naumov-andrey/linxdatacenter-assignment/internal/reader/json"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("there must be 1 command-line argument: path of file to read")
	}

	filePath := os.Args[1]
	fileExt := filepath.Ext(filePath)
	var readerConstructore func(file *os.File) reader.ProductReader

	switch fileExt {
	case ".csv":
		readerConstructore = csv.NewProductCSVReader
	case ".json":
		readerConstructore = json.NewProductJSONReader
	default:
		log.Fatalf("file extension must be '.csv' or '.json', got '%v'", fileExt)
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot open file: %w", err))
	}
	defer file.Close()

	records := make(chan model.Product)
	reader := readerConstructore(file)

	go func(records chan<- model.Product) {
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

	var maxPriceProduct, maxRatingProduct model.Product
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

	if isFirstComparing {
		fmt.Print("No product records found")
	} else {
		fmt.Printf(
			"Product with max price: %v \nProduct with max rating: %v",
			maxPriceProduct,
			maxRatingProduct,
		)
	}
}
