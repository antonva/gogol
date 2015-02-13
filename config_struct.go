package main


type Config struct {
	Bot struct {
		Server string
		Ssl    bool
		Port   int
		Nick   string
		Name   string
		Trigger string
	}
}
