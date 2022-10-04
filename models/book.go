package models

type Book struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Author    string  `json:"author"`
	Quantity  float64 `json:"quantity"`
	Available string  `json:"available"`
}

type User struct {
	ID              string `json:"id"`
	CheckedoutBooks []Book `json:"checkedoutBooks"`
}
