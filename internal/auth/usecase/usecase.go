package usecase

import (
	"chitests/internal/buildinfo"
	"chitests/internal/http/gen"
	"chitests/internal/models"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockgen -source=usecase.go -destination=../../../mocks/usecase_mock.go -package mock
type UserRepository interface {
	RegisterUser(ctx context.Context, u models.UserAccount) error
	FindUserByEmail(ctx context.Context, username string) (models.UserAccount, error)
}

//go:generate mockgen -source=usecase.go -destination=../../../mocks/usecase_mock.go -package mock
type CryptoPassword interface {
	HashPassword(password string) ([]byte, error)
	ComparePasswords(fromUser, fromDB string) bool
}

//go:generate mockgen -source=usecase.go -destination=../../../mocks/usecase_mock.go -package mock
type JWTManager interface {
	IssueToken(userID string) (string, error)
	VerifyToken(tokenString string) (*jwt.Token, error)
}

type AuthUseCase struct {
	userrepository UserRepository
	cryptopass     CryptoPassword
	jwtmanager     JWTManager
	buildINF       buildinfo.BuildInfo
}

func NewUseCase(
	ur UserRepository,
	cp CryptoPassword,
	jm JWTManager,
	bi buildinfo.BuildInfo,
) AuthUseCase {
	return AuthUseCase{
		userrepository: ur,
		cryptopass:     cp,
		jwtmanager:     jm,
		buildINF:       bi,
	}
}

func (u AuthUseCase) GetLogin(ctx context.Context, request gen.GetLoginRequestObject) (gen.GetLoginResponseObject, error) {
	user, err := u.userrepository.FindUserByEmail(ctx, request.Body.Username)

	if err != nil {
		return gen.GetLogin500JSONResponse{
			Error: err.Error(),
		}, nil
	}

	if err := validator.New().Struct(user); err != nil {
		return gen.GetLogin400JSONResponse{
			Error: err.Error(),
		}, nil
	}

	if !u.cryptopass.ComparePasswords(user.Password, request.Body.Password) {
		return gen.GetLogin401JSONResponse{Error: "unauth"}, nil
	}

	token, err := u.jwtmanager.IssueToken(user.UserName)
	if err != nil {
		return gen.GetLogin500JSONResponse{}, err
	}

	return gen.GetLogin200JSONResponse{
		AccessToken: token,
	}, nil
}

func (u AuthUseCase) PostRegister(ctx context.Context, request gen.PostRegisterRequestObject) (gen.PostRegisterResponseObject, error) {
	hashedPassword, err := u.cryptopass.HashPassword(request.Body.Password)
	if err != nil {
		return gen.PostRegister500JSONResponse{}, nil
	}

	user := models.UserAccount{
		UserName: request.Body.Username,
		Password: string(hashedPassword),
	}

	if err := validator.New().Struct(user); err != nil {
		return gen.PostRegister400JSONResponse{
			Error: err.Error(),
		}, nil
	}

	err = u.userrepository.RegisterUser(ctx, user)
	if err != nil {
		return gen.PostRegister500JSONResponse{}, nil
	}
	return gen.PostRegister201JSONResponse{
		Username: request.Body.Username,
	}, nil
}

func (u AuthUseCase) GetBuildinfo(ctx context.Context, request gen.GetBuildinfoRequestObject) (gen.GetBuildinfoResponseObject, error) {
	return gen.GetBuildinfo200JSONResponse{
		Arch:       u.buildINF.Arch,
		BuildDate:  u.buildINF.BuildDate,
		CommitHash: u.buildINF.CommitHash,
		Compiler:   u.buildINF.Compiler,
		GoVersion:  u.buildINF.GoVersion,
		Os:         u.buildINF.OS,
		Version:    u.buildINF.Version,
	}, nil
}
