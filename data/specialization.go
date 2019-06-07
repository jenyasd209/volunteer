package data

import (
	"database/sql"
	"log"
)


//Specialization struct for "specializations" table
type Specialization struct {
	ID   int
	Name string
}

//GetAllSpecialization return all rows from table "specialization"
func GetAllSpecialization() (specs []Specialization, err error) {
	rows, err := Db.Query(`SELECT id, name FROM specialization`)
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
