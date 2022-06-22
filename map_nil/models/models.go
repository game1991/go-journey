package models

import "pkg.deepin.com/golang/lib/uuid"

type Game struct {
	ID   string
	Name string
	Age  int
}

func NewGame() *Game {
	return &Game{ID: uuid.UUID32(), Name: "wechat", Age: 11}
}
