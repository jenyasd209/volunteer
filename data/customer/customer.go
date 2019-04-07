package customer

import "time"

type Customer struct {
	ID         int
	Email      string
	Password   string
	Company_id int
	Phone      string
	Facebook   string
	Skype      string
	About      string
	Rait       float32
	Created_at time.Time
}
