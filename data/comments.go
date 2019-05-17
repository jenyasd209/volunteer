package data

import (
	"log"
	"time"
)

//Comment struct for "comments" table
type Comment struct {
	ID       int
	Rait     float32
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
	err = stmt.QueryRow(comment.Text, comment.Rait, user.ID, time.Now()).Scan()
	return
}

func GetCommentByID(commentID int) (comment Comment) {
	err := Db.QueryRow(`SELECT id, rait,comment_text, user_id, created_at  FROM comments WHERE id = $1`,
						commentID).Scan(&comment.ID, &comment.Rait, &comment.Text, &comment.UserID, &comment.CreateAt)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

//CommentsDeleteAll - delete all rows in table "comments"
func CommentsDeleteAll() (err error) {
	statement := "delete from comments"
	_, err = Db.Exec(statement)
	return
}
