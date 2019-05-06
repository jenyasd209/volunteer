package data

import (
	"log"
	"time"

	"github.com/lib/pq"
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
	if err = customer.User.UpdateInformation(); err != nil {
		panic(err)
	}
	statement := `UPDATE customer SET organization WHERE id = $1`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err)
	}

	defer stmt.Close()
	err = stmt.QueryRow(customer.Organization).Scan(&customer.ID)
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

//GetCustomerUserID - return customer by user ID
func GetCustomerUserID(id int) (customer Customer, err error) {
	customer.User, err = GetUserByID(id)
	if err != nil {
		return
	}
	err = Db.QueryRow(`SELECT id, user_id, organization FROM customers
								WHERE user_id = $1`, id).Scan(&customer.ID, &customer.User.ID,
		pq.Array(&customer.Organization))
	// err = Db.QueryRow(`SELECT C.user_id, C.organization, U.email, U.password,
	// 							U.phone, U.facebook, U.skype, U.about, U.rait, U.created_at FROM customers C, users U
	// 							WHERE F.user_id = U.id and F.user_id = $1`, id).Scan(&customer.User.ID,
	// 	&customer.Organization, &customer.Email, &customer.Password,
	// 	&customer.Phone, &customer.Facebook, &customer.Skype, &customer.About,
	// 	&customer.Rating, &customer.CreatedAt)
	return
}

//CheckCustomer - check exist user in table "customers"
func CheckCustomer(userID int) (exist bool, err error) {
	err = Db.QueryRow(`SELECT EXISTS(SELECT id FROM customers WHERE user_id = $1)`, userID).Scan(&exist)
	if err != nil {
		log.Println(err)
		exist = false
		return
	}
	return
}
