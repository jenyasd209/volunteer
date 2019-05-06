package data

import (
	"errors"
	"log"
	"time"
)

//Constant status for table "Status"
const (
	OrderStatusAvailable = 1
	OrderStatusPerformed = 2
	OrderStatusDone      = 3
)

//Order struct for "orders" table
type Order struct {
	ID      int
	Title   string
	Content string
	Customer
	Status
	CreatedAt time.Time
}

//PerformedOrder struct for "performed_order" table
type PerformedOrder struct {
	ID int
	Order
	FreelancerID int
}

//CompleteOrder struct for "complete_order" table
type CompleteOrder struct {
	ID int
	Order
	FreelancerID      int
	FreelancerComment Comment
	CustomerComment   Comment
	DateComplete      time.Time
}

//Status struct for "order_status" table
type Status struct {
	ID   int
	Name string
}

//CreateOrder new row in "order" table
func (customer *Customer) CreateOrder(title string, content string) (order Order, err error) {
	statement := `insert into orders (title, content, customer_id, created_at)
								values ($1, $2, $3, $4) returning id, title, content, customer_id, created_at`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err)
	}

	defer stmt.Close()
	err = stmt.QueryRow(title, content, customer.ID, time.Now()).Scan(&order.ID, order.Title,
		order.Content, order.Customer, order.CreatedAt)
	return
}

//UpdateOrder row in "order" table
func (customer *Customer) UpdateOrder(order *Order) (err error) {
	if order.Customer.ID == customer.ID {
		if order.IsAvailable() {
			statement := `UPDATE orders SET title = $1, content = $2 WHERE id = $3
			returning id`
			stmt, err := Db.Prepare(statement)
			if err != nil {
				return err
			}

			defer stmt.Close()
			stmt.QueryRow(order.Title, order.Content, order.ID).Scan(&order.ID)
		} else {
			return errors.New("Can't change order, it status not 'Available'")
		}
	} else {
		return errors.New("Insufficient rights to update order")
	}
	return
}

// DeleteOrder row from "order" table
func (customer *Customer) DeleteOrder(order Order) (err error) {
	if order.Customer.ID == customer.ID {
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
			return errors.New("Can't change order, it status not 'Available'")
		}
	} else {
		return errors.New("Insufficient rights to delete order")
	}
	return
}

func (customer *Customer) deletePerformedOrder(orderID int) {
	statement := `DELETE FROM performed_order WHERE order_id = $1`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	stmt.QueryRow(orderID).Scan()
}

//MakeOrderPerformed row in "order" table
func (customer *Customer) MakeOrderPerformed(performedOrder *PerformedOrder) (err error) {
	if performedOrder.Customer.ID == customer.ID {
		statement := `INSERT INTO performed_orders (order_id, freelancer_id) values ($1, $2)`
		stmt, err := Db.Prepare(statement)
		if err != nil {
			log.Println(err)
		}
		defer stmt.Close()
		err = stmt.QueryRow(performedOrder.Order.ID, performedOrder.FreelancerID).Scan(&performedOrder.ID)
		if err != nil {
			log.Println(err)
		}
		performedOrder.Order.changeStatus(OrderStatusPerformed)
		customer.deletePerformedOrder(performedOrder.Order.ID)
	} else {
		return errors.New("Insufficient rights make order performed")
	}
	return
}

//MakeOrderDone row in "order" table
func (customer *Customer) MakeOrderDone(completeOrder *CompleteOrder) (err error) {
	if completeOrder.Customer.ID == customer.ID {
		statement := `INSERT INTO complete_orders (order_id, freelancer_id)
									values ($1, $2)`
		stmt, err := Db.Prepare(statement)
		if err != nil {
			log.Println(err)
			return err
		}

		err = customer.CreateComment(completeOrder.CustomerComment)
		if err != nil {
			log.Println(err)
			return err
		}
		defer stmt.Close()
		stmt.QueryRow(completeOrder.Order.ID, completeOrder.FreelancerID).Scan(&completeOrder.ID)
		completeOrder.Order.changeStatus(OrderStatusDone)
	} else {
		return errors.New("Insufficient rights make order done")
	}
	return
}

//IsAvailable - check order status, if "available" return true
func (order *Order) IsAvailable() (available bool) {
	if order.Status.Name == "Available" {
		available = true
	}
	return
}

//changeStatus - change status order in table "orders"
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
