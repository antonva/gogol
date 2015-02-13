package main

import (
	"code.google.com/p/gcfg"
	"github.com/thoj/go-ircevent"
	"crypto/tls"
	"strconv"
	"fmt"
)

func messageParser(msg string, conf Config) string {
	if len(msg) > 0 && string(msg[0]) == conf.Bot.Trigger {
		return "TRIGGER WARNING"
	}

	return ""
}

func main() {
	var conf Config
	err := gcfg.ReadFileInto(&conf, "config")
	if err != nil {
		fmt.Println(err)
	}

	bot := irc.IRC(conf.Bot.Nick, conf.Bot.Name)
	bot.VerboseCallbackHandler = true
	bot.Debug                  = true
	bot.Version                = "Pump-19 0.01. A go driven hydraulics golem"
	if conf.Bot.Ssl {
		bot.UseTLS         = true
		bot.TLSConfig      = &tls.Config{InsecureSkipVerify: true}
	}

	err = bot.Connect(conf.Bot.Server+":"+strconv.Itoa(conf.Bot.Port))
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
