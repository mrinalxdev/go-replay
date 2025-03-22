package db

import "time"

type Operation struct {
	Timestamp time.Time
	Command string
	Key string
	Value string
}
