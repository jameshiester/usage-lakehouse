package model

import "time"

type PowerRegion struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created_dttm"`
	Updated time.Time `json:"updated_dttm"`
}
