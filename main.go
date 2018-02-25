package main

func main() {
	// Start listener goroutine server from network
	done := make(chan bool)
	go startServer(done)
	<-done

	// Start user interface from frontend
	mainMenu()
}
