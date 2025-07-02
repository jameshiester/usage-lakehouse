package model

import "time"

type Account struct {
	ID      string    `json:"id"`
	LegalID *string   `json:"legal_id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created_dttm"`
	Updated time.Time `json:"updated_dttm"`
}
