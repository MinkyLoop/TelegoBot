package main

import (
	"study/internal/pkg/parser"
)

func main() {
	// открыть https://online.metro-cc.ru/category/vse_skidki-40813

	parser.AllCategoricalProductParse(parser.CategoryParse("https://online.metro-cc.ru/category/vse_skidki-40813"))
}
