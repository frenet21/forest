package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/skratchdot/open-golang/open"
)

func clearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func printBanner() {
	fmt.Println(`
         _-_                 _-_                 _-_
      /~~   ~~\           /~~   ~~\           /~~   ~~\
   /~~         ~~\     /~~         ~~\     /~~         ~~\
  {               }   {               }   {               }
   \  _-     -_  /     \  _-     -_  /     \  _-     -_  /
     ~  \\ //  ~         ~  \\ //  ~         ~  \\ //  ~
  _- -   | | _- _     _- -   | | _- _     _- -   | | _- _
    _ -  | |   -_       _ -  | |   -_       _ -  | |   -_
        // \\               // \\               // \\
 ===========================================================
                           FOREST
   Blockchain Distributed Hyper-Secure Encrypted Messenger
                     StellarTech ★ 2018 
 ===========================================================`)
}

func mainMenu() {
	clearScreen()
	printBanner()

	fmt.Println(`     
★ Main Menu ★
1. Send a message
2. View received messages
3. Add a recipient's public key
4. Add a personal private key
5. View more options ->

x. Exit Forest
`)

	fmt.Print("Make selection > ")
	var selection string
	fmt.Scan(&selection)

	switch selection {
	case "1":
		sendMessage()
	case "2":
		viewReceived()
	case "3":
		addPublicKey()
	case "4":
		addPrivateKey()
	case "5":
		moreOptions()
	case "x":
		clearScreen()
		fmt.Println("Exiting Forest. Goodbye.")
		os.Exit(0)
	default:
		fmt.Println("Unknown input. Try again.")
	}
	mainMenu()
}

func moreOptions() {
	clearScreen()
	printBanner()

	fmt.Println(`     
★ More Options... ★
1. Configuration
2. Examine the Block Forest
3. View source/issues/contributors (Opens Github)
4. <- Return to Main Menu
`)

	fmt.Print("Make selection > ")
	var selection int
	fmt.Scan(&selection)

	switch selection {
	/* Todo: add functions for secondary menu */
	case 1:
		configuration()
	case 2:
		examineForest()
	case 3:
		openGithub()
	case 4:
		mainMenu()
	default:
		fmt.Println("Unknown input. Try again.")
	}
	moreOptions()
}

func sendMessage() {

}

func viewReceived() {

}

func addPublicKey() {

}

func addPrivateKey() {

}

func configuration() {

}

func examineForest() {

}

func openGithub() {
	open.Run("https://github.com/stellar-tech/forest")
}
