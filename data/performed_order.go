package data

import (
	"errors"
	"log"
)

//PerformedOrder struct for "performed_order" table
type PerformedOrder struct {
	ID int				`json:"id"`
	Order				`json:"order"`
	Freelancer			`json:"freelancer"`
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
