package jsonstruct

type Request struct {
	CategoryID int     `json:"categoryId"`
	Filters    Filters `json:"filters"`
	Limit      int     `json:"limit"`
	Offset     int     `json:"offset"`
	Sort       Sort    `json:"sort"`
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

type Response struct {
	Categories []Categoties `json:"categories"`
}
type Categoties struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type JsonProduct struct {
	Items []Items `json:"items"`
}
type Items struct {
	Name   string `json:"name"`
	Prices Prices
}
type Prices struct {
	Price        float64 `json:"price"`
	PriceRegular float64 `json:"priceRegular"`
}
type Product struct {
	Name     string
	OldPrice float64
	NewPrice float64
}
