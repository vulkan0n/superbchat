package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Id                  int
	Username            string
	HashedPassword      []byte
	Address             string
	NameMaxChars        int
	MessageMaxChars     int
	MinDonation         float64
	IsDefaultShowAmount bool
	Token               string
	Created             time.Time
}

type AccountModel struct {
	DB *sql.DB
}

func (m *AccountModel) Insert(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	stmt := `INSERT INTO account (username, password, address, name_max_char, message_max_char, min_donation, show_amount, created)
  VALUES($1, $2, '', 25, 300, 0.001, true, CURRENT_TIMESTAMP)`
	_, err = m.DB.Exec(stmt, username, string(hashedPassword))
	if err != nil {
		var posgreSQLError *pq.Error
		if errors.As(err, &posgreSQLError) {
			if string(posgreSQLError.Code) == "23505" &&
				strings.Contains(posgreSQLError.Message, "account_username_key") {
				return ErrDuplicateUser
			}
		}
		return err
	}
	return nil
}

func (m *AccountModel) Update(account *Account) error {
	stmt := `UPDATE account
      SET address          = $1,
	      name_max_char    = $2,
		  message_max_char = $3,
		  min_donation     = $4,
		  show_amount      = $5
	  WHERE id = $6`
	_, err := m.DB.Exec(stmt, account.Address, account.NameMaxChars, account.MessageMaxChars,
		account.MinDonation, account.IsDefaultShowAmount, account.Id)
	if err != nil {
		return err
	}
	return nil
}
func (m *AccountModel) Authenticate(username, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := `SELECT id, password FROM account WHERE username = $1`
	err := m.DB.QueryRow(stmt, username).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *AccountModel) Exist(token string) bool {
	stmt := `SELECT id, username, password, address, name_max_char, message_max_char, min_donation, show_amount, token, created
	FROM account WHERE token = $1`
	_, err := GetOneByQuery(m, stmt, token)
	return (err == nil)
}

func (m *AccountModel) Get(id int) (*Account, error) {
	stmt := `SELECT id, username, password, address, name_max_char, message_max_char, min_donation, show_amount, token, created
	FROM account WHERE id = $1`

	return GetOneByQuery(m, stmt, id)
}

func (m *AccountModel) GetByUsername(username string) (*Account, error) {
	stmt := `SELECT id, username, password, address, name_max_char, message_max_char, min_donation, show_amount, token, created
	FROM account WHERE username = $1`

	return GetOneByQuery(m, stmt, username)
}

func GetOneByQuery(m *AccountModel, stmt string, id any) (*Account, error) {
	row := m.DB.QueryRow(stmt, id)

	s := &Account{}
	err := row.Scan(&s.Id, &s.Username, &s.HashedPassword, &s.Address, &s.NameMaxChars,
		&s.MessageMaxChars, &s.MinDonation, &s.IsDefaultShowAmount, &s.Token, &s.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}
