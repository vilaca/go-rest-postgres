package models

type Login struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}
