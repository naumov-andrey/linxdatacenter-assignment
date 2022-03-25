package model

import "fmt"

type Product struct {
	Name   string `json:"product"`
	Price  int    `json:"price"`
	Rating int    `json:"rating"`
}

func (p Product) String() string {
	return fmt.Sprintf("product: %s, price: %d, rating: %d", p.Name, p.Price, p.Rating)
}
