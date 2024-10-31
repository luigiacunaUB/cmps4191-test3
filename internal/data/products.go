package data

import (
	"time"
)

type Product struct {
	ID        int64
	prodName  string
	addedDate time.Time
}
