package data

import (
	"log"
	"time"
)

//Customer struct for "customers" table
type Customer struct {
	ID           int
	Organization string
	CreatedAt    time.Time
	User
}

//Create new row from "customer" table
func (customer *Customer) Create() (err error) {
	if err = customer.User.Create(); err != nil {
		return
	}
	statement := `insert into customers (user_id, organization)
								values ($1, $2) returning id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(customer.User.ID, customer.Organization).Scan(&customer.ID)

	return
}

//Update row in "customer" table
func (customer *Customer) Update() (err error) {
	log.Println(customer)
	if err = customer.User.UpdateInformation(); err != nil {
		log.Println(err)
	}
	statement := `UPDATE customers SET organization = $1 WHERE id = $2`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println(err)
	}

	defer stmt.Close()
	err = stmt.QueryRow(customer.Organization, customer.ID).Scan(&customer.ID)
	return
}

// Delete row from "customer" table
func (customer *Customer) Delete() (err error) {
	statement := "DELETE FROM customers WHERE id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}

	defer stmt.Close()
	_, err = stmt.Exec(customer.ID)
	return
}

//GetCustomerByUserID - return customer by user ID
func GetCustomerByUserID(id int) (customer Customer, err error) {
	customer.User, err = GetUserByID(id)
	if err != nil {
		log.Println(err)
		return
	}
	err = Db.QueryRow(`SELECT id, user_id, organization FROM customers
								WHERE user_id = $1`, id).Scan(&customer.ID, &customer.User.ID, &customer.Organization)
	return
}

//CheckCustomer - check exist user in table "customers"
func CheckCustomer(userID int) (exist bool) {
	err := Db.QueryRow(`SELECT EXISTS(SELECT id FROM customers WHERE user_id = $1)`, userID).Scan(&exist)
	if err != nil {
		log.Println(err)
		exist = false
		return
	}
	return
}
