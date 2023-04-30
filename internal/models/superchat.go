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
	IsHidden  bool
	Recipient int
	IsPaid    bool
	IsAlerted bool
	Created   time.Time
}

type SuperchatModel struct {
	DB *sql.DB
}

func (m *SuperchatModel) Insert(txId string, name string, message string, amount float64,
	isHidden bool, recipient int) (int, error) {
	stmt := `INSERT INTO superchat (tx_id, name, message, amount, hidden, account_id, paid, alerted, created)
    VALUES($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP) RETURNING id`
	res := m.DB.QueryRow(stmt, txId, name, message, amount, isHidden, recipient, false, false)

	var id int
	err := res.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *SuperchatModel) Get(id int) (*Superchat, error) {
	stmt := `SELECT id, tx_id, name, message, amount, hidden, account_id, paid, alerted, created
	FROM superchat WHERE id = $1`
	row := m.DB.QueryRow(stmt, id)

	s := &Superchat{}
	err := row.Scan(&s.Id, &s.TxId, &s.Name, &s.Message, &s.Amount, &s.IsHidden,
		&s.Recipient, &s.IsPaid, &s.IsAlerted, &s.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SuperchatModel) GetFromAccount(accountId int) ([]*Superchat, error) {
	stmt := `SELECT id, tx_id, name, message, amount, hidden, account_id, paid, alerted, created
	FROM superchat WHERE account_id = $1 ORDER BY created DESC`
	rows, err := m.DB.Query(stmt, accountId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	superchats := []*Superchat{}

	for rows.Next() {
		s := &Superchat{}
		err = rows.Scan(&s.Id, &s.TxId, &s.Name, &s.Message, &s.Amount, &s.IsHidden,
			&s.Recipient, &s.IsPaid, &s.IsAlerted, &s.Created)
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

func (m *SuperchatModel) SetAsPaid(txId string, superchatId int) error {
	stmt := `UPDATE superchat SET paid = true, tx_id = $1 WHERE id = $2`
	_, err := m.DB.Exec(stmt, txId, superchatId)
	if err != nil {
		return err
	}
	return nil
}

func (m *SuperchatModel) SetAsAlerted(superchatId int) error {
	stmt := `UPDATE superchat SET alerted = true WHERE id = $1`
	_, err := m.DB.Exec(stmt, superchatId)
	if err != nil {
		return err
	}
	return nil
}
