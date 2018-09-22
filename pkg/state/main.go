package state

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// State is the interface to the state object
type State struct {
	DB *sql.DB
}

// User is a user object
type User struct {
	ID       int
	Username string
	Data     UserData
}

// UserData is a key value store for embedding data
// that can be checked by an application later
type UserData map[string]interface{}

// NewClient creates a new client for accessing the state.
func NewClient(path string) (*State, error) {
	isInit := false
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Print("initializing state storage ...")
		isInit = true
	}
	db, err := sql.Open("sqlite3", path)

	if isInit {
		initStmt := `create table users (id integer not null primary key autoincrement, username text);`

		_, err = db.Exec(initStmt)
		if err != nil {
			log.Printf("Failed to initialize SQLite database: %q: %s\n", err, initStmt)
			return nil, err
		}
	}

	var state State
	state.DB = db
	return &state, err
}

// GetUserByUsername finds a user by their username
func (s *State) GetUserByUsername(username string) (*User, error) {
	stmt, err := s.DB.Prepare("select id, username from users where username = ?")
	if err != nil {
		return nil, err
	}

	var user User
	err = stmt.QueryRow(username).Scan(&user.ID, &user.Username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ListUsers returns a list of ALLLLL the users.
func (s *State) ListUsers() ([]User, error) {
	rows, err := s.DB.Query("select id, username from users")
	if err != nil {
		return nil, err
	}

	// we'll resize this later :(
	userSlice := make([]User, 0)

	defer rows.Close()
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Username)

		if err != nil {
			return nil, nil
		}

		userSlice = append(userSlice, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return userSlice, nil
}

// CreateUser creates a user.
func (s *State) CreateUser(username string) (*User, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare("insert into users(id, username) values(null, ?)")
	if err != nil {
		return nil, err
	}

	log.Print("inserted")

	defer stmt.Close()

	_, err = stmt.Exec(username)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	user, err := s.GetUserByUsername(username)
	return user, err
}

// IsNotFound checks if the error type is a "not found" error type.
func (s *State) IsNotFound(err error) bool {
	if err == sql.ErrNoRows {
		return true
	}

	return false
}
