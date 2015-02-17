package main

type PluginContainer struct {
	list  map[string]func(string) string
}
