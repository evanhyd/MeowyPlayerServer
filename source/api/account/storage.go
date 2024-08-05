package account

import (
	"encoding/json"
	"errors"
	"os"
)

// Replace with DB if hardware allows
type accountStorage struct {
	accountDir string
	accounts   accountMap
}

func makeStorage() accountStorage {
	const kAccountDir = "account.json"
	return accountStorage{accountDir: kAccountDir, accounts: accountMap{}}
}

func (s *accountStorage) initialize() error {
	data, err := os.ReadFile(s.accountDir)
	if errors.Is(err, os.ErrNotExist) {
		return s.save()
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.accounts)
}

func (s *accountStorage) save() error {
	data, err := json.Marshal(&s.accounts)
	if err != nil {
		return err
	}
	return os.WriteFile(s.accountDir, data, 0600)
}

// Register an account. Return true if success.
func (s *accountStorage) create(acc Account) bool {
	_, exist := s.accounts.LoadOrStore(acc.username, acc)
	return !exist
}

// Get an account. Return true if exist.
func (s *accountStorage) get(username string) (Account, bool) {
	val, exist := s.accounts.Load(username)
	acc := Account{}
	if exist {
		acc = val.(Account)
	}
	return acc, exist
}
