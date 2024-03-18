package database

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Email string `json:"email"`
	Id    int    `json:"id"`
	Token string `json:"token,omitempty"`
}

func (u *User) ToUserResponse(token string) UserResponse {
	return UserResponse{
		Id:    u.Id,
		Email: u.Email,
		Token: token,
	}
}

// CreateUsers creates a new user and saves it to disk
func (db *DB) CreateUsers(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	allUsers, err := db.GetAllUsers()
	if err != nil {
		return User{}, err
	}

	for _, user := range allUsers {
		if user.Email == email {
			return User{}, ErrAlreadyExists
		}
	}

	id := len(dbStructure.Users) + 1
	cryptPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	newUser := User{
		Id:       id,
		Email:    email,
		Password: string(cryptPass),
	}
	dbStructure.Users[id] = newUser

	err = db.WriteDB(dbStructure)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}

// GetUsers returns all users in the database
func (db *DB) GetUser(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}
	return user, nil
}

func (db *DB) GetAllUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	var usersArray []User
	for _, user := range dbStructure.Users {
		usersArray = append(usersArray, user)
	}

	return usersArray, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, ErrNotExist
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (db *DB) UpdateUser(id int, email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}

	cryptPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	user.Email = email
	user.Password = string(cryptPass)
	dbStructure.Users[id] = user

	err = db.WriteDB(dbStructure)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
