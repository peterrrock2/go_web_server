package database

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Red      bool   `json:"is_chirpy_red"`
}

type UserResponse struct {
	Email        string `json:"email"`
	Id           int    `json:"id"`
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Red          bool   `json:"is_chirpy_red"`
}

func (u *User) ToUserResponse(access_token, refresh_token string) UserResponse {
	return UserResponse{
		Id:           u.Id,
		Email:        u.Email,
		Token:        access_token,
		RefreshToken: refresh_token,
		Red:          u.Red,
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
		Red:      false,
	}
	dbStructure.Users[id] = newUser

	err = db.WriteDB(dbStructure)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}

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

func (db *DB) UpgradeUser(id int) error {
	dbStructure, _ := db.loadDB()
	user, ok := dbStructure.Users[id]
	if !ok {
		return ErrNotExist
	}

	user.Red = true
	dbStructure.Users[id] = user
	db.WriteDB(dbStructure)
	return nil
}
