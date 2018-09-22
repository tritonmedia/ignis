package state

import (
	"database/sql"
	"log"
	"os"
	"time"

	// used for sqlite driver
	_ "github.com/mattn/go-sqlite3"
	cache "github.com/patrickmn/go-cache"
)

// State is the interface to the state object
type State struct {
	DB *sql.DB
}

// User is a user object
type User struct {
	ID       int
	Username string
	Stage    string // stage the user is at "context"
}

// UserData is a key value store for embedding data
// that can be checked by an application later
type UserData map[string]interface{}

// stageUserLookup tracks pointers to "instances" of the cache objects
// that are used for each user's stage data.
var stageUserLookup map[int]map[string]*cache.Cache

// NewClient creates a new client for accessing the state.
func NewClient(path string) (*State, error) {
	isInit := false
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Print("initializing state storage ...")
		isInit = true
	}
	db, err := sql.Open("sqlite3", path)

	if isInit {
		initStmt := `create table users (id integer not null primary key, username text, stage text);`

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

// GetUserByUsername gets a user by their Username. DEPRECATED.
func (s *State) GetUserByUsername(username string) (*User, error) {
	stmt, err := s.DB.Prepare("select id, username, stage from users where username = ?")
	if err != nil {
		return nil, err
	}

	var user User
	err = stmt.QueryRow(username).Scan(&user.ID, &user.Username, &user.Stage)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID gets a user by their ID
func (s *State) GetUserByID(id int) (*User, error) {
	stmt, err := s.DB.Prepare("select id, username, stage from users where id = ?")
	if err != nil {
		return nil, err
	}

	var user User
	err = stmt.QueryRow(id).Scan(&user.ID, &user.Username, &user.Stage)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ListUsers returns a list of ALLLLL the users.
func (s *State) ListUsers() ([]User, error) {
	rows, err := s.DB.Query("select id, username, stage from users")
	if err != nil {
		return nil, err
	}

	// we'll resize this later :(
	userSlice := make([]User, 0)

	defer rows.Close()
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Username, &user.Stage)

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
func (s *State) CreateUser(id int, username string) (*User, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare("insert into users(id, username, stage) values(?, ?, \"init\")")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id, username)
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

// NewStageStorage creates a state storage interface that supports creating a new one,
// or using an already existing one
func NewStageStorage(u *User, stage string) *cache.Cache {
	if len(stageUserLookup) == 0 {
		log.Println("[store/stage] DEBU: initialized the first stage memory object")

		// initialize the lookup table, tracks instances of the cache object
		stageUserLookup = make(map[int]map[string]*cache.Cache)
	}

	// if not set, create the user table
	if _, ok := stageUserLookup[u.ID]; !ok {
		log.Printf("[store/stage] DEBU: initalized a fresh in-mem stage map for user: %s (uid: %d)", u.Username, u.ID)
		stageUserLookup[u.ID] = make(map[string]*cache.Cache)
	}

	// check if we have a cache object already created, get it's pointer
	if c, ok := stageUserLookup[u.ID][stage]; ok {
		if c != nil {
			log.Printf("[store/stage] DEBU: using exisiting stage cache object for user: %s (uid: %d)", u.Username, u.ID)

			return c
		}
	}

	// if we're here, we don't have one, create a new cache object and return it, storing it for later.
	log.Printf("[store/stage] DEBU: creating new cache object for user: %s (uid: %d)", u.Username, u.ID)
	stageUserLookup[u.ID][stage] = cache.New(5*time.Minute, 10*time.Minute)
	return stageUserLookup[u.ID][stage]
}
