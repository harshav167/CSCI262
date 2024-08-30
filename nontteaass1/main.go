package main

import (
	"ass1/hash"
	"ass1/models"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Handle the -i flag for user/password creation
	if len(os.Args) > 1 && os.Args[1] == "-i" {
		err := models.InitializeUser()
		if err != nil {
			log.Fatalf("Failed to initialize user: %v", err)
		}
		return
	}

	// Test MD5 hashing function
	md5Test := hash.MD5Hash("This is a test")
	fmt.Printf("MD5 (\"This is a test\") = %s\n", md5Test)

	// Signal channel to catch Ctrl + C
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("\nExiting program...")
		os.Exit(0)
	}()

	// Initialize the model
	model := models.InitialModel()

	// Load existing file system from Files.store
	if err := model.LoadFileSystem(); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No existing file system found, starting with an empty file system.")
		} else {
			log.Fatalf("Failed to load file system: %v", err)
		}
	}

	// Start the program loop
	model.Start()
}
