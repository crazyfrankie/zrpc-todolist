package entity

type User struct {
	UserID int64

	Name    string // unique name
	IconURI string // avatar URI
	IconURL string // avatar URL

	CreatedAt int64 // creation time
	UpdatedAt int64 // update time
}
