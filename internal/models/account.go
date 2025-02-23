package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Id              int
	Username        string
	HashedPassword  []byte
	Address         string
	TknsEnabled     bool
	TknAddress      string
	MessageMaxChars int
	MinDonation     float64
	ShowAmount      bool
	WidgetId        string
	Created         time.Time
}

const accountColumns = `
  id,
  username,
  password,
  address,
  tkn_enabled,
  tkn_address,
  message_max_char,
  min_donation,
  show_amount,
  widget_id,
  created
`

type AccountModel struct {
	DB *sql.DB
}

func (m *AccountModel) Insert(username, password, address, tknAddress string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	tknEnabled := tknAddress != ""

	stmt := `INSERT INTO account (
  username,
  password,
  address,
  tkn_enabled,
  tkn_address,
  message_max_char,
  min_donation,
  show_amount,
  created)
VALUES($1, $2, $3, $4, $5, 300, 0.00000547, true, CURRENT_TIMESTAMP)`
	_, err = m.DB.Exec(stmt, username, string(hashedPassword), address, tknEnabled, tknAddress)
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
	      tkn_enabled      = $2,
		  tkn_address      = $3,
		  message_max_char = $4,
		  min_donation     = $5,
		  show_amount      = $6
	  WHERE id = $7`
	_, err := m.DB.Exec(stmt,
		account.Address,
		account.TknsEnabled,
		account.TknAddress,
		account.MessageMaxChars,
		account.MinDonation,
		account.ShowAmount,
		account.Id)
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

func (m *AccountModel) Get(id int) (*Account, error) {
	stmt := fmt.Sprintf(`SELECT %s FROM account WHERE id = $1`, accountColumns)
	return GetOneByQuery(m, stmt, id)
}

func (m *AccountModel) GetByUsername(username string) (*Account, error) {
	stmt := fmt.Sprintf(`SELECT %s FROM account WHERE username = $1`, accountColumns)
	return GetOneByQuery(m, stmt, username)
}

func (m *AccountModel) GetWidgetIdByAccountId(id int) (string, error) {
	var widgetId string
	stmt := `SELECT widget_id FROM account WHERE id = $1`
	err := m.DB.QueryRow(stmt, id).Scan(&widgetId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNoRecord
		}
		return "", err
	}
	return widgetId, nil
}

func GetOneByQuery(m *AccountModel, stmt string, id any) (*Account, error) {
	row := m.DB.QueryRow(stmt, id)

	s := &Account{}
	err := row.Scan(
		&s.Id,
		&s.Username,
		&s.HashedPassword,
		&s.Address,
		&s.TknsEnabled,
		&s.TknAddress,
		&s.MessageMaxChars,
		&s.MinDonation,
		&s.ShowAmount,
		&s.WidgetId,
		&s.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}
