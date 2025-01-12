package dbo

type UserUnreadMessageCount struct {
	Name  string `db:"name"`
	Count int32  `db:"count"`
}
