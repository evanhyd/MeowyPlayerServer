package authentication

import (
	"regexp"
)

const (
	accountFile = "account.json"
	saltLength  = 32
	idPattern   = `^[0-9A-Za-z\-_]+$`
)

var idValidator = regexp.MustCompile(idPattern)

func AccountFile() string {
	return accountFile
}
