package piepan

import (
	"sort"
)

var (
	PluginExtensions []string
	Plugins          map[string]*Plugin = map[string]*Plugin{}
)

type Plugin struct {
	Name string
	New  func(*Instance) Environment
}

func Register(extension string, plugin *Plugin) {
	if plugin == nil {
		panic("piepan: plugin cannot be nil")
	}
	if _, ok := Plugins[extension]; ok {
		panic("piepan: extension " + extension + " already registered")
	}
	PluginExtensions = append(PluginExtensions, extension)
	sort.Strings(PluginExtensions)
	Plugins[extension] = plugin
}
