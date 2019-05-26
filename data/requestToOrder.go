package data

import (
	"log"
	"time"
)

//FreelancerRequest struct for "requests" table
type FreelancerRequest struct {
	ID      int
	Comment string
	//CardOrder
	OrderID int
	Freelancer
}

//CreateRequest - create new request
func (freelancer *Freelancer) CreateRequest(orderID int, addText string) (err error) {
	statement := `INSERT INTO requests (freelancer_id, order_id, comment, created_at) values ($1, $2, $3, $4)`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(freelancer.ID, orderID, addText, time.Now()).Scan()
	return
	//statement := `INSERT INTO requests (freelancer_id, order_id, created_at) values ($1, $2, $3)
    //            returning id, freelancer_id, order_id`
	//stmt, err := Db.Prepare(statement)
	//if err != nil {
	//	return
	//}
	//defer stmt.Close()
	//err = stmt.QueryRow(request.Freelancer.ID, request.CardOrder.ID, time.Now()).Scan()
	//return
}

func (order *Order) GetRequests() (orderRequests *[]FreelancerRequest){
	var freelancerID int
	orderRequests = &[]FreelancerRequest{}
	rows, err := Db.Query(`SELECT id, freelancer_id, order_id, comment, created_at FROM requests 
									WHERE order_id = $1 ORDER BY created_at ASC `, order.ID)
	if err != nil {
		return
	}
	for rows.Next() {
		orderRequest := FreelancerRequest{}
		err = rows.Scan(&orderRequest.ID, &freelancerID, &orderRequest.OrderID, &orderRequest.Comment,
			&orderRequest.CreatedAt)
		if err != nil {
			log.Println(err)
			return
		}
		orderRequest.Freelancer, err = GetFreelancerByUserID(freelancerID)
		if err != nil {
			log.Println(err)
		}
		*orderRequests = append(*orderRequests, orderRequest)
	}
	if len(*orderRequests) == 0{
		orderRequests = nil
	}
	rows.Close()
	return
}

//FreelancerRequestDeleteAll - delete all rows in table "requests"
func FreelancerRequestDeleteAll() (err error) {
	statement := "delete from request"
	_, err = Db.Exec(statement)
	return
}
