package utils

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/kylerequez/marketify/src/shared"
)

func ValidateName(name, key string) error {
	if name == "" {
		return fmt.Errorf("%s must not be empty", key)
	}

	if len(key) < shared.NAME_MIN_LENGTH {
		return fmt.Errorf("%s must not be less than %d characters", key, shared.NAME_MIN_LENGTH)
	}

	if len(key) > shared.NAME_MAX_LENGTH {
		return fmt.Errorf("%s must not be greater than %d characters", key, shared.NAME_MAX_LENGTH)
	}

	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is empty")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return err
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password is empty")
	}

	if len(password) < shared.PASSWORD_MIN_LENGTH {
		return fmt.Errorf(
			"password must not be less than %d characters",
			shared.PASSWORD_MIN_LENGTH,
		)
	}

	if len(password) > shared.PASSWORD_MAX_LENGTH {
		return fmt.Errorf("password must not exceed %d characters", shared.PASSWORD_MAX_LENGTH)
	}

	return nil
}

func ValidateAge(age uint) error {
	if int(age) == 0 {
		return errors.New("age is empty")
	}

	if int(age) < shared.AGE_MIN {
		return fmt.Errorf("age must not be less than %d", shared.AGE_MIN)
	}

	if int(age) > shared.AGE_MAX {
		return fmt.Errorf("age must not be more than %d", shared.AGE_MAX)
	}

	return nil
}

func ValidateGender(gender string) error {
	if gender == "" {
		return errors.New("gender must not be empty")
	}

	_, ok := shared.GENDER[gender]
	if !ok {
		return fmt.Errorf("the gender, %s, does not exists", gender)
	}

	return nil
}
