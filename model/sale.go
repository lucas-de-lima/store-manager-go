package model

type SaleProduct struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

type Sale struct {
	ID          int           `json:"id"`
	Date        string        `json:"date,omitempty"`
	ItemsSold   []SaleProduct `json:"itemsSold,omitempty"`
	ProductID   int           `json:"productId,omitempty"`
	Quantity    int           `json:"quantity,omitempty"`
	SaleID      int           `json:"saleId,omitempty"`
}