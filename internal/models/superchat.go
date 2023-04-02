package models

import (
	"database/sql"
	"errors"
	"time"
)

type Superchat struct {
	Id        int
	TxId      string
	Name      string
	Message   string
	Amount    float64
	Hidden    bool
	Recipient int
	Created   time.Time
}

type SuperchatModel struct {
	DB *sql.DB
}

func (m *SuperchatModel) Insert(txId string, name string, message string, amount float64, hidden bool, recipient int) (int, error) {
	stmt := `INSERT INTO superchat (tx_id, name, message, amount, hidden, user_id, created)
    VALUES(?, ?, ?, ?, ?, ?, UTC_TIMESTAMP())`
	result, err := m.DB.Exec(stmt, txId, name, message, amount, hidden, recipient)
	if err != nil {
		return 0, nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(id), nil
}

func (m *SuperchatModel) Get(id int) (*Superchat, error) {
	stmt := `SELECT id, tx_id, name, message, amount, hidden, user_id, created FROM superchat WHERE id = ?`
	row := m.DB.QueryRow(stmt, id)

	s := &Superchat{}
	err := row.Scan(&s.Id, &s.TxId, &s.Name, &s.Message, &s.Amount, &s.Hidden, &s.Recipient, &s.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SuperchatModel) GetFromUser(userId int) ([]*Superchat, error) {
	stmt := `SELECT id, tx_id, name, message, amount, hidden, user_id, created FROM superchat WHERE user_id = ?`
	rows, err := m.DB.Query(stmt, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	superchats := []*Superchat{}

	for rows.Next() {
		s := &Superchat{}
		err = rows.Scan(&s.Id, &s.TxId, &s.Name, &s.Message, &s.Amount, &s.Hidden, &s.Recipient, &s.Created)
		if err != nil {
			return nil, err
		}
		superchats = append(superchats, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return superchats, nil
}
