package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"study/internal/client"
	"study/internal/models"
)

type CategoryService interface {
	GetCategories(ID int) ([]models.Category, error)
}
type CategoryAPIService struct {
	client client.Client
}

func NewCategoryAPIService(client client.Client) *CategoryAPIService {
	return &CategoryAPIService{client: client}
}

func (c *CategoryAPIService) GetCategories(ID int) ([]models.Category, error) {
	url := fmt.Sprintf("https://lenta.com/api-gateway/v1/catalog/categories?parentId=%d&depth=2", ID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("не удалось получить категории: %s, тело: %s", resp.Status, string(body))
	}

	var categoriesResp models.CategoriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&categoriesResp); err != nil {
		return nil, err
	}

	return categoriesResp.Categories, nil
}
