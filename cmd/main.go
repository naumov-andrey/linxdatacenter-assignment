package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/naumov-andrey/linxdatacenter-assignment/internal/model"
	"github.com/naumov-andrey/linxdatacenter-assignment/internal/reader"
	"github.com/naumov-andrey/linxdatacenter-assignment/internal/reader/csv"
	"github.com/naumov-andrey/linxdatacenter-assignment/internal/reader/json"
	"github.com/naumov-andrey/linxdatacenter-assignment/internal/util"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("there must be 1 command-line argument: path of file to read")
	}

	filePath := os.Args[1]
	fileExt := filepath.Ext(filePath)
	var readerConstructor func(file *os.File) reader.ProductReader

	switch fileExt {
	case ".csv":
		readerConstructor = csv.NewProductCSVReader
	case ".json":
		readerConstructor = json.NewProductJSONReader
	default:
		log.Fatalf("file extension must be '.csv' or '.json', got '%v'", fileExt)
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot open file: %w", err))
	}
	defer file.Close()

	products := make(chan model.Product)

	go util.ParseProducts(products, readerConstructor(file))

	maxPrice, maxRating, ok := util.ProcessProducts(products)
	if !ok {
		log.Fatal("No product records found")
	}

	fmt.Printf("Product with max price: %v \nProduct with max rating: %v", maxPrice, maxRating)
}
