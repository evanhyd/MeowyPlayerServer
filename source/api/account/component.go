package account

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"regexp"
	"slices"

	"github.com/google/uuid"
)

type Component struct {
	storage           accountStorage
	usernameValidator *regexp.Regexp
	saltLength        int
}

func MakeComponent() Component {
	const (
		kUsernamePattern = `^[0-9A-Za-z\-_]+$`
		kSaltLength      = 32
	)
	return Component{
		storage:           makeStorage(),
		usernameValidator: regexp.MustCompile(kUsernamePattern),
		saltLength:        kSaltLength,
	}
}

func (c *Component) Initialize() error {
	return c.storage.initialize()
}

func (c *Component) isValidUsername(username string) bool {
	return c.usernameValidator.MatchString(username)
}

func (c *Component) computeHash(password []byte, salt []byte) []byte {
	hash := sha256.Sum256(append(password, salt...))
	return hash[:]
}

func (c *Component) Authorize(username string, password string) bool {
	acc, exist := c.storage.get(username)
	if !exist {
		return false
	}
	computedHash := c.computeHash([]byte(password), acc.salt)
	return slices.Equal(acc.hash, computedHash)
}

func (c *Component) Register(username string, password string) error {
	if !c.isValidUsername(username) {
		return fmt.Errorf("username can only contain letters, numbers, underscore or dash")
	}

	acc := Account{username: username, id: uuid.NewString()}
	acc.salt = make([]byte, c.saltLength)
	if _, err := rand.Read(acc.salt); err != nil {
		return err
	}
	acc.hash = c.computeHash([]byte(password), acc.salt)

	if !c.storage.create(acc) {
		return fmt.Errorf("username already exists")
	}
	return nil
}
