package models

type Purchase struct {
	User    string `json:"username"`
	Product string `json:"product"`
}
