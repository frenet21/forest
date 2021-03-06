package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/skratchdot/open-golang/open"
	"github.com/syndtr/goleveldb/leveldb"
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
★ Messages
  1. Send a message
  2. View received messages

★ Keys
  3. Manage recipients' public keys
  4. Manage personal private keys

★ Other
  5. Configuration
  6. Examine the Block Forest
  7. View source/issues/contributors (Opens Github)

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
		managePublicKeys()
	case "4":
		managePrivateKeys()
	case "5":
		config()
	case "6":
		examineForest()
	case "7":
		openGithub()
	case "x":
		clearScreen()
		fmt.Println("Exiting Forest. Goodbye.")
		os.Exit(0)
	default:
		fmt.Println("Unknown input. Try again.")
	}
	mainMenu()
}

func sendMessage() {

}

func viewReceived() {

}

func managePublicKeys() {
	clearScreen()
	printBanner()

	// Open the '.pubKeys' database. It is created if it does not exist.
	db, err := leveldb.OpenFile(".pubKeys", nil)
	if err != nil {
		panic("Could not open public key database file.")
	}

	fmt.Println(`
★ Main Menu -> Manage public keys ★
List of known recipient public keys:
`)

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		// Grab the pubKey (key) and name (value) from the local database
		pubKey := iter.Key()
		name := iter.Value()

		// String the byte arrays
		pubKeyString := string(pubKey[:])
		nameString := string(name[:])

		// Print the strings
		fmt.Print("NAME: " + nameString)
		fmt.Println("KEY: " + pubKeyString)
	}
	iter.Release()
	err = iter.Error()

	fmt.Println(`
1. Add new recipient public keys
2. Remove recipient public keys
3. <- Return to main menu
`)

	fmt.Print("Make selection > ")
	var selection int
	fmt.Scan(&selection)

	switch selection {
	case 1:
		fmt.Println("You will be prompted to add the recipient's public key and name it.")

		// User pastes in the public key and names it
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Paste public key > ")
		pubKey, _ := reader.ReadString('\n')
		fmt.Print("Give a name to this key > ")
		name, _ := reader.ReadString('\n')

		// Key and name are put into the 'pubKeys' database
		// 'pubKey' is the key, 'name' is the name in case of duplicates
		err = db.Put([]byte(pubKey), []byte(name), nil)
		if err != nil {
			panic("Failed to write new public key to database.")
		}
	case 2:
		// User pastes in the public key and names it
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Paste in the public key you wish to remove >")
		pubKey, _ := reader.ReadString('\n')

		data, err := db.Get([]byte(pubKey), nil)
		if err != nil {
			panic("Unknown key.")
		}
		// If the key matches, string it for printing
		dataString := string(data[:])
		// Attempt to delete the key from the database
		fmt.Print("Deleting " + dataString + "...")
		err = db.Delete([]byte(pubKey), nil)
		fmt.Print("Deleted key.")
	case 3:
		db.Close()
		return
	default:
		managePublicKeys()
	}
	db.Close()
	managePublicKeys()
}

func managePrivateKeys() {
	clearScreen()
	printBanner()

	// Open the '.pubKeys' database. It is created if it does not exist.
	db, err := leveldb.OpenFile(".priKeys", nil)
	if err != nil {
		panic("Could not open public key database file.")
	}

	fmt.Println(`
★ Main Menu -> Manage private keys ★
1. Show my private keys (Careful!)
2. Add a new private key
3. Remove a private key
4. <- Return to main menu
`)

	fmt.Print("Make selection > ")
	var selection int
	fmt.Scan(&selection)

	switch selection {
	case 1:
		fmt.Println("! ! ! Displaying keys ! ! ! ")
		iter := db.NewIterator(nil, nil)
		for iter.Next() {
			// Grab the priKey (key) and name (value) from the local database
			priKey := iter.Key()
			name := iter.Value()

			// String the byte arrays
			priKeyString := string(priKey[:])
			nameString := string(name[:])

			// Print the strings
			fmt.Print("NAME: " + nameString)
			fmt.Println("KEY: " + priKeyString)
		}
		iter.Release()
		err = iter.Error()

		fmt.Print("Press enter to hide keys...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	case 2:
		fmt.Println("You will be prompted to add your private key name it.")

		// User pastes in their private key and names it
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Paste private key > ")
		priKey, _ := reader.ReadString('\n')
		fmt.Print("Give a name to this key > ")
		name, _ := reader.ReadString('\n')

		// Key and name are put into the 'pubKeys' database
		// 'pubKey' is the key, 'name' is the name in case of duplicates
		err = db.Put([]byte(priKey), []byte(name), nil)
		if err != nil {
			panic("Failed to write new private key to database.")
		}
	case 3:
		// User pastes in the private key and names it
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Paste in the private key you wish to remove >")
		priKey, _ := reader.ReadString('\n')

		data, err := db.Get([]byte(priKey), nil)
		if err != nil {
			panic("Unknown key.")
		}
		// If the key matches, string it for printing
		dataString := string(data[:])
		// Attempt to delete the key from the database
		fmt.Print("Deleting " + dataString + "...")
		err = db.Delete([]byte(priKey), nil)
		if err != nil {
			log.Println("Could not delete key.")
		}
		fmt.Print("Deleted key.")
	case 4:
		db.Close()
		return
	default:
		managePrivateKeys()
	}
	db.Close()
	managePrivateKeys()
}

func config() {

}

func examineForest() {

}

func openGithub() {
	open.Run("https://github.com/stellar-tech/forest")
}

/* Exported functions */

/* Checks the .knownHashes database to see if it a hash exists (returns true/false) */
func CheckKnownHashes(blockID string) bool {
	db, err := leveldb.OpenFile(".knownHashes", nil)
	if err != nil {
		panic("Could not open the known hash database.")
	}

	// Set results into data
	data, err := db.Get([]byte(blockID), nil)
	db.Close()

	if len(data) > 0 {
		return true
	} else {
		return false
	}
}

/* Receives a blockID and adds it to the .knownHashes database */
func AddNewHash(blockID string) {
	db, err := leveldb.OpenFile(".knownHashes", nil)
	if err != nil {
		panic("Could not open the known hash database.")
	}

	err = db.Put([]byte(blockID), []byte(nil), nil)
	db.Close()
}

/*
func GetClientList() {
	db, err := leveldb.OpenFile(".knownClients", nil)
	if err != nil {
		panic("Could not open the known client list.")
	}

	// Set results into data

	// Todo: Get all keys in data
	data, err := db.Get([]byte(""), nil)
	db.Close()
	fmt.Print(data)
}
*/
