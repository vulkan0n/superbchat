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
	Created             time.Time
}

type AccountModel struct {
	DB *sql.DB
}

func (m *AccountModel) Insert(username, password, address string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	stmt := `INSERT INTO account (username, password, address, name_max_char, message_max_char, min_donation, show_amount, created)
  VALUES($1, $2, $3, 25, 300, 0.001, true, CURRENT_TIMESTAMP)`
	_, err = m.DB.Exec(stmt, username, string(hashedPassword), address)
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

func (m *AccountModel) Authenticate(username, password string) (int, error) {
	return 0, nil
}

func (m *AccountModel) Exist(id int) (bool, error) {
	return false, nil
}

func (m *AccountModel) Get(id int) (*Account, error) {
	stmt := `SELECT id, username, password, address, name_max_char, message_max_char, min_donation, show_amount, created
	FROM superchat WHERE id = $1`

	return GetOneByQuery(m, stmt, id)
}

func (m *AccountModel) GetByUsername(username string) (*Account, error) {
	stmt := `SELECT id, username, password, address, name_max_char, message_max_char, min_donation, show_amount, created
	FROM superchat WHERE username = $1`

	return GetOneByQuery(m, stmt, username)
}

func GetOneByQuery(m *AccountModel, stmt string, id any) (*Account, error) {
	row := m.DB.QueryRow(stmt, id)

	s := &Account{}
	err := row.Scan(&s.Id, &s.Username, &s.HashedPassword, &s.Address, &s.NameMaxChars,
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
