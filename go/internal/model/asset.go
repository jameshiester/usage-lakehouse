package model

import "time"

type Asset struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	AssetCode string    `json:"asset_code"`
	AccountID string    `json:"account_id"`
	Created   time.Time `json:"created_dttm"`
	Updated   time.Time `json:"updated_dttm"`
}
