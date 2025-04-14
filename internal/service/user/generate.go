package userservice

//go:generate minimock -i github.com/KulaginNikita/pvz-service/internal/repository/userrepo.UserRepository -o ./mocks -s "_mock.go"
//go:generate minimock -i github.com/KulaginNikita/pvz-service/pkg/jwtutil.TokenManager -o ./mocks -s "_mock.go" -n JwtManagerMock -p mocks
