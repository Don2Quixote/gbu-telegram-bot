package entity

// User is struct that represents telegram user
type User struct {
	// ID is telegram user's ID
	ID int64
	// Username is telegram user's @username (without "@")
	// It is empty string ("") if user hasn't username
	Username string
	// Name is telegram user's first and last names combined
	// into single string
	Name string
	// IsSubscribed is flag that indicates should user get
	// notifications about new posts is blog or not
	IsSubscribed bool
}
