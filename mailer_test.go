package gomail

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ConfigData struct {
	Server   string
	Username string
	Password string
	From     string
	To       string
	Subject  string
	Msg      string
}

func TestSend(t *testing.T) {
	// get the params to setup the server
	config := &ConfigData{}
	c, err := ioutil.ReadFile("config.json")
	assert.Nil(t, err)
	assert.Nil(t, json.Unmarshal(c, config))
	t.Log(config.Server, config.Username, config.Password)

	// setup a mailer
	m := New(config.Server, config.Username, config.Password)
	assert.Nil(t, m.Send(config.From, config.To, config.Subject, config.Msg, []string{"steve@raffmail.net", "paul@cycle2u.com.au"}))

	// Check that the mail ended up in the target mailbox
}
