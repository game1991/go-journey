package models

import (
	"a/uuid"
)

type Game struct {
	ID   string
	Name string
	Age  int
}

func NewGame() *Game {
	return &Game{ID: uuid.UUID32(), Name: "wechat", Age: 11}
}
