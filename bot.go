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
	fullserver string
}

func messageParser(msg string, conf config) string {
	if string(msg[0]) == conf.trigger {
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
	conf.fullserver = conf.server+":"+strconv.Itoa(conf.port);

	bot := irc.IRC(conf.nick, conf.name)
	bot.VerboseCallbackHandler = true
	bot.Debug                  = true
	if conf.ssl {
		bot.UseTLS                 = true
		bot.TLSConfig              = &tls.Config{InsecureSkipVerify: true}
	}
	bot.Version                = conf.version

	err := bot.Connect(conf.fullserver)
	if err != nil {
		fmt.Println(err.Error())
	}

	bot.AddCallback("001", func(e *irc.Event) { bot.Join("#neoplastic") })

	bot.AddCallback("PRIVMSG", func(e *irc.Event) {
		parser := string(e.Arguments[1])
		//slice  := []byte(parser)
		go func() {
			s := messageParser(parser, conf)
			if s != "" { bot.Privmsg(e.Arguments[0], s) }
		}()
	})

	bot.Loop()
}
