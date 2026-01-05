package users

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUserInput(input RegisterUserInput) (Users, error)
	LoginUserInput(input LoginUserInput) (Users, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	GetUserByID(ID int) (Users, error)
}

type service struct {
	repository Repository
}

func NewService(repository repository) *service {
	return &service{repository: &repository}
}

func (s *service) RegisterUserInput(input RegisterUserInput) (Users, error) {
	user := Users{}
	user.Name = input.Name
	user.Email = input.Email

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password),bcrypt.MinCost)
	if err != nil {
		return user,err
	}
	user.Password = string(passwordHash)
	newUser,err:= s.repository.Save(user)
	if err != nil {
		return user,err
	}

	return newUser,nil
}

func (s *service) LoginUserInput(input LoginUserInput) (Users, error) {
	email:= input.Email
	password:= input.Password

	user,err:= s.repository.FindByEmail(email)
	if err != nil {
		return user,err
	}

	if user.ID == 0 {
		return user,errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password))

	if err != nil{
		return user,err
	}

	return user,nil
}

func (s *service) IsEmailAvailable(input CheckEmailInput) (bool,error) {
	var user Users
	email:= input.Email
	user,err:= s.repository.FindByEmail(email)
	if err != nil {
		return false,err
	}
	if user.ID == 0 {
		return true,nil
	}

	return false,nil
}

func (s *service) GetUserByID(ID int) (Users,error)  {
	user,err:= s.repository.FindByID(ID)
	if err != nil {
		return user,err
	}
	if user.ID == 0 {
		return user,errors.New("Not user found")
	}
	return user,nil
}
