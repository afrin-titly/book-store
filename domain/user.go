package domain

type User struct {
	Name     string
	Email    string
	Password string
	Phone    string
	Address  string
}

type UserRepository interface {
	CreateUser(user *User) (*User, error)
	GetUser(ID int) (User, error)
	ExistsByEmail(email string) (bool, error)
}
