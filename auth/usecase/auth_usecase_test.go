package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/giovanisilqueirasantos/e-commerce-go-clean-arch/domain"
	"github.com/giovanisilqueirasantos/e-commerce-go-clean-arch/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginCheckLoginExistsError(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)

	var mockAuth domain.Auth
	mockAuth.Login = "valid login"

	mockAuthRepo.On("GetByLogin", mock.Anything, mockAuth.Login).Return(nil, errors.New("error message"))

	authUseCase := NewAuthUseCase(nil, nil, mockAuthRepo, nil)

	_, errToken := authUseCase.Login(context.Background(), &mockAuth)

	assert.Error(t, errToken)
}

func TestLoginCheckLoginExists(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)

	var mockAuth domain.Auth
	mockAuth.Login = "valid login"

	mockAuthRepo.On("GetByLogin", mock.Anything, mockAuth.Login).Return(nil, nil)

	authUseCase := NewAuthUseCase(nil, nil, mockAuthRepo, nil)

	_, errToken := authUseCase.Login(context.Background(), &mockAuth)

	assert.Error(t, errToken)
}

func TestLoginPassIsEqualHashedPassError(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)
	mockAuthService := new(mocks.MockAuthService)

	var mockAuth domain.Auth
	mockAuth.Login = "valid login"
	mockAuth.Password = "invalid password"

	mockAuthRepo.On("GetByLogin", mock.Anything, mockAuth.Login).Return("valid login", "valid password", nil)

	mockAuthService.On("PassIsEqualHashedPass", mock.Anything, mockAuth.Password, "valid password").Return(false)

	authUseCase := NewAuthUseCase(mockAuthService, nil, mockAuthRepo, nil)

	_, errToken := authUseCase.Login(context.Background(), &mockAuth)

	assert.Error(t, errToken)
}

func TestLoginSignTokenError(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)
	mockAuthService := new(mocks.MockAuthService)
	mockTokenService := new(mocks.MockTokenService)

	var mockAuth domain.Auth
	mockAuth.Login = "valid login"
	mockAuth.Password = "valid password"

	mockAuthRepo.On("GetByLogin", mock.Anything, mockAuth.Login).Return(mockAuth.Login, mockAuth.Password, nil)

	mockAuthService.On("PassIsEqualHashedPass", mock.Anything, mockAuth.Password, mockAuth.Password).Return(true)

	var thirtyDaysInMinutes int64 = 43200

	tokenInfo := domain.TokenInfo{Info: mockAuth.Login}

	mockTokenService.On("Sign", mock.Anything, tokenInfo, thirtyDaysInMinutes).Return("", errors.New("error message"))

	authUseCase := NewAuthUseCase(mockAuthService, mockTokenService, mockAuthRepo, nil)

	_, errToken := authUseCase.Login(context.Background(), &mockAuth)

	assert.Error(t, errToken)
}

func TestLoginSuccess(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)
	mockAuthService := new(mocks.MockAuthService)
	mockTokenService := new(mocks.MockTokenService)

	var mockAuth domain.Auth
	mockAuth.Login = "valid login"
	mockAuth.Password = "valid password"

	mockAuthRepo.On("GetByLogin", mock.Anything, mockAuth.Login).Return(mockAuth.Login, mockAuth.Password, nil)

	mockAuthService.On("PassIsEqualHashedPass", mock.Anything, mockAuth.Password, mockAuth.Password).Return(true)

	var thirtyDaysInMinutes int64 = 43200

	tokenInfo := domain.TokenInfo{Info: mockAuth.Login}

	mockTokenService.On("Sign", mock.Anything, tokenInfo, thirtyDaysInMinutes).Return("valid token", nil)

	authUseCase := NewAuthUseCase(mockAuthService, mockTokenService, mockAuthRepo, nil)

	token, errToken := authUseCase.Login(context.Background(), &mockAuth)

	assert.Nil(t, errToken)
	assert.Equal(t, token, domain.Token("valid token"))
}

func TestSignUpCheckLoginExistsError(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)

	var mockAuth domain.Auth
	mockAuth.Login = "valid login"

	mockAuthRepo.On("GetByLogin", mock.Anything, mockAuth.Login).Return(nil, errors.New("error message"))

	authUseCase := NewAuthUseCase(nil, nil, mockAuthRepo, nil)

	_, errToken := authUseCase.SignUp(context.Background(), &mockAuth, nil)

	assert.Error(t, errToken)
}

func TestSignUpLoginAlreadyExists(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)

	var mockAuth domain.Auth
	mockAuth.Login = "valid login"

	mockAuthRepo.On("GetByLogin", mock.Anything, mockAuth.Login).Return("valid login", "valid password", nil)

	authUseCase := NewAuthUseCase(nil, nil, mockAuthRepo, nil)

	_, errToken := authUseCase.SignUp(context.Background(), &mockAuth, nil)

	assert.Error(t, errToken)
}

