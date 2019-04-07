package user

type SessionHelper interface {
	SetUUID(string)
	GetUUID() string
	User() (UserHelper, error)
	Delete() error
	Check() (bool, error)
}
