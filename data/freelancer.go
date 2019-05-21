package data

import (
	"database/sql"
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
	log.Println(freelancer)
	if err = freelancer.User.UpdateInformation(); err != nil {
		log.Println(err)
	}
	statement := `UPDATE freelancers SET specialization = $1 WHERE id = $2`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println(err)
	}

	defer stmt.Close()
	err = stmt.QueryRow(pq.Array(freelancer.Specialization), freelancer.ID).Scan(&freelancer.ID)
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
	freelancer.User, err = GetUserByID(id)
	if err != nil {
		return
	}
	var tmp []sql.NullInt64
	err = Db.QueryRow(`SELECT id, user_id, specialization FROM freelancers
								WHERE user_id = $1`, id).Scan(&freelancer.ID, &freelancer.User.ID,
		pq.Array(&tmp))
	for _, j := range tmp {
		x := j.Int64
		freelancer.Specialization = append(freelancer.Specialization, int(x))
	}
	return
}

//CheckFreelancer - check exist user in table "freelancers"
func CheckFreelancer(userID int) (exist bool) {
	err := Db.QueryRow(`SELECT EXISTS(SELECT id FROM freelancers WHERE user_id = $1)`, userID).Scan(&exist)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func GetAllFreelancers() (freelancers []Freelancer, err error) {
	var tmp []sql.NullInt64
	rows, err := Db.Query(`SELECT id, user_id, specialization FROM freelancers`)
	if err != nil {
		return
	}
	for rows.Next() {
		freelancer := Freelancer{}
		if err != nil {
			return
		}
		err = rows.Scan(&freelancer.ID, &freelancer.User.ID, pq.Array(&tmp))
		if err != nil {
			log.Println(err)
			return
		}
		freelancer.User, err = GetUserByID(freelancer.User.ID)
		for _, j := range tmp {
			x := j.Int64
			freelancer.Specialization = append(freelancer.Specialization, int(x))
		}
		freelancers = append(freelancers, freelancer)
	}
	rows.Close()
	return
}

func GetFreelancersWhere(query string, args ...interface{}) (freelancers []Freelancer, err error) {
	var tmp []sql.NullInt64
	rows, err := Db.Query(`SELECT id, user_id, specialization, first_name, last_name, phone,
		 				   facebook, skype, about, rait, role_id, photo_url
						   FROM freelancers 
						   JOIN users USING (id)
						   WHERE ` + query, args...)
	if err != nil {
		return
	}
	for rows.Next() {
		freelancer := Freelancer{}
		if err != nil {
			return
		}
		err = rows.Scan(&freelancer.ID, &freelancer.User.ID, pq.Array(&tmp), &freelancer.FirstName, &freelancer.LastName,
						&freelancer.Phone, &freelancer.Facebook, &freelancer.Skype, &freelancer.About, &freelancer.Rait,
						&freelancer.RoleID, &freelancer.Photo)
		if err != nil {
			log.Println(err)
			return
		}
		freelancer.User, err = GetUserByID(freelancer.User.ID)
		for _, j := range tmp {
			x := j.Int64
			freelancer.Specialization = append(freelancer.Specialization, int(x))
		}
		freelancers = append(freelancers, freelancer)
	}
	rows.Close()
	return
}

//GetAllSpecialization return all rows from table "specialization"
func GetAllSpecialization() (specs []Specialization, err error) {
	rows, err := Db.Query(`SELECT id, name FROM specialization`)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		spec := Specialization{}
		err = rows.Scan(&spec.ID, &spec.Name)
		if err != nil {
			log.Println(err)
		} else {
			specs = append(specs, spec)
		}
	}
	return
}

//GetSpecializationName -
func GetSpecializationName(id int) (name string) {
	Db.QueryRow(`SELECT name FROM specialization WHERE id = $1`, id).Scan(&name)
	return
}

func (freelancer *Freelancer) ContainsSpecialization(id int) (exist bool) {
	for i := range freelancer.Specialization {
		if id == freelancer.Specialization[i] {
			return true
		}
	}
	return
}
