package data

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"time"
)

//CardOrder struct for "orders" table
type Order struct {
	ID      int								`json:"id"`
	Title   string							`json:"title"`
	Content string							`json:"content"`
	//CustomerID int
	Customer								`json:"customer"`
	Status									`json:"status"`
	SpecializationID int 					`json:"specialization_id"`
	CreatedAt time.Time						`json:"created_at"`
	FreelancerRequest *[]FreelancerRequest	`json:"freelancer_request"`
}

//PerformedOrder struct for "performed_order" table
type PerformedOrder struct {
	ID int				`json:"id"`
	Order				`json:"order"`
	Freelancer			`json:"freelancer"`
}

//CompleteOrder struct for "complete_order" table
type CompleteOrder struct {
	ID int						`json:"id"`
	Order						`json:"order"`
	Freelancer      			`json:"freelancer"`
	FreelancerComment Comment	`json:"freelancer_comment"`
	CustomerComment   Comment	`json:"customer_comment"`
	DateComplete      time.Time	`json:"date_complete"`
}

//Status struct for "order_status" table
type Status struct {
	ID   int
	Name string
}

//CreateOrder new row in "Order" table
func (customer *Customer) CreateOrder(order *Order) (err error) {
	statement := `insert into orders (title, content, customer_id, specialization_id, created_at)
								values ($1, $2, $3, $4, $5) returning id, title, content, customer_id, specialization_id, created_at`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err)
	}

	defer stmt.Close()
	err = stmt.QueryRow(order.Title, order.Content, customer.User.ID, order.SpecializationID, time.Now()).Scan(&order.ID, &order.Title,
		&order.Content, &order.Customer.User.ID, &order.SpecializationID, &order.CreatedAt)
	if err != nil {
		log.Println(customer.ID)
		log.Println(err)
	}
	return
}

//UpdateOrder row in "Order" table
func (customer *Customer) UpdateOrder(order *Order) (err error) {
	if order.Customer.User.ID == customer.User.ID {
		if order.IsAvailable() {
			statement := `UPDATE orders SET title = $1, content = $2, specialization_id = $3 WHERE id = $4
			returning id`
			stmt, err := Db.Prepare(statement)
			if err != nil {
				return err
			}

			defer stmt.Close()
			stmt.QueryRow(order.Title, order.Content, order.SpecializationID, order.ID).Scan(&order.ID)
		} else {
			return errors.New("can't change Order, it status not 'Available'")
		}
	} else {
		return errors.New("insufficient rights to update Order")
	}
	return
}

// DeleteOrder row from "Order" table
func (customer *Customer) DeleteOrder(order Order) (err error) {
	if order.Customer.User.ID == customer.User.ID {
		if order.IsAvailable() {
			statement := "DELETE FROM orders WHERE id = $1"
			stmt, err := Db.Prepare(statement)
			if err != nil {
				return err
			}

			defer stmt.Close()
			_, err = stmt.Exec(order.ID)
			if err != nil {
				return err
			}
		} else {
			return errors.New("can't change Order, it status not 'Available'")
		}
	} else {
		return errors.New("insufficient rights to delete Order")
	}
	return
}

func (customer *Customer) deletePerformedOrder(orderID int) {
	statement := `DELETE FROM performed_orders WHERE order_id = $1`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	stmt.QueryRow(orderID).Scan()
}

