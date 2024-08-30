Go Application for Authentication and Access Control System

Description This Go program is designed to act as an authentication and access control system.

Features Authentication: The auth package handles user authentication, including verifying login credentials and managing user sessions. It retrieves user information such as usernames and hashed passwords from files (like salt.txt and shadow.txt) and ensures that users are authenticated correctly before granting access to the system.

Hashing: The hash package is responsible for generating secure password hashes using the MD5 algorithm. It creates a unique salt for each user, combines it with the user's password, and generates a hash that is stored securely. This ensures that passwords are not stored in plain text and provides additional security through the use of salted hashes.

Models: The models package contains the core data structures and functions for managing the file system within the application. It includes operations for creating, reading, writing, appending, and listing files, with access control based on the user's security clearance level. The package also manages in-memory storage of files during a session and handles the saving of these files back to disk when necessary.

Installation

Steps

Download the zip file containing the binary executive file named main.

Navigate to the directory containing the source code files.

To compile the code and build an executable file, navigate to the directory containing the source code files and run the command: go build main.go. It is required that the user has go version 1.22.2. Go can be installed by running the following command: sudo apt install golang-go.

NOTE TO TUTOR: CAPA may have blocked permissions to download and install go. For testing and marking, the executable file has been built on local system and located in the same directory as the source code, named main. The marker is able to proceed with steps 4-5 for testing of the authentication system.

Input the command ./main -i to run the filesystem in sign-up mode.

Input the command ./main to run the filesystem in login mode.

Description of the Secure File System

User Creation During initialization (./main -i), the program prompts for a username, password, and security clearance level. The password must be at least 8 characters long. If the user does not enter a valid password, the program will display an error message and terminate.

Once the input is validated, the program generates an 8-digit salt using the generateSalt() function. The salt is concatenated with the user's password and hashed using the MD5Hash function from the hash package. The resulting hash, along with the user's security clearance, is stored in shadow.txt. The username and salt value are stored in salt.txt.

The program terminates after the user creation process is completed.

User Login During a standard run (./main), the program prompts for the username and password. The system retrieves the salt associated with the username from salt.txt and uses it to generate the MD5 hash of the entered password. This generated hash is then compared with the stored hash in shadow.txt. If the values match, the user is authenticated and logged in.

If the user is not authenticated, the program terminates with an appropriate error message.

Main Menu Initialization After a successful login, the Start() function is called, which initializes the main menu. The program tests the MD5 hashing function by calling MD5Hash from the hash package with a test string.

The program then calls the LoadFileSystem() function. This function checks if a Files.store exists. If the file does not exist, a new Files.store is created. If there is an error creating the file, an error message is displayed.

The Files.store file contains information in the following format: [filename]:[username]:[clearance]:[content]

If the file exists, it extracts each file's data, creates a File struct, and stores it in the in-memory file system.

Main Menu Functions

"C" - Create a File: The program calls handleCreate(). The user is prompted for a filename. If the file does not exist in the in-memory file system, it is created with the current user's username as the owner and the user's clearance level as the file’s classification. The created file's classification level matches the user's clearance level.

"A" - Append to a File: The program calls handleAppend(). The user is prompted for a filename. If the file exists, it checks if the user’s clearance level allows them to append to the file (user clearance level ≥ file classification level). If allowed, the user is prompted to enter content, which is appended to the file.

"R" - Read a File: The program calls handleRead(). The user is prompted for a filename. If the file exists, it checks if the user’s clearance level allows them to read the file (user clearance level ≤ file classification level). If allowed, the file's content is displayed.

"W" - Write to a File: The program calls handleWrite(). The user is prompted for a filename. If the file exists, it checks if the user’s clearance level matches the file’s classification level (user clearance level == file classification level). If allowed, the user is prompted to enter new content, which replaces the existing content.

"L" - List All Files: The program calls listFiles(). This displays the names of all files that the user has clearance to access. If no files are present, it informs the user that no accessible files are available.

"S" - Save Files to Disk: The program calls saveFileSystem(). This writes the in-memory file system back to Files.store, saving any changes made during the session.

"E" - Exit the System: The program calls handleExit(). The user is prompted to confirm if they want to exit. If confirmed, the program exits; otherwise, it returns to the main menu.

Security Model Implementation The Secure File System follows the Bell-LaPadula model for access control:

Read Access: Users can only read files if their clearance level is greater than or equal to the file's classification level (user clearance level >= file classification level).

Write Access: Users can only write to files if their clearance level is exactly equal to the file's classification level (user clearance level == file classification level).

Append Access: Users can append to files if their clearance level is less than or equal to the file's classification level (user clearance level <= file classification level).

The file system enforces these rules to ensure that users cannot access or modify files beyond their authorized clearance level.
