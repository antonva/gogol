package main

import (
    "fmt"
    "os"
    "os/user"
    "crypto/tls"
    "strconv"
    "strings"
    "log"
    "time"
    "encoding/json"
    "github.com/thoj/go-ircevent" // Consider eventually writing your own.
    "github.com/antonva/gogol/plugins"
)

func messageParser(e *irc.Event, conf Config, buffer *Msgbuf, p *plugins.Plugins) []string {
    msg := string(e.Arguments[1])
    nick := string(e.Nick)
    channel := string(e.Arguments[0])
    buffer.Append(nick, channel, msg)
    if len(msg) > 1 {
        if string(msg[0]) == conf.Bot.Trigger || string(msg[0]) == ".bots" {
            split := strings.Split(msg[1:], " ")
            reply := p.Call(split)
            return reply
        }
    }
    return []string{}
}

// Looks for a config file in $HOME/.config/gogol/config as per the XDG spec.
// TODO: custom flag for other config location.
// TODO: default fallback when no config is present or human readable error.
// TODO: Replace gcfg with something else.
func loadConfig() Config {
    fmt.Println("conf")
    var usr, err = user.Current()
    file, _ := os.Open(usr.HomeDir + "/.config/gogol/config.json")
    decoder := json.NewDecoder(file)
    conf := Config{}
    err = decoder.Decode(&conf)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(conf)
    return conf
}

func main() {

    buffer :=  Msgbuf{Buflength: 512}
    minversion := "0.2"
    conf := loadConfig()
    plugins := plugins.NewPlugins()

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
            if len(s) > 0 {
                for _, el := range s {
                    con.Privmsg(e.Arguments[0], el)
                    time.Sleep( time.Second / 2 )
                }
            }
        }()
    })

    con.Loop()
}
