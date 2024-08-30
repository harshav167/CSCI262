package models

import (
	"ass1/auth"
	"ass1/hash"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/rand"
)

type model struct {
	state       state
	username    string
	password    string
	clearance   int
	loggedIn    bool
	fileSystem  map[string]File
	currentFile string
	errMsg      string
	successMsg  string
}

type File struct {
	owner          string
	classification int
	content        string
}

type state int

const (
	initialState state = iota
	usernameState
	passwordState
	menuState
	createState
	appendState
	readState
	writeState
	listState
	saveState
	exitState
	appendContentState
)

// InitializeUser handles the creation of a new user
func InitializeUser() error {
	var username, password, confirmPassword string
	var clearance int

	fmt.Print("Username: ")
	fmt.Scanln(&username)

	// Check if username already exists
	if userExists(username) {
		return errors.New("username already exists")
	}

	fmt.Print("Password: ")
	fmt.Scanln(&password)

	fmt.Print("Confirm Password: ")
	fmt.Scanln(&confirmPassword)

	if password != confirmPassword {
		return errors.New("passwords do not match")
	}

	// Check password requirements
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	fmt.Print("User clearance (0, 1, 2, 3): ")
	fmt.Scanln(&clearance)

	if clearance < 0 || clearance > 3 {
		return errors.New("invalid clearance level")
	}

	// Generate salt and hash password
	salt := generateSalt()
	passSaltHash := hash.MD5Hash(password + salt)

	// Save to salt.txt and shadow.txt
	if err := saveToFile("salt.txt", fmt.Sprintf("%s:%s\n", username, salt)); err != nil {
		return err
	}
	if err := saveToFile("shadow.txt", fmt.Sprintf("%s:%s:%d\n", username, passSaltHash, clearance)); err != nil {
		return err
	}

	fmt.Println("User created successfully!")
	return nil
}

func userExists(username string) bool {
	file, err := os.ReadFile("salt.txt")
	if err != nil {
		return false
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, username+":") {
			return true
		}
	}
	return false
}

func saveToFile(filename, content string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

func generateSalt() string {
	return fmt.Sprintf("%08d", rand.Intn(100000000))
}

// InitialModel returns the initial state of the model
func InitialModel() model {
	return model{
		state:      initialState,
		fileSystem: make(map[string]File),
	}
}

func (m *model) Start() {
	m.LoadFileSystem()

	for {
		m.View()

		switch m.state {
		case initialState:
			m.state = usernameState

		case usernameState:
			m.handleTextInput(&m.username, passwordState)

		case passwordState:
			m.handlePasswordInput()

		case menuState:
			m.handleMenu()

		case createState:
			m.handleCreate()

		case appendState:
			m.handleAppend()

		case appendContentState:
			m.handleAppendContent()

		case readState:
			m.handleRead()

		case writeState:
			m.handleWrite()

		case listState:
			m.listFiles()

		case saveState:
			m.saveFileSystem()
			m.state = menuState

		case exitState:
			m.handleExit()
		}
	}
}

func (m *model) View() {
	if m.errMsg != "" {
		fmt.Println("Error:", m.errMsg)
		m.errMsg = ""
	}

	if m.successMsg != "" {
		fmt.Println("Success:", m.successMsg)
		m.successMsg = ""
	}

	switch m.state {
	case initialState:
		fmt.Println("Welcome to the Secure File System")
		fmt.Println("Press Enter to continue...")

	case usernameState:
		fmt.Print("Enter Username: ")

	case passwordState:
		fmt.Print("Enter Password: ")

	case menuState:
		fmt.Println("Options: (C)reate, (A)ppend, (R)ead, (W)rite, (L)ist, (S)ave or (E)xit.")

	case createState:
		fmt.Print("Enter Filename to Create: ")

	case appendState:
		fmt.Print("Enter Filename to Append: ")

	case appendContentState:
		fmt.Print("Enter content to append and press Enter: ")

	case readState:
		fmt.Print("Enter Filename to Read: ")

	case writeState:
		fmt.Print("Enter Filename to Write: ")

	case listState:
		fmt.Println(m.listFiles())

	case saveState:
		fmt.Println("Saving files...")

	case exitState:
		fmt.Print("Shut down the FileSystem? (Y)es or (N)o: ")

	default:
		fmt.Println("Unknown state")
	}
}

func (m *model) handleTextInput(target *string, nextState state) {
	var input string
	fmt.Scanln(&input)
	*target = strings.TrimSpace(input)
	m.state = nextState
}

func (m *model) handlePasswordInput() {
	var input string
	fmt.Scanln(&input)
	m.password = strings.TrimSpace(input)

	authSuccess, clearance, err := auth.AuthenticateUser(m.username, m.password)
	if err != nil {
		m.errMsg = "Authentication failed: " + err.Error()
		m.state = initialState
	} else if authSuccess {
		m.clearance = clearance
		m.loggedIn = true
		m.state = menuState
	} else {
		m.errMsg = "Invalid credentials"
		m.state = initialState
	}
}

func (m *model) handleMenu() {
	var input string
	fmt.Scanln(&input)
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "c":
		m.state = createState
	case "a":
		m.state = appendState
	case "r":
		m.state = readState
	case "w":
		m.state = writeState
	case "l":
		m.state = listState
		fmt.Println(m.listFiles()) // List files and return to menu state
		m.state = menuState
	case "s":
		m.state = saveState
	case "e":
		m.state = exitState
	default:
		m.errMsg = "Invalid option. Try again."
		m.state = menuState
	}
}

