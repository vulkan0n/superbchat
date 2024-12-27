package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Superchat struct {
	Id        int       `json:"id"`
	TxId      string    `json:"txId"`
	Name      string    `json:"name"`
	Message   string    `json:"message"`
	Amount    float64   `json:"amount"`
	TknSymbol string    `json:"tknSymbol"`
	TknAmount float64   `json:"tknAmount"`
	IsHidden  bool      `json:"isHidden"`
	Recipient int       `json:"recipient"`
	IsPaid    bool      `json:"isPaid"`
	IsAlerted bool      `json:"isAlerted"`
	Created   time.Time `json:"created"`
}

type SuperchatModel struct {
	DB *sql.DB
}

const superchatColumns = `
  tx_id,
  name,
  message,
  amount,
  tkn_symbol,
  tkn_amount,
  hidden,
  account_id,
  paid,
  alerted,
  created
`

func (m *SuperchatModel) Insert(txId string, name string, message string, amount float64, tknSymbol string,
	tknAmount float64, isHidden bool, recipient int) (int, error) {
	stmt := fmt.Sprintf(`INSERT INTO superchat (%s)
    VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP) RETURNING id`, superchatColumns)
	res := m.DB.QueryRow(stmt, txId, name, message, amount, tknSymbol, tknAmount, isHidden, recipient, false, false)

	var id int
	err := res.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func scanSuperchat(row *sql.Row) (*Superchat, error) {
	s := &Superchat{}
	err := row.Scan(
		&s.Id,
		&s.TxId,
		&s.Name,
		&s.Message,
		&s.Amount,
		&s.TknSymbol,
		&s.TknAmount,
		&s.IsHidden,
		&s.Recipient,
		&s.IsPaid,
		&s.IsAlerted,
		&s.Created,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return s, nil
}

func (m *SuperchatModel) Get(id int) (*Superchat, error) {
	stmt := fmt.Sprintf(`SELECT id, %s FROM superchat WHERE id = $1`, superchatColumns)
	row := m.DB.QueryRow(stmt, id)
	return scanSuperchat(row)
}

func (m *SuperchatModel) GetFromAccount(accountId int) ([]*Superchat, error) {
	stmt := fmt.Sprintf(`SELECT id, %s FROM superchat WHERE account_id = $1 ORDER BY created DESC`, superchatColumns)
	rows, err := m.DB.Query(stmt, accountId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	superchats := []*Superchat{}

	for rows.Next() {
		s := &Superchat{}
		err = rows.Scan(
			&s.Id,
			&s.TxId,
			&s.Name,
			&s.Message,
			&s.Amount,
			&s.TknSymbol,
			&s.TknAmount,
			&s.IsHidden,
			&s.Recipient,
			&s.IsPaid,
			&s.IsAlerted,
			&s.Created,
		)
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

func (m *SuperchatModel) GetOldestNotAlerted(token string) (*Superchat, error) {
	stmt := `SELECT superchat.id,
					superchat.name,
					superchat.message,
					superchat.amount,
					superchat.hidden
			 FROM superchat
			 INNER JOIN account ON superchat.account_id = account.id
			 WHERE account.token = $1 AND
				   superchat.paid = true AND
				   superchat.alerted = false
			 ORDER BY superchat.created`
	row := m.DB.QueryRow(stmt, token)

	s := &Superchat{}
	err := row.Scan(&s.Id, &s.Name, &s.Message, &s.Amount, &s.IsHidden)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
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
