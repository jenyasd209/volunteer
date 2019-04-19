package data

import (
	"errors"
	"log"
	"time"
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
}

//CompleteOrder struct for "complete_order" table
type CompleteOrder struct {
	ID int
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

//MakeOrderPerformed row in "order" table
func (customer *Customer) MakeOrderPerformed(order *Order) (err error) {
	if order.Customer.ID == customer.ID {
		statement := `UPDATE orders SET status_id = 2 WHERE id = $1
									returning id`
		stmt, err := Db.Prepare(statement)
		if err != nil {
			log.Println(err)
		}

		defer stmt.Close()
		stmt.QueryRow(order.ID).Scan(&order.ID)
	} else {
		return errors.New("Insufficient rights make order performed")
	}
	return
}

//MakeOrderDone row in "order" table
func (customer *Customer) MakeOrderDone(order *Order) (err error) {
	if order.Customer.ID == customer.ID {
		statement := `UPDATE orders SET status_id = 2 WHERE id = $1
									returning id`
		stmt, err := Db.Prepare(statement)
		if err != nil {
			log.Println(err)
		}

		defer stmt.Close()
		stmt.QueryRow(order.ID).Scan(&order.ID)
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
