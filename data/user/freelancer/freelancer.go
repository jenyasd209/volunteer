package freelancer

import (
	"graduate/data"
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
	CreatedAt time.Time
}

//Session truct for "freelancer_session" table
type Session struct {
	ID           int
	Email        string
	FreelancerID int
	CreatedAt    time.Time
}

//Create new row from "freelancer" table
func (freelancer *Freelancer) Create() (err error) {
	statement := `insert into freelancers (first_name, last_name, email, password, created_at)
								values ($1, $2, $3, $4, $5) returning id, created_at`
	stmt, err := data.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(freelancer.FirstName, freelancer.LastName,
		freelancer.Email, data.Encrypt(freelancer.Password),
		time.Now()).Scan(&freelancer.ID, &freelancer.CreatedAt)

	return
}

//Delete row from "freelancer" table
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

//GetAllUsers return all rows from table "freelancer"
func GetAllUsers() (freelancers []Freelancer, r error) {
	rows, err := data.Db.Query(`SELECT id, first_name, last_name, email, password,
															phone, facebook, skype, about, created_at FROM freelancer`)
	if err != nil {
		return
	}

	for rows.Next() {
		freelancer := Freelancer{}

		if err = rows.Scan(&freelancer.ID, &freelancer.FirstName, &freelancer.LastName,
			&freelancer.Email, &freelancer.Password, &freelancer.Phone, &freelancer.Facebook,
			&freelancer.Skype, &freelancer.About, &freelancer.CreatedAt); err != nil {
			freelancers = append(freelancers, freelancer)
		}
	}

	rows.Close()
	return
}

//GetUserByEmail return rows with required email
func GetUserByEmail(email string) (freelancer Freelancer, err error) {
	freelancer = Freelancer{}
	err = data.Db.QueryRow(`SELECT first_name, last_name, email, password, phone,
		 											facebook, skype, about, created_at FROM freelancers
													WHERE email = $1`, email).Scan(&freelancer.FirstName,
		&freelancer.LastName, &freelancer.Email, &freelancer.Password,
		&freelancer.Phone, &freelancer.Facebook, &freelancer.Skype,
		&freelancer.About, &freelancer.CreatedAt)
	return
}

//GetUserByID return rows with required ID
func GetUserByID(id int) (freelancer Freelancer, err error) {
	freelancer = Freelancer{}
	err = data.Db.QueryRow(`SELECT first_name, last_name, email, password, phone,
		 											facebook, skype, about,created_at FROM freelancers
													WHERE id = $1`, id).Scan(&freelancer.FirstName,
		&freelancer.LastName, &freelancer.Email, &freelancer.Password,
		&freelancer.Phone, &freelancer.Facebook, &freelancer.Skype,
		&freelancer.About, &freelancer.CreatedAt)
	return
}
