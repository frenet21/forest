package main

import (
	"fmt"
	"os"
)

func main(){
	offset := true
	for offset {
		fmt.Println("")
		fmt.Println("★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★")
		fmt.Printf("★★★★★★ Welcome to Forest! You can begin your chatting now ★★★★★★\n")
		fmt.Println("★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★")
		fmt.Printf("1. Write new message;\n")
		fmt.Printf("2. Store a public key;\n")
		fmt.Printf("3. Exit Forest. \n")

		fmt.Print("★★★(Enter any number of options above)\n>>")

		//get input order from user
		var order int
		fmt.Scan(&order)

		//fmt.Println(reflect.TypeOf(order))

		switch order {
		case 1:
			fmt.Printf("You entered 1\n")
			newM()
	
		case 2: 	
			fmt.Printf("You entered 2\n")
			fmt.Print("Enter the public key you want to store: \n>>")
			var pubKey string
			fmt.Scan(&pubKey)
			fmt.Print("Enter the name of this public key: \n>>")
			var uName string
			fmt.Scan(&uName)
			fmt.Print("The public key you entered is: "+pubKey+", the name you entered is: "+uName+"\n")
			storeAPublicKey(pubKey, uName)
	
		case 3: 
			fmt.Printf("★Thank you for using Forest, you already exited ★\n")
			offset = false
		}
	}
}

/*To-Do: write a new message*/
func newM(){
	fmt.Println("func!!!!!!!!!!!")
}

/* Used by storeAPublicKey function */
var (
	fileInfo *os.FileInfo
	err error
)
/*To-Do: storeAPublicKey is storing public_key-user_Name pairs into a txt file, 
so that newM() function could use to select receivers from one of the list.*/
func storeAPublicKey(publicKey string, userName string){
	/*CHECK if the test.txt is already existed, if not, create the file. */
	fileInfo, err := os.Stat("test.txt")
    if err != nil {
        newFile, err := os.Create("test.txt")
    		if err != nil {
        		fmt.Println(err)
			}
			newFile.Close()
	} else {
		fmt.Println("File does exist. File information:")
		fmt.Println(fileInfo)
	}

	/*write public key and user name into the file. */
	file, err := os.OpenFile(
        "test.txt",
        os.O_WRONLY|os.O_APPEND|os.O_CREATE,
        0666,
    )
    if err != nil {
        fmt.Println(err)
    }
    defer file.Close()

	li, err := file.WriteString(publicKey+","+userName+"\n")
    fmt.Printf("wrote %d into the file.\n", li)

}

