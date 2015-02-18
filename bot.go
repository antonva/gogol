package main

import (
	"crypto/tls"
	"strconv"
	"strings"
	"log"
	"regexp"
	"time"
	"code.google.com/p/gcfg"
	"github.com/thoj/go-ircevent"
	"github.com/aarzilli/golua/lua"
)

func messageParser(e *irc.Event, conf Config, buffer *Msgbuf) string {
	msg := string(e.Arguments[1])
	nick := string(e.Nick)
	channel := string(e.Arguments[0])
	buffer.Append(nick, channel, msg)
	if len(msg) > 1 {
		if string(msg[0]) == conf.Bot.Trigger {
			//split := strings.Split(msg[1:], " ")
			//return split[0]
	/*
	 * Sed style replace on last line spoken
	 * TODO: refactor so you can correct others.
	 */
		} else if string(msg[0:2]) == "s/" {
			for i := len(buffer.Buffer)-2; i >= 0; i-- {
				if nick == buffer.Buffer[i].Nick {
					return "I think you meant: " +
					sed(msg[2:], buffer.Buffer[i].Message)
				}
			}
		}
	}
	return ""
}

func sed(rgx string,msg string) string {
	split_slice := strings.Split(rgx, "/")
	if len(split_slice) > 1 {
		//TODO: Make it work inline and not just for one off \/'s.
		if split_slice[0] == "\\" { split_slice[0] = "/" }
		re := regexp.MustCompile(split_slice[0])
		return re.ReplaceAllString(msg, split_slice[1])
	}
	return ""
}


func main() {
	/* 
	 * The bot's message buffer.
	 * TODO: Move Buflength to config.
	 */
	buffer :=  Msgbuf{Buflength: 512}

	/* Config parsing */
	var conf Config
	err := gcfg.ReadFileInto(&conf, "config")
	if err != nil {
		log.Fatal(err)
	}

	/* Plugin */
	L := lua.NewState()
	L.OpenLibs()
	defer L.Close()

	err = L.DoString("print(\"Hello world\")")
	var plugin PluginContainer
	plugin.list = make(map[string]func(string) string)

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

	con.AddCallback("001", func(e *irc.Event) {
		con.Privmsg("nickserv", "identify pump19")
		time.Sleep( 2 * time.Second )
		con.Join("#frost")
	})

	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		s := messageParser(e, conf, &buffer)
		if s != "" { con.Privmsg(e.Arguments[0], s) }
	})

	con.Loop()
}
