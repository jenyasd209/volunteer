package data

import (
	"log"

	"github.com/lib/pq"
)

//Freelancer struct for "freelancer" table
type Freelancer struct {
	ID             int
	Specialization []int
	User
}

//Specialization struct for "specializations" table
type Specialization struct {
	ID   int
	Name string
}

//Create new row from "freelancer" table
func (freelancer *Freelancer) Create() (err error) {
	if err = freelancer.User.Create(); err != nil {
		return
	}
	statement := `insert into freelancers (user_id, specialization)
								values ($1, $2) returning id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}

	defer stmt.Close()
	err = stmt.QueryRow(freelancer.User.ID, pq.Array(freelancer.Specialization)).Scan(&freelancer.ID)
	return
}

//Update row in "freelancer" table
func (freelancer *Freelancer) Update() (err error) {
	if err = freelancer.User.Update(); err != nil {
		panic(err)
	}
	statement := `UPDATE freelancers SET spetialization WHERE id = $1`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err)
	}

	defer stmt.Close()
	err = stmt.QueryRow(freelancer.Specialization).Scan(&freelancer.ID)
	return
}

// Delete row from "freelancer" table
func (freelancer *Freelancer) Delete() (err error) {
	statement := "DELETE FROM freelancers WHERE id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}

	defer stmt.Close()
	_, err = stmt.Exec(freelancer.ID)
	return
}

//GetFreelancerByUserID - return user with set ID
func GetFreelancerByUserID(id int) (freelancer Freelancer, err error) {
	err = Db.QueryRow(`SELECT F.user_id, F.specialization, U.email, U.password,
								U.phone, U.facebook, U.skype, U.about, U.rait, U.created_at FROM freelancers F, users U
								WHERE F.user_id = U.id and F.user_id = $1`, id).Scan(&freelancer.User.ID,
		&freelancer.Specialization, &freelancer.Email, &freelancer.Password,
		&freelancer.Phone, &freelancer.Facebook, &freelancer.Skype, &freelancer.About,
		&freelancer.Rait, &freelancer.CreatedAt)
	return
}

//CheckFreelancer - check exist user in table "freelancers"
func CheckFreelancer(userID int) (exist bool, err error) {
	err = Db.QueryRow(`SELECT EXISTS(SELECT id FROM freelancers WHERE user_id = $1)`, userID).Scan(&exist)
	if err != nil {
		log.Println(err)
		exist = false
		return
	}
	return
}
