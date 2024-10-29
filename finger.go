package finger

import (
	"fmt"
	"strings"
	"sync"
)

type User struct {
	Username  string
	FullName  string
	LastLogin string
	IdleTime  string
	Plan      string
}

func (u User) String() string {
	return fmt.Sprintf("[%s, %s, %s, %s, %s]", u.Username, u.FullName, u.LastLogin, u.IdleTime, u.Plan)
}

func (u User) PartialString() string {
	return fmt.Sprintf("[%s, %s]", u.Username, u.FullName)
}

type UserDatabase struct {
	users []User
	mu    sync.Mutex
}

func NewUserDatabase() *UserDatabase {
	database := &UserDatabase{users: make([]User, 0)}
	database.users = append(database.users, User{Username: "touqan", FullName: "Zahi Touqan", LastLogin: "Sun", IdleTime: "1d"})
	database.users = append(database.users, User{Username: "mustermann", FullName: "Max Mustermann", LastLogin: "Oct 2", IdleTime: "26d"})
	database.users = append(database.users, User{Username: "touerika", FullName: "Erika Mustermann", LastLogin: "Mon", IdleTime: "", Plan: "Test Plan"})
	return database
}

func (db *UserDatabase) GetAllUsers(verbose bool) []User {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.users
}

func (db *UserDatabase) GetUser(username string) (User, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, user := range db.users {
		if user.Username == username {
			return user, true
		}
	}
	return User{}, false
}

func (db *UserDatabase) GetUserAmbiguous(username string) []User {
	db.mu.Lock()
	defer db.mu.Unlock()

	var usersList []User
	for _, user := range db.users {
		if strings.Contains(user.Username, username) {
			usersList = append(usersList, user)
		}
	}
	return usersList
}
