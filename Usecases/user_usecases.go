package usecases

import (
	domain "task-manager/Domain"
	infrastructure "task-manager/Infrastructure"
	repositories "task-manager/Repositories"
)

type UserUsecase interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
	PromoteUser(userID string) error
}

type userUsecase struct {
	userRepo        repositories.UserRepository
	passwordService infrastructure.PasswordService
	jwtService      infrastructure.JWTService
}

func NewUserUsecase(userRepo repositories.UserRepository, passwordService infrastructure.PasswordService, jwtService infrastructure.JWTService) UserUsecase {
	return &userUsecase{
		userRepo:        userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

func (u *userUsecase) Register(username, password string) error {
	if username == "" || password == "" {
		return &domain.BadRequestError{Message: "username and password are required"}
	}

	_, err := u.userRepo.FindByUsername(username)
	if err == nil {
		return &domain.BadRequestError{Message: "username already exists"}
	} else if _, ok := err.(*domain.NotFoundError); !ok {
		return err
	}

	hashedPassword, err := u.passwordService.HashPassword(password)
	if err != nil {
		return &domain.InternalServerError{Message: "error hashing password"}
	}

	user := domain.User{
		Username: username,
		Password: hashedPassword,
		Role:     "user",
	}
	// If first user, promote to admin
	count, err := u.userRepo.CountUsers()
	if err != nil {
		return err
	}

	if count == 0 {
		user.Role = "admin"
	}

	return u.userRepo.CreateUser(user)
}

func (u *userUsecase) Login(username, password string) (string, error) {
	user, err := u.userRepo.FindByUsername(username)
	if err != nil {
		if _, ok := err.(*domain.NotFoundError); ok {
			return "", &domain.BadRequestError{Message: "invalid username or password"}
		}
		return "", &domain.InternalServerError{Message: "error authenticating user"}
	}

	if err := u.passwordService.ComparePasswords(user.Password, password); err != nil {
		return "", &domain.BadRequestError{Message: "invalid username or password"}
	}

	token, err := u.jwtService.GenerateToken(user.Username, user.Role)
	if err != nil {
		return "", &domain.InternalServerError{Message: "error generating token"}
	}

	return token, nil
}


func (u *userUsecase) PromoteUser(username string) error {
	user, err := u.userRepo.FindByUsername(username)
	if err != nil {
		return err
	}

	if user.Role == "admin" {
		return &domain.BadRequestError{Message: "user is already an admin"}
	}

	user.Role = "admin"
	return u.userRepo.UpdateUser(user.ID, user)
}
