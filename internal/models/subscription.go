package models

type Subscription struct {
	ID                int
	SubscribedUserID  int
	SubscribingUserID int
}
