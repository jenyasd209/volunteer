package customer

import "time"

type Customer struct {
	ID         int
	Email      string
	Password   string
	Company_id int
	Created_at time.Time
}
