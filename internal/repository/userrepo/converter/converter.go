package converter
import (
	"github.com/KulaginNikita/pvz-service/internal/domain/user"
	"github.com/KulaginNikita/pvz-service/internal/repository/userrepo/model"
)

func ToDB(u *user.User) *model.User {
	return &model.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
		Role:     string(u.Role),
	}
}

func FromDB(m *model.User) *user.User {
	return &user.User{
		ID:       m.ID,
		Email:    m.Email,
		Password: m.Password,
		Role:     user.Role(m.Role),
	}
}
