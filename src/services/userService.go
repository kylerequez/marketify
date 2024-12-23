package services

import (
	"log"
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

	if user.Status != shared.STATUS["ACTIVE"] {
		form.Errors["form"] = "user is currently " + user.Status
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
		Value:   session.ID,
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
		Birthdate  string
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
	form.Birthdate = body.Birthdate
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

	birthdate, err := time.Parse("2006-01-02", body.Birthdate)
	if err != nil {
		hasError = true
		form.Errors["birthdate"] = err.Error()
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
		Birthdate:   birthdate,
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
		LoggedInUser: utils.RetrieveLoggedInUser(c),
	}

	users, err := us.Ur.GetAllUsers(c.Context())
	if err != nil {
		return utils.Render(c, pages.Users(info, *users))
	}

	return utils.Render(c, pages.Users(info, *users))
}

func (us *UserService) GetUserPage(c fiber.Ctx) error {
	info := shared.PageInfo{
		Title:        "User",
		Path:         c.Path(),
		LoggedInUser: utils.RetrieveLoggedInUser(c),
	}

	id := c.Params("id")
	if err := utils.ValidateId(id); err != nil {
		return c.Status(http.StatusBadRequest).Redirect().Back()
	}

	userId, err := uuid.Parse(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).Redirect().Back()
	}

	user, err := us.Ur.GetUserById(c.Context(), userId)
	if err != nil {
		return c.Status(http.StatusBadRequest).Redirect().Back()
	}

	return utils.Render(c, pages.User(info, *user))
}

func (us *UserService) GetUserEditForm(c fiber.Ctx) error {
	id := c.Params("id")
	if err := utils.ValidateId(id); err != nil {
		log.Println(err)
		return c.Redirect().Back()
	}

	userId, err := uuid.Parse(id)
	if err != nil {
		log.Println(err)
		return c.Redirect().Back()
	}

	user, err := us.Ur.GetUserById(c.Context(), userId)
	if err != nil {
		log.Println(err)
		return c.Redirect().Back()
	}

	form := shared.EditUserFormData{
		ID:         user.ID,
		Firstname:  user.Firstname,
		Middlename: user.Middlename,
		Lastname:   user.Lastname,
		Gender:     user.Gender,
		Email:      user.Email,
		Errors:     make(map[string]string),
	}

	return utils.Render(c, components.UserEditForm(form))
}

func (us *UserService) GetUserDeleteForm(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Redirect().Back()
	}

	userId, err := uuid.Parse(id)
	if err != nil {
		return c.Redirect().Back()
	}

	return utils.Render(c, components.UserDeleteForm(userId))
}

func (us *UserService) UpdateUser(c fiber.Ctx) error {
	form := shared.EditUserFormData{
		Errors: make(map[string]string),
	}

	id := c.Params("id")
	if id == "" {
		return c.Redirect().Back()
	}

	userId, err := uuid.Parse(id)
	if err != nil {
		form.Errors["form"] = err.Error()
		return utils.Render(c, components.UserEditForm(form))
	}

	type RequestBody struct {
		Firstname  string
		Middlename string
		Lastname   string
		Gender     string
		Email      string
	}

	body := new(RequestBody)
	if err := c.Bind().Body(body); err != nil {
		form.Errors["form"] = err.Error()
		return utils.Render(c, components.UserEditForm(form))
	}

	form.Firstname = body.Firstname
	form.Middlename = body.Middlename
	form.Lastname = body.Lastname
	form.Gender = body.Gender
	form.Email = body.Email

	hasErrors := false
	if err := utils.ValidateName(body.Firstname, "firstname"); err != nil {
		hasErrors = true
		form.Errors["firstname"] = err.Error()
	}

	if err := utils.ValidateName(body.Lastname, "lastname"); err != nil {
		hasErrors = true
		form.Errors["lastname"] = err.Error()
	}

	if err := utils.ValidateGender(body.Gender); err != nil {
		hasErrors = true
		form.Errors["gender"] = err.Error()
	}

	if err := utils.ValidateEmail(body.Email); err != nil {
		hasErrors = true
		form.Errors["email"] = err.Error()
	}

	if hasErrors {
		return utils.Render(c, components.UserEditForm(form))
	}

	user, err := us.Ur.GetUserById(c.Context(), userId)
	if err != nil {
		form.Errors["form"] = err.Error()
		return utils.Render(c, components.UserEditForm(form))
	}

	isExists, err := us.Ur.GetUserByEmail(c.Context(), body.Email)
	if err != nil {
		form.Errors["form"] = err.Error()
		return utils.Render(c, components.UserEditForm(form))
	}

	if isExists != nil && user.ID != isExists.ID {
		form.Errors["email"] = "email already exists"
		return utils.Render(c, components.UserEditForm(form))
	}

	updatedUser := models.User{
		ID:         userId,
		Firstname:  body.Firstname,
		Middlename: body.Middlename,
		Lastname:   body.Lastname,
		Gender:     body.Gender,
		Email:      body.Email,
	}

	if err := us.Ur.UpdateUser(c.Context(), updatedUser); err != nil {
		form.Errors["form"] = err.Error()
		return utils.Render(c, components.UserEditForm(form))
	}

	return utils.Redirect(c, 200, "/dashboard/users/"+userId.String())
}

func (us *UserService) DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Redirect().Back()
	}

	userId, err := uuid.Parse(id)
	if err != nil {
		return c.Redirect().Back()
	}

	if err := us.Ur.DeleteUserById(c.Context(), userId); err != nil {
		return nil
	}

	return utils.Render(c, components.UserDeleteForm(userId))
}
