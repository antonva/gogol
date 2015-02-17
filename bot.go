package main

import (
	"crypto/tls"
	"strconv"
	"strings"
	"log"
	"code.google.com/p/gcfg"
	"github.com/thoj/go-ircevent"
	"fmt"
)

func messageParser(e *irc.Event, conf Config, buffer *Msgbuf) string {
	msg := string(e.Arguments[1])
	nick := string(e.Nick)
	channel := string(e.Arguments[0])
	buffer.Append(nick, channel, msg)
	if len(msg) > 1 {
		if string(msg[0]) == conf.Bot.Trigger {
			split := strings.Split(msg[1:], " ")
			return split[0]
		/*TODO: split into own function/plugin.*/
		} else if string(msg[0:2]) == "s/" {
		}
	}
	fmt.Println(buffer.Buffer)
	return ""
}

func configParser(conf Config) {
}

func main() {
	var conf Config
	//configParser(conf)
	err := gcfg.ReadFileInto(&conf, "config")
	if err != nil {
		log.Fatal(err)
	}
	var plugin PluginContainer
	plugin.list = make(map[string]func(string) string)
	buffer :=  Msgbuf{Buflength: 512}

	con := irc.IRC(conf.Bot.Nick, conf.Bot.Name)
	con.VerboseCallbackHandler = conf.Bot.VerboseDebug
	con.Debug                  = conf.Bot.Debug
	con.Version                = "Pump-19 0.01. A go driven hydraulics golem"

	if conf.Bot.Ssl {
		con.UseTLS         = true
		con.TLSConfig      = &tls.Config{InsecureSkipVerify: true}
	}

	err = con.Connect(conf.Bot.Server+":"+strconv.Itoa(conf.Bot.Port))
	if err != nil {
		log.Fatal(err)
	}

	con.AddCallback("001", func(e *irc.Event) { con.Join("#neoplastic") })

	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		s := messageParser(e, conf, &buffer)
		if s != "" { con.Privmsg(e.Arguments[0], s) }
	})

	con.Loop()
}
