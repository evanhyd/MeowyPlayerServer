package account

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

// Replace with DB if hardware allows
type accountStorage struct {
	accountDir string
	accounts   accountMap
	fileLock   sync.Mutex
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
	s.fileLock.Lock()
	defer s.fileLock.Unlock()

	data, err := json.Marshal(s.accounts)
	if err != nil {
		return err
	}
	return os.WriteFile(s.accountDir, data, 0600)
}

// Register an account. Return true if success.
func (s *accountStorage) store(acc Account) bool {
	if _, exist := s.accounts.LoadOrStore(acc.Username, acc); exist {
		return false
	}

	if err := s.save(); err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Get an account. Return true if exist.
func (s *accountStorage) load(username string) (Account, bool) {
	val, exist := s.accounts.Load(username)
	acc := Account{}
	if exist {
		acc = val.(Account)
	}
	return acc, exist
}