//MakeOrderPerformed row in "Order" table
func (customer *Customer) MakeOrderPerformed(performedOrder *PerformedOrder) (err error) {
	if performedOrder.Order.Status.ID == OrderStatusAvailable {
		if performedOrder.Customer.User.ID == customer.User.ID {
			statement := `INSERT INTO performed_orders (order_id, freelancer_id) values ($1, $2)`
			stmt, err := Db.Prepare(statement)
			if err != nil {
				log.Println(err)
			}
			defer stmt.Close()
			err = stmt.QueryRow(performedOrder.Order.ID, performedOrder.Freelancer.User.ID).Scan(&performedOrder.ID)
			if err != nil {
				log.Println(err)
			}
			performedOrder.Order.changeStatus(OrderStatusPerformed)
			//customer.deletePerformedOrder(performedOrder.Order.ID)
		} else {
			return errors.New("insufficient rights make performedOrder performed")
		}
	}else {return errors.New("performedOrder is not available")}
	return
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

func (user *User) updateRating(rait float32)  {
	rait = float32(math.Round(float64(rait*100)) / 100)
	statement := `UPDATE users SET rait = $1 WHERE id = $2 returning rait`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println(err)
	}

	defer stmt.Close()
	err = stmt.QueryRow(rait, user.ID).Scan(&user.Rait)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func (freelancer *Freelancer) calcRating(newMark float32) (countMark sql.NullInt64) {
	err := Db.QueryRow(`SELECT COUNT(customer_comment_id) FROM complete_orders 
							   WHERE freelancer_id = $1`, freelancer.User.ID).Scan(&countMark)
	if err != nil {
		log.Println(err)
		return
	}
	rait := (freelancer.Rait * float32(countMark.Int64 - 1) + newMark) / float32(countMark.Int64)
	freelancer.User.updateRating(rait)
	return
}

func (customer *Customer) calcRating(newMark float32) (countMark sql.NullInt64) {
	err := Db.QueryRow(`SELECT COUNT(freelancer_comment_id) FROM complete_orders 
							   WHERE order_id IN (
							       SELECT id 
							       FROM orders WHERE status_id = 3 and customer_id = $1)`, customer.User.ID).Scan(&countMark)
	if err != nil {
		log.Println(err)
		return
	}
	rait := (customer.Rait * float32(countMark.Int64 - 1) + newMark) / float32(countMark.Int64)
	customer.User.updateRating(rait)
	return
}

//IsAvailable - check Order status, if "available" return true
func (order *Order) IsAvailable() (available bool) {
	if order.Status.ID == OrderStatusAvailable{
		available = true
	}
	return
}

//changeStatus - change status Order in table "orders"
func (order *Order) changeStatus(status int) {
	statement := `UPDATE orders SET status_id = $1 WHERE id = $2 returning id;`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	stmt.QueryRow(status, order.ID).Scan()
}

//OrderDeleteAll - delete all rows in table "orders"
func OrderDeleteAll() (err error) {
	statement := "delete from orders"
	_, err = Db.Exec(statement)
	return
}

func (customer *Customer) Orders() (orders []Order) {
	rows, err := Db.Query(`SELECT id, title, content, customer_id, specialization_id, status_id, created_at FROM orders 
									WHERE customer_id = $1 ORDER BY created_at ASC `, customer.User.ID)
	if err != nil {
		return
	}
	for rows.Next() {
		order := Order{}
		err = rows.Scan(&order.ID, &order.Title, &order.Content, &order.Customer.User.ID, &order.SpecializationID,
						&order.Status.ID, &order.CreatedAt)
		if err != nil {
			log.Println(err)
			return
		}
		order.Status = GetStatusByID(order.Status.ID)
		order.FreelancerRequest = order.GetRequests()
		orders = append(orders, order)
	}
	rows.Close()
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

func (freelancer *Freelancer) PerformingOrders() (performedOrders []PerformedOrder) {
	rows, err := Db.Query(`SELECT id, order_id, freelancer_id
								  FROM performed_orders
								  WHERE freelancer_id = $1;`, freelancer.User.ID)
	if err != nil {
		return
	}
	for rows.Next() {
		performedOrder := PerformedOrder{}
		err = rows.Scan(&performedOrder.ID, &performedOrder.Order.ID, &performedOrder.Freelancer.User.ID)
		if err != nil {
			log.Println(err)
			return
		}
		performedOrder.Freelancer, _ = GetFreelancerByUserID(performedOrder.Freelancer.User.ID)
		performedOrder.Order = GetOrderByID(performedOrder.Order.ID)
		performedOrders = append(performedOrders, performedOrder)
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
		fmt.Println(completeOrder.Freelancer)
		completeOrder.CustomerComment = GetCommentByID(completeOrder.CustomerComment.ID)
		completeOrder.Order = GetOrderByID(completeOrder.Order.ID)
		completeOrders = append(completeOrders, completeOrder)
	}
	rows.Close()
	return
}

func (customer *Customer) PerformedOrders() (performedOrders []PerformedOrder) {
	rows, err := Db.Query(`SELECT id, order_id, freelancer_id
								  FROM performed_orders
								  WHERE
								  	order_id IN (SELECT id FROM orders WHERE customer_id = $1);`, customer.User.ID)
	if err != nil {
		return
	}
	for rows.Next() {
		performedOrder := PerformedOrder{}
		err = rows.Scan(&performedOrder.ID, &performedOrder.Order.ID, &performedOrder.Freelancer.User.ID)
		if err != nil {
			log.Println(err)
			return
		}
		performedOrder.Freelancer, _ = GetFreelancerByUserID(performedOrder.Freelancer.User.ID)
		performedOrder.Order = GetOrderByID(performedOrder.Order.ID)
		performedOrders = append(performedOrders, performedOrder)
	}
	rows.Close()
	return
}

func GetPerformedOrdersByID(orderID int) (performedOrder PerformedOrder) {
	err := Db.QueryRow(`SELECT id, order_id, freelancer_id
								  FROM performed_orders
								  WHERE order_id = $1;`, orderID).Scan(&performedOrder.ID, &performedOrder.Order.ID,
								  	&performedOrder.Freelancer.User.ID)
	if err != nil {
		log.Println(err)
		return
	}
	performedOrder.Freelancer, _ = GetFreelancerByUserID(performedOrder.Freelancer.User.ID)
	performedOrder.Order = GetOrderByID(performedOrder.Order.ID)
	return
}

func ExistOffer (freelancerID, orderID int) (exist bool){
	err := Db.QueryRow(`SELECT EXISTS(SELECT id FROM requests WHERE freelancer_id = $1 and order_id = $2)`,
								freelancerID, orderID).Scan(&exist)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func GetAllOrders() (orders []Order, err error) {
	rows, err := Db.Query(`SELECT id, title, content, customer_id, status_id, specialization_id, created_at FROM orders 
									ORDER BY created_at ASC `)
	if err != nil {
		return
	}
	for rows.Next() {
		order := Order{}
		err = rows.Scan(&order.ID, &order.Title, &order.Content, &order.Customer.User.ID, &order.Status.ID,
						&order.SpecializationID, &order.CreatedAt)
		if err != nil {
			log.Println(err)
			return
		}
		order.Status = GetStatusByID(order.Status.ID)
		order.FreelancerRequest = order.GetRequests()
		orders = append(orders, order)
	}
	rows.Close()
	return
}

func GetOrdersWhere(query string, args ...interface{}) (orders []Order, err error) {
	rows, err := Db.Query(`SELECT id, title, content, customer_id, status_id, specialization_id, created_at FROM orders WHERE ` +
							query + ` ORDER BY created_at ASC`, args...)
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		order := Order{}
		err = rows.Scan(&order.ID, &order.Title, &order.Content, &order.Customer.User.ID, &order.Status.ID,
						&order.SpecializationID, &order.CreatedAt)
		if err != nil {
			log.Println(err)
			return
		}
		order.Status = GetStatusByID(order.Status.ID)
		order.FreelancerRequest = order.GetRequests()
		orders = append(orders, order)
	}
	rows.Close()
	return
}

func GetStatusByID(id int) (status Status) {
	err := Db.QueryRow(`SELECT id, name FROM order_status WHERE id = $1`, id).Scan(&status.ID, &status.Name)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func GetOrderByID(id int) (order Order) {
	err := Db.QueryRow(`SELECT id, title, content, customer_id, specialization_id, status_id, created_at FROM orders
						WHERE id = $1`, id).Scan(&order.ID, &order.Title, &order.Content, &order.Customer.User.ID,
												 &order.SpecializationID, &order.Status.ID, &order.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return
	}
	order.Status = GetStatusByID(order.Status.ID)
	order.Customer, _ = GetCustomerByUserID(order.Customer.User.ID)
	return
}

func (customer *Customer) GetOrdersByStatus(status int) (orders []Order) {
	rows, err := Db.Query(`SELECT id, title, content, customer_id, specialization_id, status_id, created_at FROM orders 
									WHERE customer_id = $1 AND status_id = $2`, customer.User.ID, status)
	if err != nil {
		return
	}
	for rows.Next() {
		order := Order{}
		err = rows.Scan(&order.ID, &order.Title, &order.Content, &order.Customer.User.ID, &order.SpecializationID,
						&order.Status.ID, &order.CreatedAt)
		if err != nil {
			log.Println(err)
			return
		}
		order.Status = GetStatusByID(order.Status.ID)
		order.FreelancerRequest = order.GetRequests()
		orders = append(orders, order)
	}
	rows.Close()
	return
}