package users

import (
	"time"
)

type Users struct {
	ID 				 int
	Name 			 string
	Email 		 string
	Password 	 string
	Created_At time.Time
	Updated_At time.Time
	Token 		 string
}