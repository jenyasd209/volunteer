package user

import "graduate/data"

//Freelancer truct for "freelancer" table
type Freelancer struct {
	ID        int
	FirstName string
	LastName  string
	User
}

//Create new row from "freelancer" table
func (freelancer *Freelancer) Create() (err error) {
	if err = freelancer.User.Create(); err != nil {
		panic(err)
	}
	statement := `insert into freelancers (user_id, first_name, last_name)
								values ($1, $2, $3) returning id`
	stmt, err := data.Db.Prepare(statement)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(freelancer.User.ID, freelancer.FirstName, freelancer.LastName).Scan(&freelancer.ID)

	return
}

func (freelancer *Freelancer) Update() (err error) {
	return
}

func GetByUserID(id int) (freelancer Freelancer, err error) {
	err = data.Db.QueryRow(`SELECT F.user_id, F.first_name, F.last_name, U.email, U.password,
								U.phone, U.facebook, U.skype, U.about, U.rait, U.created_at FROM freelancers F, users U
								WHERE F.user_id = U.id and F.user_id = $1`, id).Scan(&freelancer.User.ID,
		&freelancer.FirstName, &freelancer.LastName, &freelancer.Email, &freelancer.Password,
		&freelancer.Phone, &freelancer.Facebook, &freelancer.Skype, &freelancer.About,
		&freelancer.Rait, &freelancer.CreatedAt)
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
