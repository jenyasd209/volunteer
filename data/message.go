package data

import (
	"log"
	"time"
)

//Message struct for "messages" table
type Message struct {
	ID         int
	SenderID   int
	ReceiverID int
	Text       string
	Read       bool
	DateSend   time.Time
}

//SendMessage - create new message
func (user *User) SendMessage(message Message) (err error) {
	statement := `INSERT INTO messages (sender_id, receiver_id, text_message, date_send)
								values ($1, $2, $3, $4)`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(message.SenderID, message.ReceiverID, message.Text, message.DateSend).Scan()
	if err != nil {
		log.Println(err)
		return
	}
	return
}

//ReadMessage - make message read "True"
func (user *User) ReadMessage(messageID int) (err error) {
	statement := `UPDATE messages SET read = true WHERE id = $1 receiver_id = $2`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(messageID, user.ID).Scan()
	if err != nil {
		log.Println(err)
		return
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
			&message.Read, &message.DateSend); err != nil {
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
			&message.Read, &message.DateSend); err != nil {
			messages = append(messages, message)
		}
	}
	return
}
