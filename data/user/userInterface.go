package user

type Userable interface {
	Create() error
	Delete() error
	CreateSession() (Sessionable, error)
}
