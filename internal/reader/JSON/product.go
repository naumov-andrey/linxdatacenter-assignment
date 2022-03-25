package json

import (
	"encoding/json"
	"io"
	"os"

	"github.com/naumov-andrey/linxdatacenter-assignment/internal/model"
	"github.com/naumov-andrey/linxdatacenter-assignment/internal/reader"
)

type ProductJSONReader struct {
	decoder *json.Decoder
}

func NewProductJSONReader(file *os.File) reader.ProductReader {
	decoder := json.NewDecoder(file)
	// skip array begin token
	decoder.Token()
	return &ProductJSONReader{decoder}
}

func (r *ProductJSONReader) Read() (model.Product, error) {
	var result model.Product
	if !r.decoder.More() {
		return result, io.EOF
	}
	err := r.decoder.Decode(&result)
	return result, err
}
