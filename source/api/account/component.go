package account

import (
	"crypto/rand"
	"crypto/sha256"
	"log"
	"slices"

	"github.com/google/uuid"
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

func (c *Component) computeHash(password []byte, salt []byte) []byte {
	hash := sha256.Sum256(append(password, salt...))
	return hash[:]
}

func (c *Component) Authorize(username string, password string) bool {
	if !c.isValidUsername(username) {
		return false
	}

	acc, exist := c.storage.load(username)
	if !exist {
		return false
	}
	computedHash := c.computeHash([]byte(password), acc.salt)
	return slices.Equal(acc.hash, computedHash)
}

func (c *Component) Register(username string, password string) bool {
	if !c.isValidUsername(username) {
		return false
	}

	acc := Account{username: username, id: uuid.NewString()}
	acc.salt = make([]byte, c.saltLength)
	if _, err := rand.Read(acc.salt); err != nil {
		log.Println(err)
		return false
	}
	acc.hash = c.computeHash([]byte(password), acc.salt)
	return c.storage.store(acc)
}
