package credentials

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
)

type Credentials struct {
	Username string
	Password string
}

type Ask int

const (
	NeverAsk Ask = iota
	CanAsk
	MustAsk
)

func GetUsernamePassword(ask Ask, usernameEnvar, passwordEnvar string) (*Credentials, error) {

	stdin := int(syscall.Stdin)
	istty := terminal.IsTerminal(stdin)

	username := os.Getenv(usernameEnvar)
	if ask == MustAsk || ask == CanAsk && username == "" && istty {
		fmt.Print("Enter username: ")
		bytes, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return nil, errors.Wrap(err, "reading username from terminal")
		}
		username = strings.TrimRight(string(bytes), "\n\r")
	}

	password := os.Getenv(passwordEnvar)
	if ask == MustAsk || ask == CanAsk && password == "" && istty {
		fmt.Print("Enter password: ")
		bytes, err := terminal.ReadPassword(stdin)
		if err != nil {
			return nil, errors.Wrap(err, "reading password from terminal")
		}
		fmt.Println()
		password = strings.TrimRight(string(bytes), "\n\r")
	}

	return &Credentials{Username: username, Password: password}, nil

}
