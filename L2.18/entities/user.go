package entities

import "time"

// User структура
type User struct {
	ID     int
	Events map[time.Time][]string
}
