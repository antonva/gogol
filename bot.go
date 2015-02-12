package main

import ( 
	"github.com/thoj/go-ircevent"
	"crypto/tls"
	"strconv"
	"fmt"
)

type config struct {
	server string
	ssl    bool
	port   int
	nick   string
	name   string
	version string
	trigger string
}

func messageParser(msg string, conf config) string {
	if len(msg) > 0 && string(msg[0]) == conf.trigger {
		return "TRIGGER WARNING"
	}

	return ""
}

func main() {
	conf := config {
		server:  "irc.rizon.net",
		ssl:     true,
		port:    6697,
		nick:    "Mr-Pump",
		name:    "Pump-19",
		version: "Pump-19 0.01. A go driven hydraulics golem",
		trigger: ".",
	}

	bot := irc.IRC(conf.nick, conf.name)
	bot.VerboseCallbackHandler = true
	bot.Debug                  = true
	bot.Version                = conf.version
	if conf.ssl {
		bot.UseTLS         = true
		bot.TLSConfig      = &tls.Config{InsecureSkipVerify: true}
	}

	err := bot.Connect(conf.server+":"+strconv.Itoa(conf.port))
	if err != nil {
		fmt.Println(err.Error())
	}

	bot.AddCallback("001", func(e *irc.Event) { bot.Join("#neoplastic") })

	bot.AddCallback("PRIVMSG", func(e *irc.Event) {
		msg := string(e.Arguments[1])
		go func() {
			s := messageParser(msg, conf)
			if s != "" { bot.Privmsg(e.Arguments[0], s) }
		}()
	})

	bot.Loop()
}
