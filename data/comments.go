package data

import "time"

//Comment struct for "comments" table
type Comment struct {
	ID       int
	Rating   float32
	Text     string
	UserID   int
	CreateAt time.Time
}

//CreateComment - create new comment
func (user *User) CreateComment(comment Comment) (err error) {
	statement := `INSERT INTO comments (comment_text, rait, user_id, created_at)
								values ($1, $2, $3, $4) returning id, comment_text, rait, user_id, created_at`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(comment.Text, comment.Rating, user.ID, time.Now()).Scan()
	return
}

//CommentsDeleteAll - delete all rows in table "comments"
func CommentsDeleteAll() (err error) {
	statement := "delete from comments"
	_, err = Db.Exec(statement)
	return
}
