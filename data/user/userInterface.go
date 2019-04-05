package user

type User interface {
	Create() error
	Delete() error
	CreateSession() (Session, error)
}
