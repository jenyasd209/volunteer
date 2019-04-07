package user

type UserDBHelper interface {
	Create() error
	Delete() error
	Update() error
}

type UserSessionHelper interface {
	CreateSession() (SessionHelper, error)
}

type UserHelper interface {
	UserDBHelper
	UserSessionHelper
}
