package services

import (
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/kylerequez/marketify/src/db"
	"github.com/kylerequez/marketify/src/models"
	"github.com/kylerequez/marketify/src/repositories"
	"github.com/kylerequez/marketify/src/shared"
	"github.com/kylerequez/marketify/src/utils"
	"github.com/kylerequez/marketify/src/views/components"
	"github.com/kylerequez/marketify/src/views/pages"
)

type UserService struct {
	Ur    *repositories.UserRepository
	Store *db.PostgresStorage
}

func NewUserService(ur *repositories.UserRepository, store *db.PostgresStorage) *UserService {
	return &UserService{
		Ur:    ur,
		Store: store,
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
		return utils.Render(c, components.LoginForm(form))
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

	session, err := us.Store.CreateNewSession(*user)
	if err != nil {
		form.Errors["form"] = err.Error()
		return utils.Render(c, components.LoginForm(form))
	}

	c.Cookie(&fiber.Cookie{
		Name:    "marketify-user-session",
		Value:   user.ID.String(),
		Expires: session.Expiration,
	})

	return utils.Redirect(c, http.StatusOK, "/dashboard/users")
}

func (us *UserService) GetSignupPage(c fiber.Ctx) error {
	info := shared.PageInfo{
		Title: "Sign Up",
		Path:  c.Path(),
	}

	form := shared.SignupFormData{}

	return utils.Render(c, pages.Signup(info, form))
}

func (us *UserService) CreateUser(c fiber.Ctx) error {
	form := shared.SignupFormData{
		Errors: make(map[string]string),
	}

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
		form.Errors["form"] = err.Error()
		return utils.Render(c, components.SignupForm(form))
	}

	form.Firstname = body.Firstname
	form.Middlename = body.Middlename
	form.Lastname = body.Lastname
	form.Age = body.Age
	form.Gender = body.Gender
	form.Email = body.Email
	form.Password = body.Password
	form.RePassword = body.RePassword

	hasError := false
	if err := utils.ValidateName(body.Firstname, "firstname"); err != nil {
		hasError = true
		form.Errors["firstname"] = err.Error()
	}

	if err := utils.ValidateName(body.Lastname, "lastname"); err != nil {
		hasError = true
		form.Errors["lastname"] = err.Error()
	}

	if err := utils.ValidateAge(body.Age); err != nil {
		hasError = true
		form.Errors["age"] = err.Error()
	}

	if err := utils.ValidateGender(body.Gender); err != nil {
		hasError = true
		form.Errors["gender"] = err.Error()
	}

	if err := utils.ValidateEmail(body.Email); err != nil {
		hasError = true
		form.Errors["email"] = err.Error()
	}

	if err := utils.ValidatePassword(body.Password); err != nil {
		hasError = true
		form.Errors["password"] = err.Error()
	}

	if err := utils.ValidatePassword(body.RePassword); err != nil {
		hasError = true
		form.Errors["rePassword"] = err.Error()
	}

	if !strings.EqualFold(body.Password, body.RePassword) {
		hasError = true
		form.Errors["rePassword"] = "passwords are not the same"
	}

	hashedPassword, err := utils.EncryptPassword(body.Password, bcrypt.DefaultCost)
	if err != nil {
		hasError = true
		form.Errors["password"] = err.Error()
	}

	if hasError {
		return utils.Render(c, components.SignupForm(form))
	}

	user := models.User{
		Firstname:   body.Firstname,
		Middlename:  body.Middlename,
		Lastname:    body.Lastname,
		Age:         body.Age,
		Gender:      body.Gender,
		Email:       body.Email,
		Password:    hashedPassword,
		Authorities: []string{shared.ROLES["ADMIN"]},
		Status:      shared.STATUS["UNVERIFIED"],
	}

	if err := us.Ur.CreateUser(c.Context(), user); err != nil {
		form.Errors["form"] = err.Error()
		return utils.Render(c, components.SignupForm(form))
	}

	return c.SendStatus(http.StatusOK)
}

func (us *UserService) LogoutUser(c fiber.Ctx) error {
	cookie := c.Cookies("marketify-user-session", "")
	if cookie == "" {
		return utils.Redirect(c, http.StatusBadRequest, "/login")
	}
	c.Cookie(&fiber.Cookie{
		Name:    "marketify-user-session",
		Expires: time.Now().Add(-(time.Hour * 2)),
	})

	id, err := uuid.Parse(cookie)
	if err != nil {
		return utils.Redirect(c, http.StatusBadRequest, "/login")
	}

	isAlive, err := us.Store.IsSessionAlive(id)
	if err != nil {
		return utils.Redirect(c, http.StatusBadRequest, "/login")
	}

	if !isAlive {
		return utils.Redirect(c, http.StatusBadRequest, "/login")
	}

	if err := us.Store.RemoveSession(id); err != nil {
		return utils.Redirect(c, http.StatusBadRequest, "/login")
	}

	return utils.Redirect(c, http.StatusOK, "/login")
}

func (us *UserService) GetUsersPage(c fiber.Ctx) error {
	info := shared.PageInfo{
		Title:        "Users",
		Path:         c.Path(),
		LoggedInUser: utils.RetrieveLoggedInUser(c, us.Ur),
	}

	var users []models.User = []models.User{}

	return utils.Render(c, pages.Users(info, users))
}
