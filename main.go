package main

import (
	"fmt"
	"os"
	//"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

const masterPassFile = "masterpass.txt"

var passwords map[string]string

func main() {
	// Check if the master passphrase has already been set
	if _, err := os.Stat(masterPassFile); err == nil {
		passwords = make(map[string]string)

		// Initialize or load passwords from a file/database

		// Authenticate user
		fmt.Print("Enter master passphrase: ")
		passphrase, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("Failed to read passphrase:", err)
			return
		}
		fmt.Println()

		// Read the stored hashed passphrase from file
		hashedPassphrase, err := os.ReadFile(masterPassFile)
		if err != nil {
			fmt.Println("Failed to read stored master passphrase:", err)
			return
		}

		// Compare the entered passphrase with the stored hashed passphrase
		err = bcrypt.CompareHashAndPassword(hashedPassphrase, passphrase)
		if err != nil {
			fmt.Println("Invalid master passphrase")
			return
		}

		menu()
	} else {
		// Master passphrase not set, prompt user to set it
		setMasterPassphrase()
	}
}
func menu() {
	// Present menu to the user
	fmt.Println("Welcome to the Password Manager CLI")
	fmt.Println("1. Add Password")
	fmt.Println("2. Retrieve Password")
	fmt.Println("3. Exit")

	var choice int
	fmt.Print("Enter your choice: ")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		addPassword()
	case 2:
		retrievePassword()
	case 3:
		os.Exit(0)
	default:
		fmt.Println("Invalid choice")
	}
}

// func authenticate(passphrase string) bool {
// 	// Check if passphrase matches the stored hash
// 	// You should store and compare the hashed passphrase securely
// 	storedHash := "$2a$10$H04PTgg3T8k.1O5nUDV3V.PajivucrJF0bQ51o7M6oCQ9ZVLwmt1O" // Example hash
// 	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(passphrase))
// 	return err == nil
// }

func addPassword() {
	var name string
	fmt.Print("Enter name for the password: ")
	fmt.Scanln(&name)
	fmt.Print("Enter password: ")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Password added failed")
		menu()
	}
	fmt.Println()

	// Store the password securely
	passwords[name] = string(password)
	fmt.Println("Password added successfully")
	menu()
}

func retrievePassword() {
	var name string
	fmt.Print("Enter name of the password to retrieve: ")
	fmt.Scanln(&name)

	// Retrieve and display the password
	password, ok := passwords[name]
	if !ok {
		fmt.Println("Password not found")
		menu()
	}

	fmt.Println("Password:", password)
	menu()
}

func setMasterPassphrase() {
	fmt.Println("Welcome to the Password Manager CLI")

	// Prompt the user to enter the master passphrase
	fmt.Print("Enter master passphrase: ")
	passphrase, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Failed to read passphrase:", err)
		return
	}
	fmt.Println()

	// Hash the master passphrase
	hashedPassphrase, err := bcrypt.GenerateFromPassword(passphrase, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Failed to hash passphrase:", err)
		return
	}

	// Store the hashed passphrase in a file
	err = os.WriteFile(masterPassFile, hashedPassphrase, 0600)
	if err != nil {
		fmt.Println("Failed to store master passphrase:", err)
		return
	}

	fmt.Println("Master passphrase set successfully!")
	menu()
}
