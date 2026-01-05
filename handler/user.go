package handler

import (
	"log"
	"net/http"
	"users/auth"
	"users/users"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService users.Service
	authService auth.Service
}

func NewUserHandler(userService users.Service, authService auth.Service) *userHandler{
	return &userHandler{userService: userService,authService: authService}
} 

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input users.RegisterUserInput

	err:= c.ShouldBindJSON(&input)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"status" : "error",
			"message" : err.Error(),
		})
		return
	}

	newUser,err := h.userService.RegisterUserInput(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	token,err:= h.authService.GenerateToken(
		newUser.ID,
		c.Request.Context(),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}	
	formatter:= users.FormatUser(newUser,token)

	c.JSON(http.StatusOK, gin.H{
		"status" : "success",
		"message" : "success",
		"data": formatter,
	})
}


func (h *userHandler) LoginUser(c *gin.Context) {
	var input users.LoginUserInput
	err:= c.ShouldBindJSON(&input)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"status" : "error",
			"message" : err.Error(),
		})
		return
	}

	loginUser,err:= h.userService.LoginUserInput(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	token,err:= h.authService.GenerateToken(
		loginUser.ID,
		c.Request.Context(),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}


	formatter:= users.FormatUser(loginUser,token)
	
	c.JSON(http.StatusOK, gin.H{
		"status" : "success",
		"message" : "success",
		"data": formatter,
	})
}

func MetaHandler(c *gin.Context)  {
	log.Println("META HANDLER HIT")
	user:= c.MustGet("user").(users.Users)
	c.JSON(200,gin.H{
		"message": "Hi " + user.Name,
	})
	log.Println("STATUS SENT:", c.Writer.Status())
}

