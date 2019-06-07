package data

import (
	"errors"
	"fmt"
	"log"
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