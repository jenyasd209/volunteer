package data

import (
	"fmt"
	"log"
	"time"
)

type Dialog struct {
	ID          int 		`json:"id"`
	UserCurrent User 		`json:"user_current"`
	UserTwo     User 		`json:"user_two"`
	CreatedAt   time.Time	`json:"created_At"`
	Messages    []*Message 	`json:"messages"`
}

func createDialog(userOneID, userTwoID int) (dialogID int) {
	statement := `INSERT INTO dialogs (user1_id, user2_id, date_created)
								values ($1, $2, $3) returning id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(userOneID, userTwoID, time.Now()).Scan(&dialogID)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func getDialog(userOneID, userTwoID int) (dialogID int) {
	//err := Db.QueryRow(`SELECT EXISTS(SELECT id FROM dialogs WHERE (user1_id = $1 and user2_id = $2)
	//							or (user1_id = $2 and user2_id = $1))`, userOneID, userTwoID).Scan(&exist)
	err := Db.QueryRow(`SELECT id FROM dialogs WHERE (user1_id = $1 and user2_id = $2)
    							or (user1_id = $2 and user2_id = $1)`, userOneID, userTwoID).Scan(&dialogID)
	if err != nil {
		log.Println(err)
		fmt.Println("Dialog no found")
		dialogID = createDialog(userOneID, userTwoID)
		return
	}
	fmt.Println("Dialog ")
	return
}

func (user *User) Dialogs() (dialogs []*Dialog) {
	var userOne, userTwo int
	rows, err := Db.Query(`SELECT id, user1_id, user2_id, date_created
                         FROM dialogs WHERE user2_id = $1 or user1_id = $1 ORDER BY  date_created DESC `, user.ID)
	if err != nil {
		log.Println(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		dialog := Dialog{}
		err = rows.Scan(&dialog.ID, &userOne, &userTwo, &dialog.CreatedAt)
		if user.ID == userOne{
			dialog.UserCurrent, _ = GetUserByID(userOne)
			dialog.UserTwo, _ = GetUserByID(userTwo)
		}else {
			dialog.UserCurrent, _ = GetUserByID(userTwo)
			dialog.UserTwo, _ = GetUserByID(userOne)
		}

		if err == nil {
			dialog.Messages = dialog.GetMessages()
			dialogs = append(dialogs, &dialog)
		}else if err != nil {
			log.Println(err)
		}
	}
	return
}

func (user *User) DialogByID(dialogID int) (dialog Dialog) {
	var userOne, userTwo int
	err := Db.QueryRow(`SELECT id, user1_id, user2_id, date_created
                         FROM dialogs WHERE (user2_id = $1 or user1_id = $1) and id = $2`, user.ID,
		dialogID).Scan(&dialog.ID, &userOne, &userTwo, &dialog.CreatedAt)
	if err != nil {
		log.Println(err)
		fmt.Println("return")
		return
	}
	if user.ID == userOne{
		dialog.UserCurrent, _ = GetUserByID(userOne)
		dialog.UserTwo, _ = GetUserByID(userTwo)
	}else {
		dialog.UserCurrent, _ = GetUserByID(userTwo)
		dialog.UserTwo, _ = GetUserByID(userOne)
	}

	if err == nil {
		dialog.Messages = dialog.GetMessages()
	}else if err != nil {
		log.Println(err)
	}
	return
}
