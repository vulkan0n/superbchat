package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Superchat struct {
	Id          int       `json:"id"`
	TxId        string    `json:"txId"`
	Name        string    `json:"name"`
	Message     string    `json:"message"`
	Amount      float64   `json:"amount"`
	TknCategory string    `json:"tknCategory"`
	TknSymbol   string    `json:"tknSymbol"`
	TknLogo     string    `json:"tknLogo"`
	IsHidden    bool      `json:"isHidden"`
	Recipient   int       `json:"recipient"`
	Created     time.Time `json:"created"`
}

type SuperchatModel struct {
	DB *sql.DB
}

const superchatColumns = `
tx_id,
name,
message,
amount,
tkn_category,
tkn_symbol,
tkn_logo,
hidden,
account_id,
created
`

func (m *SuperchatModel) Insert(txId string, name string, message string, amount float64,
	tknCategory string, tknSymbol string, tknLogo string, isHidden bool, recipient int) (int, error) {
	stmt := fmt.Sprintf(`INSERT INTO superchat (%s)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP) RETURNING id`, superchatColumns)
	res := m.DB.QueryRow(stmt, txId, name, message, amount, tknCategory, tknSymbol, tknLogo, isHidden, recipient)

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
		&s.TknCategory,
		&s.TknSymbol,
		&s.TknLogo,
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
			&s.TknCategory,
			&s.TknSymbol,
			&s.TknLogo,
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
