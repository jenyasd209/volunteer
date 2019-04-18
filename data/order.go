package data

import (
	"log"
	"time"
)

type Order struct {
	ID      int
	Title   string
	Content string
	Customer
	Status
	CreatedAt time.Time
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
func (customer *Customer) UpdateOrder(order Order) (err error) {
	if order.Customer.ID == customer.ID {
		statement := `UPDATE orders SET title = $1, content = $2 WHERE id = $3`
		stmt, err := Db.Prepare(statement)
		if err != nil {
			panic(err)
		}

		defer stmt.Close()
		return stmt.QueryRow(order.Title, order.Content).Scan(&order.ID)
	} else {
		log.Println("Insufficient rights to delete order")
	}
	return
}

// DeleteOrder row from "order" table
func (customer *Customer) DeleteOrder(order Order) (err error) {
	if order.Customer.ID == customer.ID {
		statement := "DELETE FROM orders WHERE id = $1"
		stmt, err := Db.Prepare(statement)
		if err != nil {
			log.Println(err)
		}

		defer stmt.Close()
		_, err = stmt.Exec(order.ID)
	} else {
		log.Println("Insufficient rights to delete order")
	}
	return
}
