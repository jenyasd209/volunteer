package data

import (
	"fmt"
	"log"
	"time"
)

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
