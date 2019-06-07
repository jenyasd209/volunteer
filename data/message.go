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

//Message struct for "messages" table
type Message struct {
	ID         int 			`json:"id"`
	SenderID   int 			`json:"sender_id"`
	ReceiverID int 			`json:"receiver_id"`
	DialogID   int 			`json:"dialog_id"`
	Text       string 		`json:"text"`
	Read       bool 		`json:"read"`
	DateSend   time.Time 	`json:"date_send"`
}

//SendMessage - create new message
func (user *User) SendMessage(message Message) (err error) {
	statement := `INSERT INTO messages (sender_id, receiver_id, dialog_id, text_message, date_send)
								values ($1, $2, $3, $4, $5)`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(user.ID, message.ReceiverID, getDialog(user.ID, message.ReceiverID), message.Text, time.Now().UTC()).Scan()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("send msg")
	return
}

//ReadMessage - make message read "True"
func (user *User) ReadMessage(dialogID int) (err error) {
	statement := `UPDATE messages SET read = true WHERE dialog_id = $1 and receiver_id = $2`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(dialogID, user.ID).Scan()
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

func (dialog *Dialog) GetMessages() (messages []*Message) {
	rows, err := Db.Query(`SELECT id, sender_id, receiver_id, dialog_id, text_message, read, date_send
                         FROM messages WHERE dialog_id = $1 order by date_send ASC `, dialog.ID)
	if err != nil {
		log.Println(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		message := Message{}

		err = rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.DialogID, &message.Text,
			&message.Read, &message.DateSend)
		if err == nil {
			messages = append(messages, &message)
			fmt.Println(message)
		}
	}
	return
}

//GetAllMessageBySenderID - return array message by senderID
func GetAllMessageBySenderID(senderID int) (messages []Message, err error) {
	rows, err := Db.Query(`SELECT id, sender_id, receiver_id, text_message, read, date_send
                         FROM messages WHERE sender_id = $1`, senderID)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		message := Message{}

		if err = rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Text,
			&message.Read, &message.DateSend); err == nil {
			messages = append(messages, message)
		}
	}
	return
}

//GetAllMessageByReceiverID - return array message by receiverID
func GetAllMessageByReceiverID(receiverID int) (messages []Message, err error) {
	rows, err := Db.Query(`SELECT id, sender_id, receiver_id, text_message, read, date_send
                         FROM messages WHERE receiver_id = $1`, receiverID)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		message := Message{}

		if err = rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Text,
			&message.Read, &message.DateSend); err == nil {
			messages = append(messages, message)
		}
	}
	return
}

//MessagesDeleteAll - delete all rows in table "messages"
func MessagesDeleteAll() (err error) {
	statement := "delete from messages"
	_, err = Db.Exec(statement)
	return
}
