package data

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

//CompleteOrder struct for "complete_order" table
type CompleteOrder struct {
	ID int						`json:"id"`
	Order						`json:"order"`
	Freelancer      			`json:"freelancer"`
	FreelancerComment Comment	`json:"freelancer_comment"`
	CustomerComment   Comment	`json:"customer_comment"`
	DateComplete      time.Time	`json:"date_complete"`
}

//MakeOrderDone row in "Order" table
func (customer *Customer) MakeOrderDone(completeOrder *CompleteOrder) (err error) {
	if completeOrder.Order.Status.ID == OrderStatusPerformed {
		if completeOrder.Customer.User.ID == customer.User.ID {
			statement := `INSERT INTO complete_orders (order_id, freelancer_id, customer_comment_id, date_complete)
									values ($1, $2, $3, $4)`
			stmt, err := Db.Prepare(statement)
			if err != nil {
				log.Println(err)
				return err
			}

			err = customer.CreateComment(&completeOrder.CustomerComment)
			if err != nil {
				log.Println(err)
			}
			defer stmt.Close()
			err = stmt.QueryRow(completeOrder.Order.ID, completeOrder.Freelancer.User.ID, completeOrder.CustomerComment.ID, time.Now()).Scan()
			if err != nil {
				log.Println(err)
			}
			completeOrder.Freelancer, _ = GetFreelancerByUserID(completeOrder.Freelancer.User.ID)
			completeOrder.Order.changeStatus(OrderStatusDone)
			customer.deletePerformedOrder(completeOrder.Order.ID)
			//customer.updateRating(completeOrder.FreelancerComment.Rait)

			completeOrder.Freelancer.calcRating(completeOrder.CustomerComment.Rait)
		} else {
			return errors.New("insufficient rights make order done")
		}
	}else {return errors.New("order is not preformed")}
	return
}

func (freelancer *Freelancer) FinishWorks() (completeOrders []CompleteOrder) {
	var freelancerCommentID sql.NullInt64
	rows, err := Db.Query(`SELECT id, order_id, freelancer_id, freelancer_comment_id, customer_comment_id,
       								date_complete FROM complete_orders 
									WHERE freelancer_id = $1 ORDER BY date_complete ASC `, freelancer.User.ID)
	if err != nil {
		return
	}
	for rows.Next() {
		completeOrder := CompleteOrder{}
		err = rows.Scan(&completeOrder.ID, &completeOrder.Order.ID, &completeOrder.Freelancer.User.ID,
			&freelancerCommentID, &completeOrder.CustomerComment.ID, &completeOrder.DateComplete)
		if err != nil {
			log.Println(err)
			return
		}
		if freelancerCommentID.Valid{
			completeOrder.FreelancerComment = GetCommentByID(int(freelancerCommentID.Int64))
		}
		completeOrder.CustomerComment = GetCommentByID(completeOrder.CustomerComment.ID)
		completeOrder.Order = GetOrderByID(completeOrder.Order.ID)
		completeOrders = append(completeOrders, completeOrder)
	}
	rows.Close()
	return
}

func (customer *Customer) CompleteOrders() (completeOrders []CompleteOrder) {
	rows, err := Db.Query(`SELECT id, order_id, freelancer_id, freelancer_comment_id, customer_comment_id,
       							  	date_complete
								  FROM complete_orders
								  WHERE
								  	order_id IN (SELECT id FROM orders WHERE customer_id = $1);`, customer.User.ID)
	if err != nil {
		return
	}
	var freelancerCommentID sql.NullInt64
	for rows.Next() {
		completeOrder := CompleteOrder{}
		err = rows.Scan(&completeOrder.ID, &completeOrder.Order.ID, &completeOrder.Freelancer.User.ID,
			&freelancerCommentID, &completeOrder.CustomerComment.ID, &completeOrder.DateComplete)
		if err != nil {
			log.Println(err)
		}
		if freelancerCommentID.Valid{
			completeOrder.FreelancerComment = GetCommentByID(int(freelancerCommentID.Int64))
		}
		completeOrder.Freelancer, _ = GetFreelancerByUserID(completeOrder.Freelancer.User.ID)
		completeOrder.CustomerComment = GetCommentByID(completeOrder.CustomerComment.ID)
		completeOrder.Order = GetOrderByID(completeOrder.Order.ID)
		completeOrders = append(completeOrders, completeOrder)
	}
	rows.Close()
	return
}

func GetDoneOrderByID(orderID int) (doneOrder CompleteOrder) {
	var freelancerCommentID sql.NullInt64
	err := Db.QueryRow(`SELECT id, order_id, freelancer_id, freelancer_comment_id, customer_comment_id
								  FROM complete_orders
								  WHERE order_id = $1;`, orderID).Scan(&doneOrder.ID, &doneOrder.Order.ID,
		&doneOrder.Freelancer.User.ID, &freelancerCommentID, &doneOrder.CustomerComment.ID)
	if err != nil {
		log.Println(err)
		return
	}
	doneOrder.FreelancerComment.ID = int(freelancerCommentID.Int64)
	doneOrder.Freelancer, _ = GetFreelancerByUserID(doneOrder.Freelancer.User.ID)
	doneOrder.Order = GetOrderByID(doneOrder.Order.ID)
	doneOrder.CustomerComment = GetCommentByID(doneOrder.CustomerComment.ID)
	return
}
