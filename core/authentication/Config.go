package authentication

import (
	"regexp"
)

const (
	accountFile     = "account.json"
	saltLength      = 32
	usernamePattern = `^[0-9A-Za-z\-_]+$`
)

var usernameValidator = regexp.MustCompile(usernamePattern)

func AccountFile() string {
	return accountFile
}
