package freelancer

import (
	"fmt"
	"graduate/data"
	"graduate/data/user"
	"time"
)

//Freelancer truct for "freelancer" table
type Freelancer struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Password  string
	Phone     string
	Facebook  string
	Skype     string
	About     string
	Rait      float32
	CreatedAt time.Time
}

const DB_TABLE_NAME = "freelancer_session"
const DB_FIELD_NAME = "freelancer_id"

//Create new row from "freelancer" table
func (freelancer *Freelancer) Create() (err error) {
	statement := `insert into freelancers (first_name, last_name, email, password, created_at)
								values ($1, $2, $3, $4, $5) returning id, created_at`
	stmt, err := data.Db.Prepare(statement)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(freelancer.FirstName, freelancer.LastName,
		freelancer.Email, data.Encrypt(freelancer.Password),
		time.Now()).Scan(&freelancer.ID, &freelancer.CreatedAt)

	return
}

// Delete row from "freelancer" table
func (freelancer *Freelancer) Delete() (err error) {
	statement := "DELETE FROM freelancers WHERE id = $1"
	stmt, err := data.Db.Prepare(statement)
	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(freelancer.ID)
	return
}

func (freelancer *Freelancer) CreateSession() (sess user.Sessionable, err error) {
	session := Session{}

	statement := `INSERT INTO freelancer_session (uuid, email, freelancer_id, created_at) values
	                ($1, $2, $3, $4) returning id, uuid, email, freelancer_id, created_at`
	stmt, err := data.Db.Prepare(statement)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer stmt.Close()
	err = stmt.QueryRow(data.CreateUUID(), freelancer.Email, freelancer.ID, time.Now()).Scan(&session.ID,
		&session.UUID, &session.Email, &session.UserID, &session.CreatrdAt)

	if err != nil {
		fmt.Println(err)
		return
	}

	sess = &session

	return
}

//GetAllUsers return all rows from table "freelancer"
func GetAllUsers() (freelancers []Freelancer, r error) {
	rows, err := data.Db.Query(`SELECT id, first_name, last_name, email, password,
															phone, facebook, skype, about, rait, created_at FROM freelancer`)
	if err != nil {
		return
	}

	for rows.Next() {
		freelancer := Freelancer{}

		if err = rows.Scan(&freelancer.ID, &freelancer.FirstName, &freelancer.LastName,
			&freelancer.Email, &freelancer.Password, &freelancer.Phone, &freelancer.Facebook,
			&freelancer.Skype, &freelancer.About, &freelancer.Rait, &freelancer.CreatedAt); err != nil {
			freelancers = append(freelancers, freelancer)
		}
	}

	rows.Close()
	return
}

//GetUserByEmail return rows with required email
func GetUserByEmail(email string) (freelancer Freelancer, err error) {
	// var freelancer *Freelancer
	err = data.Db.QueryRow(`SELECT id, first_name, last_name, email, password, phone,
		 											facebook, skype, about, rait, created_at FROM freelancers
													WHERE email = $1`, email).Scan(&freelancer.ID, &freelancer.FirstName,
		&freelancer.LastName, &freelancer.Email, &freelancer.Password,
		&freelancer.Phone, &freelancer.Facebook, &freelancer.Skype,
		&freelancer.About, &freelancer.Rait, &freelancer.CreatedAt)
	return
}

//GetUserByID return rows with required ID
func GetUserByID(id int) (freelancer Freelancer, err error) {
	// freelancer = Freelancer{}
	err = data.Db.QueryRow(`SELECT id, first_name, last_name, email, password, phone,
		 											facebook, skype, about, rait, created_at FROM freelancers
													WHERE id = $1`, id).Scan(&freelancer.ID, &freelancer.FirstName,
		&freelancer.LastName, &freelancer.Email, &freelancer.Password, &freelancer.Phone,
		&freelancer.Facebook, &freelancer.Skype, &freelancer.About, &freelancer.Rait,
		&freelancer.CreatedAt)
	return
}
