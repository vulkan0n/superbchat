package models

import (
	"database/sql"
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

func (m *SuperchatModel) Get(id string) (*Superchat, error) {
	return nil, nil
}
