package shared

import "github.com/kylerequez/marketify/src/models"

type PageInfo struct {
	Title        string
	Path         string
	LoggedInUser *models.User
}
