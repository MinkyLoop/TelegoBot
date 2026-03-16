package main

import (
	"fmt"
	"study/internal/parser"
)

func main() {
	// открыть https://lenta.com/catalog/osobenno-vygodno-1893/
	products := parser.Pagination(1893)

	for _, v := range products {
		fmt.Println(v.Name, v.OldPrice, v.NewPrice)
	}
}
