package models

// ErrNotAdmin represents an error when a user is not an admin.
type ErrNotAdmin struct {
}

func (e ErrNotAdmin) Error() string {
	return "User is not an admin"
}

// ErrAlreadyAdmin represents an error when a user is already an admin.
type ErrAlreadyAdmin struct {
}

func (e ErrAlreadyAdmin) Error() string {
	return "User is already an admin"
}




