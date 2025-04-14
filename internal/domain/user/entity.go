package user

type Role string

const (
	RoleEmployee  Role = "employee"
	RoleModerator Role = "moderator"
)

type User struct {
	ID       int64
	Email    string
	Password string
	Role     Role
}
