package models

type NotifySetting struct {
	ID               int `json:"id"`
	UserID           int `json:"user_id"`
	DaysBeforeNotify int `json:"days_before_notify"`
}
