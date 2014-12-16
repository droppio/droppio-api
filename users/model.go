package users

import (
	"code.google.com/p/go.crypto/bcrypt"
	"dropler-new/store"
	"log"
	"strings"
	"time"
)

// TimeStamp struct for the Model Interface
type TimeStamp struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User model struct
type User struct {
	ID             int64  `json:"id"db:"Id"`
	Name           string `form:"name"json:"name"`
	Email          string `form:"email"json:"email"`
	HashedPassword string `json:"hashed_password"`
	TimeStamp
}

type UserList []User

func (u *UserList) List() error {
	_, err := store.Db.Select(u, "SELECT * FROM USERS ORDER BY CreatedAt DESC")
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// Insert Method to create a new user from the models User struct
func (u *User) Insert(password string) error {

	// Lowercase email
	u.Email = strings.ToLower(u.Email)

	// run the SetPassword method on the user model
	// if a password is provided
	if password != "" {
		err := u.SetPassword(password)
		if err != nil {
			return err
		}
	}

	// run the UpdateTime ethod on the user model
	u.UpdateTime()

	// run the DB insert function
	err := store.Db.Insert(u)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetById(id string) error {
	err := store.Db.SelectOne(u, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateTime Method for setting/updating the time
// struct elements
func (t *TimeStamp) UpdateTime() {

	currentTime := time.Now().UTC()
	if !t.CreatedAt.IsZero() {
		t.UpdatedAt = currentTime
		return
	}
	t.CreatedAt = currentTime
	t.UpdatedAt = currentTime
	return
}

// SetPassword Method on User model for setting the hashed password
func (u *User) SetPassword(password string) error {

	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.HashedPassword = string(b)

	return nil
}

// CheckPassword Method to check if password matches the stored hash
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
}
