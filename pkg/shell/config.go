package shell

import (
	"log"
	"os"
	"os/user"
	"strings"
)

type Config struct {
	CurrentDir   string
	CurrentUser  *user.User
	CurrentGroup *user.Group
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	checkError(err)
	return dir
}

func NewConfig() (result *Config) {

	currUsr, err := user.Current()
	checkError(err)
	group, err := user.LookupGroupId(currUsr.Gid)
	checkError(err)

	result = &Config{
		CurrentDir:   getCurrentDir(),
		CurrentUser:  currUsr,
		CurrentGroup: group,
	}

	return result
}

func (c *Config) UiString() string {
	var b strings.Builder

	b.WriteString(c.CurrentUser.Name)
	b.WriteRune('@')
	b.WriteString(c.CurrentGroup.Name)
	b.WriteRune(':')
	b.WriteString(c.CurrentDir)
	b.WriteRune('$')

	return b.String()
}
