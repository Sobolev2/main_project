package models

type Friend struct {
	ID         int `json:"id"`
	UserFirst  int `json:"user_first"`
	UserSecond int `json:"user_second"`
}
