// Package savestates takes care of serializing and unserializing the game RAM
// to the host filesystem.
package savestates

import (
	"io/ioutil"
	"os/user"
	"path/filepath"

	"github.com/libretro/ludo/state"
)

func name() string {
	name := filepath.Base(state.Global.GamePath)
	ext := filepath.Ext(name)
	name = name[0 : len(name)-len(ext)]
	return name + ".state"
}

// Save the current state to the filesystem
func Save() error {
	usr, _ := user.Current()
	s := state.Global.Core.SerializeSize()
	bytes, err := state.Global.Core.Serialize(s)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(usr.HomeDir+"/.ludo/savestates/"+name(), bytes, 0644)
}

// Load the state from the filesystem
func Load() error {
	usr, _ := user.Current()
	s := state.Global.Core.SerializeSize()
	bytes, err := ioutil.ReadFile(usr.HomeDir + "/.ludo/savestates/" + name())
	if err != nil {
		return err
	}
	err = state.Global.Core.Unserialize(bytes, s)
	return err
}
