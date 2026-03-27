package models

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ParentID    int    `json:"parentId"`
	HasChildren bool   `json:"hasChildren"`
}

type CategoriesResponse struct {
	Categories []Category `json:"categories"`
}

type ItemsRequest struct {
	CategoryID int     `json:"categoryId"`
	Filters    Filters `json:"filters"`
	Sort       Sort    `json:"sort"`
	Limit      int     `json:"limit"`
	Offset     int     `json:"offset"`
}

type Filters struct {
	Checkbox      []interface{} `json:"checkbox"`
	Multicheckbox []interface{} `json:"multicheckbox"`
	Range         []interface{} `json:"range"`
}

type Sort struct {
	Type  string `json:"type"`
	Order string `json:"order"`
}

type ItemsResponse struct {
	Items []Item `json:"items"`
	Total int    `json:"total"`
}

type Item struct {
	ID         int    `json:"id"`
	Title      string `json:"name"`
	Prices     Prices `json:"prices"`
	CategoryID int    `json:"categoryId"`
}

type Prices struct {
	Cost         int `json:"cost"`
	Price        int `json:"price"`
	PriceRegular int `json:"priceRegular"`
}
