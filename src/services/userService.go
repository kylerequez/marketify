package services

import (
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/kylerequez/marketify/src/repositories"
	"github.com/kylerequez/marketify/src/shared"
	"github.com/kylerequez/marketify/src/utils"
	"github.com/kylerequez/marketify/src/views/components"
	"github.com/kylerequez/marketify/src/views/pages"
)

type UserService struct {
	Ur *repositories.UserRepository
}

func NewUserService(ur *repositories.UserRepository) *UserService {
	return &UserService{
		Ur: ur,
	}
}

func (us *UserService) GetLoginPage(c fiber.Ctx) error {
	info := shared.PageInfo{
		Title: "Login",
		Path:  c.Path(),
	}

	form := shared.LoginFormData{}

	return utils.Render(c, pages.Login(info, form))
}

func (us *UserService) LoginUser(c fiber.Ctx) error {
	info := shared.PageInfo{
		Title: "Login",
		Path:  c.Path(),
	}

	form := shared.LoginFormData{
		Errors: make(map[string]string),
	}

	type RequestBody struct {
		Email    string
		Password string
	}

	body := new(RequestBody)
	if err := c.Bind().Body(body); err != nil {
		form.Errors["form"] = err.Error()
		return utils.Render(c, pages.Login(info, form))
	}

	form.Email = body.Email
	form.Password = body.Password

	hasError := false
	if err := utils.ValidateEmail(body.Email); err != nil {
		hasError = true
		form.Errors["email"] = err.Error()
	}

	if err := utils.ValidatePassword(body.Password); err != nil {
		hasError = true
		form.Errors["password"] = err.Error()
	}

	if hasError {
		return utils.Render(c, components.LoginForm(form))
	}

	user, err := us.Ur.GetUserByEmail(c.Context(), body.Email)
	if err != nil {
		form.Errors["form"] = err.Error()
		return utils.Render(c, components.LoginForm(form))
	}

	if err := utils.ComparePassword(user.Password, []byte(body.Password)); err != nil {
		form.Errors["password"] = err.Error()
		return utils.Render(c, components.LoginForm(form))
	}

	log.Print("HER")
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
