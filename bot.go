package main

import (
    "os/user"
	"crypto/tls"
	"strconv"
    "strings"
	"log"
	"time"
	"gopkg.in/gcfg.v1"
	"github.com/thoj/go-ircevent"
    "github.com/antonva/gogol/plugins"
)

func messageParser(e *irc.Event, conf Config, buffer *Msgbuf, p *plugins.Plugins) string {
	msg := string(e.Arguments[1])
	nick := string(e.Nick)
	channel := string(e.Arguments[0])
	buffer.Append(nick, channel, msg)
	if len(msg) > 1 {
		if string(msg[0]) == conf.Bot.Trigger {
			split := strings.Split(msg[1:], " ")
            reply := p.Call(split[0])
			return reply
        }
    }
	return ""
}

func loadPlugins() *plugins.Plugins {
    p := plugins.NewPlugins()
    p.Register()
    return p
}

// Looks for a config file in $HOME/.config/denton/config as per the XDG spec.
// TODO: custom flag for other config location.
// TODO: default fallback when no config is present or human readable error.
func loadConfig() Config {
    var usr, err = user.Current()
	var conf Config
	err = gcfg.ReadFileInto(&conf, usr.HomeDir + "/.config/gogol/config")
	if err != nil {
		log.Fatal(err)
	}
    return conf
}

func main() {

	buffer :=  Msgbuf{Buflength: 512}
    minversion := "0.1"
    conf := loadConfig()
    plugins := loadPlugins()

	con := irc.IRC(conf.Bot.Nick, conf.Bot.Name)
	con.VerboseCallbackHandler = conf.Bot.VerboseDebug
	con.Debug                  = conf.Bot.Debug
    con.Version                = "gogol - golang bot version: " + minversion

	if conf.Bot.Ssl {
		con.UseTLS         = true
		con.TLSConfig      = &tls.Config{InsecureSkipVerify: true}
	}

    err := con.Connect(conf.Bot.Server+":"+strconv.Itoa(conf.Bot.Port))
	if err != nil {
		log.Fatal(err)
	}

	con.AddCallback("001", func(e *irc.Event) {
        if conf.Bot.Password != "" {
            con.Privmsg("nickserv", "identify " + conf.Bot.Password)
        }
		time.Sleep( 2 * time.Second )
        for _, element := range conf.Bot.Channels {
            con.Join(element)
        }
	})

	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		go func() {
            s := messageParser(e, conf, &buffer, plugins)
            if s != "" { con.Privmsg(e.Arguments[0], s) }
        }()
	})

	con.Loop()
}
