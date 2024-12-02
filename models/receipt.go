package models

// External receipt structure sent by client
type ExtReceipt struct {
	Retailer     string `json:"retailer" binding:"required"`
	PurchaseDate string `json:"purchaseDate" binding:"required"`
	PurchaseTime string `json:"purchaseTime" binding:"required"`
	Items        []Item `json:"items" binding:"required,dive"`
	Total        string `json:"total" binding:"required"`
}

// Internal receipt structure used internally
type Receipt struct {
	ID           string
	Retailer     string
	PurchaseDate string
	PurchaseTime string
	Items        []Item
	Total        string
}

// A single item purchased in a receipt
type Item struct {
	ShortDescription string `json:"shortDescription" binding:"required"`
	Price            string `json:"price" binding:"required"`
}
