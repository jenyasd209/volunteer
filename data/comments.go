package data

import (
	"log"
	"time"
)

//Comment struct for "comments" table
type Comment struct {
	ID       int		`json:"id"`
	Rait     float32	`json:"rait"`
	Text     string		`json:"text"`
	UserID   int		`json:"user_id"`
	CreateAt time.Time	`json:"create_at"`
}

//CreateComment - create new comment
func (user *User) CreateComment(comment *Comment) (err error) {
	statement := `INSERT INTO comments (comment_text, rait, user_id, created_at)
								values ($1, $2, $3, $4) returning id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(comment.Text, comment.Rait, user.ID, time.Now()).Scan(&comment.ID)
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

func (completeOrder *CompleteOrder) UpdateFreelancerComment() {
	statement := `UPDATE complete_orders SET freelancer_comment_id = $1 WHERE id = $2`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println(err)
	}

	defer stmt.Close()
	stmt.QueryRow(&completeOrder.FreelancerComment.ID, &completeOrder.ID).Scan()
	completeOrder.FreelancerComment = GetCommentByID(completeOrder.FreelancerComment.ID)

	completeOrder.Order.Customer.calcRating(completeOrder.FreelancerComment.Rait)
	return
}

//CommentsDeleteAll - delete all rows in table "comments"
func CommentsDeleteAll() (err error) {
	statement := "delete from comments"
	_, err = Db.Exec(statement)
	return
}
