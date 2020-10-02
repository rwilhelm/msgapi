package db

import (
	"database/sql"
	"git.sr.ht/~rxw/msgapi/models"
)

func (db Database) GetAllMsgs() (*models.MsgList, error) {
	list := &models.MsgList{}

	rows, err := db.Conn.Query("SELECT * FROM msgs ORDER BY ID DESC")
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var msg models.Msg
		err := rows.Scan(&msg.ID, &msg.Name, &msg.Email, &msg.Body, &msg.CreatedAt)
		if err != nil {
			return list, err
		}
		list.Msgs = append(list.Msgs, msg)
	}
	return list, nil
}

func (db Database) AddMsg(msg *models.Msg) error {
	var id int
	var createdAt string
	query := `INSERT INTO msgs (name, email, body) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := db.Conn.QueryRow(query, msg.Name, msg.Email, msg.Body).Scan(&id, &createdAt)
	if err != nil {
		return err
	}

	msg.ID = id
	msg.CreatedAt = createdAt
	return nil
}

func (db Database) GetMsgById(msgId int) (models.Msg, error) {
	msg := models.Msg{}

	query := `SELECT * FROM msgs WHERE id = $1;`
	row := db.Conn.QueryRow(query, msgId)
	switch err := row.Scan(&msg.ID, &msg.Name, &msg.Email, &msg.Body, &msg.CreatedAt); err {
	case sql.ErrNoRows:
		return msg, ErrNoMatch
	default:
		return msg, err
	}
}

func (db Database) DeleteMsg(msgId int) error {
	query := `DELETE FROM msgs WHERE id = $1;`
	_, err := db.Conn.Exec(query, msgId)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}

func (db Database) UpdateMsg(msgId int, msgData models.Msg) (models.Msg, error) {
	msg := models.Msg{}
	query := `UPDATE msgs SET name=$1, email=$2, body=$3 WHERE id=$4 RETURNING id, name, email, body, created_at;`
	err := db.Conn.QueryRow(query, msgData.Name, msgData.Email, msgData.Body, msgId).Scan(&msg.ID, &msg.Name, &msg.Email, &msg.Body, &msg.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return msg, ErrNoMatch
		}
		return msg, err
	}

	return msg, nil
}
