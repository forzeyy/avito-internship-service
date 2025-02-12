package models

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type TransactionHistoryItem struct {
	FromUser string `json:"fromUser,omitempty"`
	ToUser   string `json:"toUser,omitempty"`
	Amount   int    `json:"amount"`
}

type CoinHistory struct {
	Received []TransactionHistoryItem `json:"received"`
	Sent     []TransactionHistoryItem `json:"sent"`
}

type InfoResponse struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}
