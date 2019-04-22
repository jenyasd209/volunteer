package data

import "time"

//FreelancerRequest struct for "requests" table
type FreelancerRequest struct {
	ID int
	Order
	Freelancer
}

//CreateRequest - create new request
func (freelancer *Freelancer) CreateRequest(request FreelancerRequest) (err error) {
	statement := `INSERT INTO requests (freelancer_id, order_id, created_at) values ($1, $2, $3)
                returning id, freelancer_id, order_id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(request.Freelancer.ID, request.Order.ID, time.Now()).Scan()
	return
}

//FreelancerRequestDeleteAll - delete all rows in table "requests"
func FreelancerRequestDeleteAll() (err error) {
	statement := "delete from request"
	_, err = Db.Exec(statement)
	return
}
