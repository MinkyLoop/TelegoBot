package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"study/internal/client"
	"study/internal/models"
	"time"
)

type ItemService interface {
	GetAllItems(ID int) ([]models.Item, error)
}
type ItemAPIService struct {
	client client.Client
}

func NewItemAPIService(client client.Client) *ItemAPIService {
	return &ItemAPIService{client: client}
}

func (c *ItemAPIService) GetAllItems(ID int) ([]models.Item, error) {
	const limit = 40
	offset := 0
	var allItems []models.Item
	url := "https://lenta.com/api-gateway/v1/catalog/items"

	for {
		reqBody := models.ItemsRequest{
			CategoryID: ID,
			Filters:    models.Filters{Checkbox: []interface{}{}, Multicheckbox: []interface{}{}, Range: []interface{}{}},
			Sort:       models.Sort{Type: "popular", Order: "desc"},
			Limit:      limit,
			Offset:     offset,
		}
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("статус %s, тело: %s", resp.Status, string(body))
		}

		var respData models.ItemsResponse
		if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			return nil, err
		}

		for i := range respData.Items {
			respData.Items[i].CategoryID = ID
		}

		allItems = append(allItems, respData.Items...)

		if len(respData.Items) < limit {
			break
		}

		offset += limit

		time.Sleep(500 * time.Millisecond)
	}

	return allItems, nil
}
