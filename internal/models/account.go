package models

import (
	"database/sql"
	"errors"
	"time"
)

type Account struct {
	Id                  int
	Username            string
	Password            string
	Address             string
	NameMaxChars        int
	MessageMaxChars     int
	MinDonation         float64
	IsDefaultShowAmount bool
	Created             time.Time
}

type AccountModel struct {
	DB *sql.DB
}

func (m *AccountModel) Insert(username string, password string, address string, nameMaxChars int,
	messageMaxChars int, minDonation float64, defaultShowAmount bool) error {
	stmt := `INSERT INTO account (username, password, address, name_max_char, message_max_char, min_donation, show_amount, created)
  VALUES($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)`
	_, err := m.DB.Exec(stmt, username, password, address, nameMaxChars, messageMaxChars, minDonation, defaultShowAmount)
	if err != nil {
		return err
	}
	return nil
}

func (m *AccountModel) Get(id int) (*Account, error) {
	stmt := `SELECT id, username, password, address, name_max_char, message_max_char, min_donation, show_amount, created
	FROM superchat WHERE id = $1`
	row := m.DB.QueryRow(stmt, id)

	s := &Account{}
	err := row.Scan(&s.Id, &s.Username, &s.Password, &s.Address, &s.NameMaxChars,
		&s.MessageMaxChars, &s.MinDonation, &s.IsDefaultShowAmount, &s.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}
