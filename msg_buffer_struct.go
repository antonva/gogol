package main

import (
	"fmt"
)

type Msg struct {
	Nick    string
	Channel string
	Message string
}

type Msgbuf struct {
	Buffer    []Msg
	Buflength int
}

func (m *Msgbuf) Append(nick string, channel string, message string) {
	var tmp Msg
	tmp.Nick    = nick
	tmp.Channel = channel
	tmp.Message = message
	if len(m.Buffer) == m.Buflength {
		m.Buffer = m.Buffer[1:]
	}
	m.Buffer = append(m.Buffer, tmp)
	fmt.Println(m.Buffer)
}

func (m Msg) String() string {
	return "Nick: " + m.Nick + " Channel: " + m.Channel + " Msg: " + m.Message
}

