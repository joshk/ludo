package menu

import (
	"io/ioutil"
	"path/filepath"

	"github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/utils"
)

type screenExplorer struct {
	entry
}

func buildExplorer(path string, exts []string, cb func(string), dirAction *entry) Scene {
	var list screenExplorer
	list.label = "Explorer"

	files, err := ioutil.ReadDir(path)
	if err != nil {
		notifications.DisplayAndLog("Menu", err.Error())
	}

	if dirAction != nil && dirAction.label != "" {
		dirAction.callbackOK = func() { cb(path) }
		list.children = append(list.children, *dirAction)
	}

	for _, f := range files {
		f := f
		icon := "file"

		// Check whether or not we are to display hidden files.
		if f.Name()[:1] == "." && settings.Settings.ShowHiddenFiles {
			continue
		}

		if f.IsDir() {
			icon = "folder"
		}
		list.children = append(list.children, entry{
			label: f.Name(),
			icon:  icon,
			callbackOK: func() {
				if f.IsDir() {
					list.segueNext()
					newPath := path + "/" + f.Name()
					if dirAction != nil {
						dirAction.callbackOK = func() { cb(newPath) }
					}
					menu.stack = append(menu.stack, buildExplorer(newPath, exts, cb, dirAction))
				} else if cb != nil && (exts == nil || utils.StringInSlice(filepath.Ext(f.Name()), exts)) {
					cb(path + "/" + f.Name())
				}
			},
		})
	}

	if len(files) == 0 {
		list.children = append(list.children, entry{
			label: "Empty",
			icon:  "subsetting",
		})
	}

	list.segueMount()

	return &list
}

func (explorer *screenExplorer) Entry() *entry {
	return &explorer.entry
}

func (explorer *screenExplorer) segueMount() {
	genericSegueMount(&explorer.entry)
}

func (explorer *screenExplorer) segueNext() {
	genericSegueNext(&explorer.entry)
}

func (explorer *screenExplorer) segueBack() {
	genericAnimate(&explorer.entry)
}

func (explorer *screenExplorer) update() {
	genericInput(&explorer.entry)
}

func (explorer *screenExplorer) render() {
	genericRender(&explorer.entry)
}
