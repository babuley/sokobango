package main

import (
	"github.com/google/uuid"
)

type Player struct {
	X, Y int
}

type Target struct {
	X, Y int
	ID   uuid.UUID
}

type Boulder struct {
	X, Y int
	ID   uuid.UUID
}
