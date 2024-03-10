package user

type UserManager interface {
	CreateUser(user User) (User, error)
	ReadUser(id int) (User, error)
	UpdateUser(toUpdateUser User) (User, error)
	DeleteUser(id int) error
}
