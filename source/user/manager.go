package user

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"net/http"
	"regexp"
	"slices"

	"github.com/google/uuid"
)

type userManager struct {
	storage           *accountStorage
	usernameValidator *regexp.Regexp
	saltLength        int
}

func NewUserManager() *userManager {
	const (
		kUsernamePattern = `^[0-9A-Za-z\-_]+$`
		kSaltLength      = 32
	)
	return &userManager{
		storage:           newAccountStorage(),
		usernameValidator: regexp.MustCompile(kUsernamePattern),
		saltLength:        kSaltLength,
	}
}

func (m *userManager) initialize() error {
	return m.storage.initialize()
}

func (m *userManager) isValidUsername(username string) bool {
	return m.usernameValidator.MatchString(username)
}

func (m *userManager) computeHash(password []byte, salt []byte) []byte {
	hash := sha256.Sum256(append(password, salt...))
	return hash[:]
}

func (m *userManager) authorize(username string, password string) bool {
	acc, exist := m.storage.get(username)
	if !exist {
		return false
	}
	computedHash := m.computeHash([]byte(password), acc.salt)
	return slices.Equal(acc.hash, computedHash)
}

func (m *userManager) register(username string, password string) error {
	if !m.isValidUsername(username) {
		return fmt.Errorf("username can only contain letters, numbers, underscore or dash")
	}

	acc := account{username: username, id: uuid.NewString()}
	acc.salt = make([]byte, m.saltLength)
	if _, err := rand.Read(acc.salt); err != nil {
		return err
	}
	acc.hash = m.computeHash([]byte(password), acc.salt)

	if !m.storage.create(acc) {
		return fmt.Errorf("username already exists")
	}
	return nil
}

func (m *userManager) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.User.Username()
	password, ok := r.URL.User.Password()
	if !ok {
		fmt.Fprintf(w, "missing password field")
		return
	}
	fmt.Fprintf(w, "%v", m.register(username, password))
}
