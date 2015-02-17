package main


type Config struct {
	Bot struct {
		Debug        bool
		Ssl          bool
		VerboseDebug bool
		Port         int
		Name         string
		Nick         string
		Server       string
		Trigger      string
	}
}
