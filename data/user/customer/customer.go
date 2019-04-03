package customer

import "graduate/data"

//Company struct fro "companies" table
type Company struct {
	ID    int
	Title string
}

//Create new row from "company" table
func (company *Company) Create() (err error) {
	statement := "INSERT INTO companies (title) values ($1) returning id"
	stmt, err := data.Db.Prepare(statement)
	if err != nil {
		return
	}

	stmt.Close()

	err = stmt.QueryRow(company.Title).Scan(&company.ID)
	return
}

//Update row in "companies" table
func (company *Company) Update() (err error) {
	statement := `UPDATE company
                SET title = $2
                WHERE id = $1`
	stmt, err := data.Db.Prepare(statement)
	if err != nil {
		return
	}

	defer stmt.Close()
	_, err = stmt.Exec(company.ID, company.Title)

	return
}

//Delete row from "companies" table
func (company *Company) Delete() {
	statement := "DELETE FROM company WHERE id = $1"
	stmt, err := data.Db.Prepare(statement)
	if err != nil {
		return
	}

	defer stmt.Close()
	_, err = stmt.Exec(company.ID)

	return
}