func TestSignUpCheckUserExistsError(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)
	mockUserRepo := new(mocks.MockUserRepository)

	var mockAuth domain.Auth
	mockAuth.Login = "valid login"

	var mockUser domain.User
	mockUser.Email = "valid email"

	mockAuthRepo.On("GetByLogin", mock.Anything, mockAuth.Login).Return(nil, nil)

	mockUserRepo.On("GetByEmail", mock.Anything, mockUser.Email).Return(nil, errors.New("error message"))

	authUseCase := NewAuthUseCase(nil, nil, mockAuthRepo, mockUserRepo)

	_, errToken := authUseCase.SignUp(context.Background(), &mockAuth, &mockUser)

	assert.Error(t, errToken)
}

func TestSignUpCheckUserExists(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)
	mockUserRepo := new(mocks.MockUserRepository)

	var mockAuth domain.Auth
	mockAuth.Login = "valid login"

	var mockUser domain.User
	mockUser.Email = "valid email"

	mockAuthRepo.On("GetByLogin", mock.Anything, mockAuth.Login).Return(nil, nil)

	mockUserRepo.On("GetByEmail", mock.Anything, mockUser.Email).Return("user email", "user first name", "user last name", "user phone number", "user addres", nil)

	authUseCase := NewAuthUseCase(nil, nil, mockAuthRepo, mockUserRepo)

	_, errToken := authUseCase.SignUp(context.Background(), &mockAuth, &mockUser)

	assert.Error(t, errToken)
}

func TestSignUpStoreUserError(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthService := new(mocks.MockAuthService)

	var mockAuth domain.Auth
	mockAuth.Login = "valid login"
	mockAuth.Password = "valid password"

	var mockUser domain.User
	mockUser.Email = "user email"

	mockAuthService.On("EncodePass", mock.Anything, mockAuth.Password).Return("hashed password")

	mockUserRepo.On("GetByEmail", mock.Anything, mockUser.Email).Return(nil, nil)

	mockAuthRepo.On("GetByLogin", mock.Anything, mockAuth.Login).Return("valid login", "valid password", nil)
	mockAuthRepo.On("StoreWithUser", mock.Anything, &domain.Auth{Login: mockAuth.Login, Password: "hashed password"}, &mockUser).Return(errors.New("error message"))

	authUseCase := NewAuthUseCase(nil, nil, mockAuthRepo, mockUserRepo)

	_, errToken := authUseCase.SignUp(context.Background(), &mockAuth, &mockUser)

	assert.Error(t, errToken)
}

func TestSignUpSignTokenError(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockTokenService := new(mocks.MockTokenService)
	mockAuthService := new(mocks.MockAuthService)

	var mockAuth domain.Auth
	mockAuth.Login = "valid login"
	mockAuth.Password = "valid password"

	var mockUser domain.User
	mockUser.Email = "user email"

	mockAuthService.On("EncodePass", mock.Anything, mockAuth.Password).Return("hashed password")

	mockUserRepo.On("GetByEmail", mock.Anything, mockUser.Email).Return(nil, nil)

	mockAuthRepo.On("GetByLogin", mock.Anything, mockAuth.Login).Return("valid login", "valid password", nil)
	mockAuthRepo.On("StoreWithUser", mock.Anything, &domain.Auth{Login: mockAuth.Login, Password: "hashed password"}, &mockUser).Return(nil)

	var thirtyDaysInMinutes int64 = 43200

	tokenInfo := domain.TokenInfo{Info: mockAuth.Login}

	mockTokenService.On("Sign", mock.Anything, tokenInfo, thirtyDaysInMinutes).Return("", errors.New("error message"))

	authUseCase := NewAuthUseCase(nil, mockTokenService, mockAuthRepo, mockUserRepo)

	_, errToken := authUseCase.SignUp(context.Background(), &mockAuth, &mockUser)

	assert.Error(t, errToken)
}

func TestSignUpSuccess(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockTokenService := new(mocks.MockTokenService)
	mockAuthService := new(mocks.MockAuthService)

	var mockAuth domain.Auth
	mockAuth.Login = "valid login"
	mockAuth.Password = "valid password"

	var mockUser domain.User
	mockUser.Email = "user email"

	mockAuthService.On("EncodePass", mock.Anything, mockAuth.Password).Return("hashed password")

	mockUserRepo.On("GetByEmail", mock.Anything, mockUser.Email).Return(nil, nil)

	mockAuthRepo.On("GetByLogin", mock.Anything, mockAuth.Login).Return(nil, nil)
	mockAuthRepo.On("StoreWithUser", mock.Anything, &domain.Auth{Login: mockAuth.Login, Password: "hashed password"}, &mockUser).Return(nil)

	var thirtyDaysInMinutes int64 = 43200

	tokenInfo := domain.TokenInfo{Info: mockAuth.Login}

	mockTokenService.On("Sign", mock.Anything, tokenInfo, thirtyDaysInMinutes).Return("valid token", nil)

	authUseCase := NewAuthUseCase(mockAuthService, mockTokenService, mockAuthRepo, mockUserRepo)

	token, errToken := authUseCase.SignUp(context.Background(), &mockAuth, &mockUser)

	assert.Nil(t, errToken)
	assert.Equal(t, token, domain.Token("valid token"))
}
