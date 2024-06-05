package models

type Notification struct {
	SubscribingUserEmail string
	DaysBeforeNotify     int
	SubscribedUserEmail  string
}