func (m *model) handleCreate() {
	var filename string
	fmt.Scanln(&filename)
	filename = strings.TrimSpace(filename)

	if filename == "" {
		m.errMsg = "Filename cannot be empty."
		m.state = menuState
		return
	}

	if _, exists := m.fileSystem[filename]; exists {
		m.errMsg = "File already exists."
		m.state = menuState
		return
	}

	m.fileSystem[filename] = File{owner: m.username, classification: m.clearance, content: ""}
	m.successMsg = fmt.Sprintf("File '%s' created successfully.", filename)
	m.state = menuState
}

func (m *model) handleAppend() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Filename to Append: ")
	filename, _ := reader.ReadString('\n')
	filename = strings.TrimSpace(filename)

	file, exists := m.fileSystem[filename]
	if !exists {
		m.errMsg = fmt.Sprintf("Error: File '%s' does not exist.", filename)
		m.state = menuState
		return
	}

	// Bell-LaPadula: No Write Down - user cannot append to files classified lower than their clearance level
	if file.classification < m.clearance {
		m.errMsg = "Access denied. Cannot append to files classified lower than your clearance level."
		m.state = menuState
		return
	}

	fmt.Print("Enter content to append and press Enter: ")
	content, _ := reader.ReadString('\n')
	content = strings.TrimSpace(content)

	file.content += " " + content
	m.fileSystem[filename] = file
	fmt.Printf("Content successfully appended to file '%s'.\n", filename)
	m.state = menuState
}

func (m *model) handleAppendContent() {
	reader := bufio.NewReader(os.Stdin)
	content, _ := reader.ReadString('\n')

	file := m.fileSystem[m.currentFile]
	file.content += "\n" + strings.TrimSpace(content)
	m.fileSystem[m.currentFile] = file
	m.successMsg = fmt.Sprintf("Content successfully appended to file '%s'.", m.currentFile)

	// Transition back to menu state after successful append
	m.state = menuState
}

func (m *model) handleRead() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Filename to Read: ")
	filename, _ := reader.ReadString('\n')
	filename = strings.TrimSpace(filename)

	file, exists := m.fileSystem[filename]
	if !exists {
		m.errMsg = "File does not exist."
		m.state = menuState
		return
	}

	// Bell-LaPadula: No Read Up - user cannot read files classified higher than their clearance
	if file.classification > m.clearance {
		m.errMsg = "Access denied. Clearance level too low."
		m.state = menuState
		return
	}

	// Display the file content
	fmt.Printf("File content:\n%s\n", file.content)
	m.state = menuState
}

func (m *model) handleWrite() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Filename to Write: ")
	filename, _ := reader.ReadString('\n')
	filename = strings.TrimSpace(filename)

	file, exists := m.fileSystem[filename]
	if !exists {
		m.errMsg = fmt.Sprintf("Error: File '%s' does not exist.", filename)
		m.state = menuState
		return
	}

	// Bell-LaPadula: No Write Down - user cannot write to files classified lower than their clearance
	if file.classification < m.clearance {
		m.errMsg = "Access denied. Cannot write to files classified lower than your clearance level."
		m.state = menuState
		return
	}

	fmt.Print("Enter content to write: ")
	content, _ := reader.ReadString('\n')
	content = strings.TrimSpace(content)

	file.content = content
	m.fileSystem[filename] = file
	m.successMsg = fmt.Sprintf("File '%s' content overwritten.", filename)
	m.state = menuState
}

func (m *model) listFiles() string {
	var builder strings.Builder
	for filename, file := range m.fileSystem {
		if file.classification <= m.clearance {
			builder.WriteString(fmt.Sprintf("File: %s, Owner: %s, Classification: %d\n", filename, file.owner, file.classification))
		}
	}
	if builder.Len() == 0 {
		return "No accessible files found.\n"
	}
	return builder.String()
}

func (m *model) saveFileSystem() error {
	f, err := os.Create("Files.store")
	if err != nil {
		return err
	}
	defer f.Close()

	for filename, file := range m.fileSystem {
		// Save the content as a single line in the store file
		_, err := fmt.Fprintf(f, "%s:%s:%d:%s\n", strings.TrimSpace(filename), file.owner, file.classification, strings.ReplaceAll(file.content, "\n", " "))
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *model) LoadFileSystem() error {
	file, err := os.ReadFile("Files.store")
	if err != nil {
		return err
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if len(line) > 0 {
			parts := strings.SplitN(line, ":", 4)
			if len(parts) < 4 {
				continue // Skip invalid entries
			}
			classification, err := strconv.Atoi(parts[2])
			if err != nil || classification < 0 || classification > 3 {
				continue // Skip invalid entries
			}
			m.fileSystem[strings.TrimSpace(parts[0])] = File{
				owner:          parts[1],
				classification: classification,
				content:        parts[3],
			}
		}
	}
	return nil
}

func (m *model) handleExit() {
	var input string
	fmt.Scanln(&input)
	if strings.ToLower(strings.TrimSpace(input)) == "y" {
		fmt.Println("Exiting...")
		os.Exit(0)
	} else {
		m.state = menuState
	}
}
