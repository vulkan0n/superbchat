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
	IsHidden  bool      `json:"isHidden"`
	Recipient int       `json:"recipient"`
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
  hidden,
  account_id,
  created
`

func (m *SuperchatModel) Insert(txId string, name string, message string, amount float64, tknSymbol string,
	isHidden bool, recipient int) (int, error) {
	stmt := fmt.Sprintf(`INSERT INTO superchat (%s)
    VALUES($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP) RETURNING id`, superchatColumns)
	res := m.DB.QueryRow(stmt, txId, name, message, amount, tknSymbol, isHidden, recipient)

	var id int
	err := res.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *SuperchatModel) Delete(id int) error {
	query := `DELETE FROM superchat WHERE id = $1`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Superchat not found")
	}

	return nil
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
		&s.IsHidden,
		&s.Recipient,
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
			&s.IsHidden,
			&s.Recipient,
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
