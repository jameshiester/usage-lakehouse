package dbentity

import "time"

type Account struct {
	ID      string
	LegalID string
	Name    string
	Created time.Time
	Updated time.Time
}
