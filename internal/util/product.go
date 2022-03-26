package util

import (
	"fmt"
	"io"
	"log"

	"github.com/naumov-andrey/linxdatacenter-assignment/internal/model"
	"github.com/naumov-andrey/linxdatacenter-assignment/internal/reader"
)

func ParseProducts(records chan<- model.Product, reader reader.ProductReader) {
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
}

func ProcessProducts(records <-chan model.Product) (model.Product, model.Product, bool) {
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

	return maxPriceProduct, maxRatingProduct, !isFirstComparing
}
