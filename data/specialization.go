package data

import (
	"database/sql"
	"log"
)


//Specialization struct for "specializations" table
type Specialization struct {
	ID   int 	`json:"id"`
	Name string `json:"name"`
}

func (specialization *Specialization) Create () (err error) {
	statement := `INSERT INTO specialization (name) values ($1) returning id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()
	err = stmt.QueryRow(specialization.Name).Scan(&specialization.ID)
	if err != nil {
		log.Println(err)
	}
	return
}

func (specialization *Specialization) Update () (err error) {
	statement := `UPDATE specialization SET name = $1 WHERE id = $2 returning name`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()
	err = stmt.QueryRow(specialization.Name, specialization.ID).Scan(&specialization.Name)
	if err != nil {
		log.Println(err)
	}
	return
}

func (specialization *Specialization) Delete () (err error) {
	statement := `DELETE FROM specialization WHERE id = $1`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()
	err = stmt.QueryRow(specialization.ID).Scan()
	if err != nil {
		log.Println(err)
	}
	return
}

//GetAllSpecialization return all rows from table "specialization"
func GetAllSpecialization() (specs []Specialization, err error) {
	rows, err := Db.Query(`SELECT id, name FROM specialization order by name`)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

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

func (freelancer *Freelancer) calcRating(newMark float32) (countMark sql.NullInt64) {
	err := Db.QueryRow(`SELECT COUNT(customer_comment_id) FROM complete_orders 
							   WHERE freelancer_id = $1`, freelancer.User.ID).Scan(&countMark)
	if err != nil {
		log.Println(err)
		return
	}
	rait := (freelancer.Rait * float32(countMark.Int64 - 1) + newMark) / float32(countMark.Int64)
	freelancer.User.updateRating(rait)
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
