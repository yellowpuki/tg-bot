// Package telegram
package telegram

// UpdatesResponce ...
type UpdatesResponce struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// Update ...
type Update struct {
	ID      int              `json:"update_id"`
	Message *IncommingMessage `json:"message"`
}

// IncommingMessage ...
type IncommingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

// From ...
type From struct {
	UserName string `json:"username"`
}

// Chat ...
type Chat struct {
	ID int `json:"id"`
}
