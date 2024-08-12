package account

import (
	"crypto/rand"
	"crypto/sha256"
	"log"
	"slices"
)

type Component struct {
	storage    accountStorage
	saltLength int
}

func MakeComponent() Component {
	const kSaltLength = 32
	return Component{
		storage:    makeStorage(),
		saltLength: kSaltLength,
	}
}

func (c *Component) Initialize() error {
	return c.storage.initialize()
}

func (c *Component) isValidUsername(username string) bool {
	const (
		kMinLen = 1
		kMaxLen = 24
	)
	return kMinLen <= len(username) && len(username) <= kMaxLen
}

func (c *Component) generateHash(password string) ([]byte, []byte) {
	salt := make([]byte, c.saltLength)
	if _, err := rand.Read(salt); err != nil {
		log.Println(err)
	}
	return c.computeHash([]byte(password), salt), salt
}

func (c *Component) computeHash(password []byte, salt []byte) []byte {
	hash := sha256.Sum256(append(password, salt...))
	return hash[:]
}

func (c *Component) Authenticate(username string, password string) (UserID, bool) {
	if !c.isValidUsername(username) {
		return UserID{}, false
	}

	acc, exist := c.storage.load(username)
	if !exist {
		return UserID{}, false
	}
	computedHash := c.computeHash([]byte(password), acc.Salt)
	return acc.UserID, slices.Equal(acc.Hash, computedHash)
}

func (c *Component) Register(username string, password string) bool {
	if !c.isValidUsername(username) {
		return false
	}

	acc := Account{Username: username, UserID: NewUserID()}
	acc.Hash, acc.Salt = c.generateHash(password)
	return c.storage.store(acc)
}
