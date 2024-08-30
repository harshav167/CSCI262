package auth

import (
	"ass1/hash"
	"errors"
	"fmt"
	"os"
	"strings"
)

func AuthenticateUser(username, password string) (bool, int, error) {
	salt, err := getSalt(username)
	if err != nil {
		return false, 0, err
	}

	passSaltHash := hash.MD5Hash(password + salt)

	clearance, err := getClearance(username, passSaltHash)
	if err != nil {
		return false, 0, err
	}

	fmt.Println("Authentication for user", username, "complete.")
	fmt.Println("The clearance for", username, "is", clearance)

	return true, clearance, nil
}

func getSalt(username string) (string, error) {
	file, err := os.ReadFile("salt.txt")
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, username+":") {
			return strings.Split(line, ":")[1], nil
		}
	}

	return "", errors.New("username not found in salt.txt")
}

func getClearance(username, passSaltHash string) (int, error) {
	file, err := os.ReadFile("shadow.txt")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if parts[0] == username && parts[1] == passSaltHash {
			clearance := parts[2]
			return int(clearance[0] - '0'), nil
		}
	}

	return 0, errors.New("authentication failed")
}
