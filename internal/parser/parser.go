package parser

import (
	"fmt"
	"study/internal/api"
	"study/internal/models"
	"study/internal/storage"
	"time"
)

type Parser struct {
	itemService     api.ItemService
	categoryService api.CategoryService
	storage         storage.Storage
	rootID          int
}

func NewParser(itemService api.ItemService,
	categoryService api.CategoryService,
	storage storage.Storage,
	rootID int,
) *Parser {
	return &Parser{
		itemService:     itemService,
		categoryService: categoryService,
		storage:         storage,
		rootID:          rootID,
	}
}

func (p *Parser) Run() ([]models.Item, error) {
	fmt.Printf("Получение категорий из parentID => %d\n", p.rootID)
	categories, err := p.categoryService.GetCategories(p.rootID)
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения категорий %w", err)
	}

	fmt.Println("Получено ", len(categories), " категорий")

	var allItems []models.Item

	for _, category := range categories {
		items, err := p.itemService.GetAllItems(category.ID)
		if err != nil {
			continue
		}

		allItems = append(allItems, items...)
		time.Sleep(1 * time.Second)
	}

	if err := p.storage.SaveItems(allItems); err != nil {
		return nil, err
	}

	return allItems, nil
}
