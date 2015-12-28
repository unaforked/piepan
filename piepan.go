package piepan

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/unascribed/gumble/gumble"
	"github.com/unascribed/gumble/gumble_ffmpeg"
)

type Environment interface {
	gumble.EventListener
	LoadScriptFile(filename string) error
}

type Instance struct {
	Client *gumble.Client
	Audio  *gumble_ffmpeg.Stream

	envs map[string]Environment
}

func New(client *gumble.Client) *Instance {
	in := &Instance{
		Client: client,
		envs:   make(map[string]Environment),
	}
	return in
}

// [type:[environment:]]filename
func (in *Instance) LoadScript(name string) error {
	var filename, filetype, environment string
	pieces := filepath.SplitList(name)
	switch len(pieces) {
	case 1:
		filename = pieces[0]
		for _, ext := range PluginExtensions {
			if strings.HasSuffix(filename, ext) {
				filetype = ext
				break
			}
		}
		environment = filetype
	case 2:
		filename = pieces[1]
		filetype = pieces[0]
		environment = filetype
	case 3:
		filename = pieces[2]
		filetype = pieces[0]
		environment = pieces[1]
	default:
		return errors.New("unknown script name format")
	}
	plugin := Plugins[filetype]
	if plugin == nil {
		return errors.New("unknown filetype")
	}
	env := in.envs[environment]
	if env == nil {
		env = plugin.New(in)
		in.Client.Attach(env)
		in.envs[environment] = env
	}
	if err := env.LoadScriptFile(filename); err != nil {
		return err
	}
	return nil
}
