package services

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/kylerequez/marketify/src/repositories"
	"github.com/kylerequez/marketify/src/utils"
)

type UserService struct {
	Ur *repositories.UserRepository
}

func NewUserService(ur *repositories.UserRepository) *UserService {
	return &UserService{
		Ur: ur,
	}
}

func (us *UserService) LoginUser(c fiber.Ctx) error {
	type RequestBody struct {
		Email    string
		Password string
	}

	body := new(RequestBody)
	if err := c.Bind().Body(body); err != nil {
		return err
	}

	hasError := false
	errors := make(map[string]string)
	if err := utils.ValidateEmail(body.Email); err != nil {
		hasError = true
		errors["email"] = err.Error()
	}

	if err := utils.ValidatePassword(body.Password); err != nil {
		hasError = true
		errors["password"] = err.Error()
	}

	if hasError {
		return c.JSON(
			fiber.Map{
				"success": "false",
				"errors":  errors,
			},
		)
	}

	user, err := us.Ur.GetUserByEmail(c.Context(), body.Email)
	if err != nil {
		return c.JSON(
			fiber.Map{
				"success": "false",
				"error":   err.Error(),
			},
		)
	}

	if err := utils.ComparePassword(user.Password, []byte(body.Password)); err != nil {
		return c.JSON(
			fiber.Map{
				"success": "false",
				"error":   err.Error(),
			},
		)
	}

	return c.SendStatus(http.StatusOK)
}

func (us *UserService) CreateUser(c fiber.Ctx) error {
	type RequestBody struct {
		Firstname  string
		Middlename string
		Lastname   string
		Age        uint
		Gender     string
		Email      string
		Password   string
		RePassword string
	}

	body := new(RequestBody)
	if err := c.Bind().Body(body); err != nil {
		return err
	}

	hasError := false
	errs := make(map[string]string)
	if err := utils.ValidateName(body.Firstname, "firstname"); err != nil {
		hasError = true
		errs["firstname"] = err.Error()
	}

	if err := utils.ValidateName(body.Lastname, "lastname"); err != nil {
		hasError = true
		errs["lastname"] = err.Error()
	}

	if err := utils.ValidateEmail(body.Email); err != nil {
		hasError = true
		errs["email"] = err.Error()
	}

	if err := utils.ValidatePassword(body.Password); err != nil {
		hasError = true
		errs["password"] = err.Error()
	}

	if err := utils.ValidatePassword(body.RePassword); err != nil {
		hasError = true
		errs["rePassword"] = err.Error()
	}

	if !strings.EqualFold(body.Password, body.RePassword) {
		hasError = true
		errs["rePassword"] = "passwords are not the same"
	}

	if hasError {
		return c.JSON(fiber.Map{
			"success": false,
			"errors":  errs,
		})
	}

	return c.SendStatus(http.StatusOK)
}
