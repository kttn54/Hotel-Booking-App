package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/term"
)

var wg = sync.WaitGroup{}
var loggedInCustomerFirstName string
var loggedInCustomerLastName string
var loggedInCustomerEmail string
var enteredEmail string
var enteredPassword string
var userChoice int
var db *sql.DB

func main() {
	wg.Add(1)
	db = connectToDatabase()

	displayIntroMessage()

	wg.Wait()
	db.Close()
}

func insertCustomerRecordIntoDatabase(db *sql.DB, firstName string, lastName string, email string, password string) {
	query := "INSERT INTO customers (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)"

	// Execute the SQL statement with user registration data
	_, err := db.Exec(query, firstName, lastName, email, password)
	if err != nil {
		fmt.Println("Error is 1")
		log.Fatal(err)
	} else {
		fmt.Print("\nUser registration is successful.\n")
	}

	displayIntroMessage()
}

func displayLoginScreen(db *sql.DB) {
	fmt.Print("\nEnter your email: ")
	fmt.Scan(&enteredEmail)
	fmt.Print("Enter your password: ")
	enteredPassword, enteredPasswordErr := term.ReadPassword(int(os.Stdin.Fd()))
	if enteredPasswordErr != nil {
		log.Fatal(enteredPasswordErr)
	}

	query := "SELECT password FROM customers WHERE email=$1"

	// Execute the SQL statement with user login data
	row := db.QueryRow(query, enteredEmail)

	var password string
	err := row.Scan(&password)
	if err == sql.ErrNoRows {
		fmt.Println("User not found")
		displayLoginScreen(db)
	} else if err != nil {
		fmt.Println("Error occurred while fetching the password.")
		log.Fatal(err)
		displayLoginScreen(db)
	} else {
		if password == string(enteredPassword) {
			fmt.Print("\nUser login is successful.\n")
			getCustomerDetails()
			displayHomeScreen(db)
		} else {
			fmt.Print("\nIncorrect email/password.\n")
			displayLoginScreen(db)
		}
	}
}

func getCustomerDetails() {
	firstNameQuery := "SELECT first_name FROM customers WHERE email=$1"

	firstNameRow := db.QueryRow(firstNameQuery, enteredEmail)
	firstNameError := firstNameRow.Scan(&loggedInCustomerFirstName)
	if firstNameError != nil {
		log.Fatal(firstNameError)
	}

	lastNameQuery := "SELECT last_name FROM customers WHERE email=$1"

	lastNameRow := db.QueryRow(lastNameQuery, enteredEmail)
	lastNameError := lastNameRow.Scan(&loggedInCustomerLastName)
	if lastNameError != nil {
		log.Fatal(lastNameError)
	}
}

func displayRegisterScreen(db *sql.DB) {
	var registerCustomerFirstName string
	var registerCustomerLastName string
	var registerCustomerEmail string

	fmt.Print("\nWhat is your first name? ")
	fmt.Scan(&registerCustomerFirstName)
	fmt.Print("What is your last name? ")
	fmt.Scan(&registerCustomerLastName)
	fmt.Print("What is your email? ")
	fmt.Scan(&registerCustomerEmail)
	fmt.Print("What is your password? ")
	customerPassword, customerPasswordErr := term.ReadPassword(int(os.Stdin.Fd()))
	if customerPasswordErr != nil {
		log.Fatal(customerPasswordErr)
	}
	fmt.Print("Re-enter your password: ")
	reEnteredCustomerPassword, reEnteredCustomerPasswordErr := term.ReadPassword(int(os.Stdin.Fd()))
	if reEnteredCustomerPasswordErr != nil {
		log.Fatal(reEnteredCustomerPasswordErr)
	}
	if string(customerPassword) != string(reEnteredCustomerPassword) {
		fmt.Print("The passwords do not match. Try again.\n")
		displayRegisterScreen(db)
	}

	insertCustomerRecordIntoDatabase(db, registerCustomerFirstName, registerCustomerLastName, registerCustomerEmail, string(customerPassword))
}

func displayIntroMessage() {
	fmt.Print("\n##################################\n")
	fmt.Print("Welcome to the Hotel Booking App!\n")
	fmt.Print("##################################\n\n")

	fmt.Print("    1. Register\n")
	fmt.Print("    2. Login\n")
	fmt.Print("    3. Exit\n\n")
	fmt.Print("What would you like to do? ")

	navigateUserOptionIntroScreen()
}

func displayHomeScreen(db *sql.DB) {
	fmt.Print("\n##################################\n")
	fmt.Printf("Welcome %v to the Hotel Booking App!\n", loggedInCustomerFirstName)
	fmt.Print("##################################\n\n")

	fmt.Print("    1. Book a room\n")
	fmt.Print("    2. See my booked rooms\n")
	fmt.Print("    3. Exit\n\n")
	fmt.Print("What would you like to do? ")

	navigateUserOptionHomeScreen()
}

func navigateUserOptionIntroScreen() {
	fmt.Scan(&userChoice)

	switch userChoice {
	case 1:
		displayRegisterScreen(db)
	case 2:
		displayLoginScreen(db)
	case 3:
		os.Exit(0)
	}
}

func navigateUserOptionHomeScreen() {
	fmt.Scan(&userChoice)

	switch userChoice {
	case 1:
		getCustomerDetails()
		//displayUserRoomFilters(db)
	case 2:
		//displayUserRooms(db)
	case 3:
		os.Exit(0)
	}
}

func connectToDatabase() *sql.DB {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	username := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	host := os.Getenv("PGHOST")
	port := os.Getenv("PGPORT")
	dbname := os.Getenv("PGDATABASE")

	// Database connection string
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", username, password, host, port, dbname)

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database")
	wg.Done()

	return db
}
