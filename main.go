package main

import (
	// "errors"
	"log"
	// "net/http"
	"os"
	// "strings"
	middleware "users/Middleware"
	"users/auth"
	"users/handler"
	"users/users"

	// "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load("./.env")

	dsn:= os.Getenv("DATABASE_URL")
	redisConnect:= os.Getenv("REDIS_PORT")
	db,err:= gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil{
		log.Fatal(err)
	}

	userRepository:= users.NewRepository(db)
	userService:= users.NewService(*userRepository)



	redisClient:= redis.NewClient(&redis.Options{
		Addr: redisConnect,
		DB: 0,
	})

	authService:= auth.NewService(redisClient)
	userHandler:= handler.NewUserHandler(userService,authService)

	router:= gin.Default()
	api:= router.Group("/api/v1")

	api.POST("/users",userHandler.RegisterUser)
	api.POST("/login", userHandler.LoginUser)

	api.GET("/me",middleware.AuthMiddleware(authService,db),handler.MetaHandler)

	router.Run(":3000")
}

// func authMiddleware(authSerice auth.Service, userService users.Service) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		authHeader:= ctx.GetHeader("Authorization")

// 		if ! strings.Contains(authHeader,"Bearer") {
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
// 				"status" : "error",
// 				"message": errors.New("Error"),
// 			})
// 			return
// 		}
// 		var tokenString string
// 		arrayToken:= strings.Split(authHeader," ")
// 		if len(arrayToken) == 2 {
// 			tokenString = arrayToken[1]
// 		}
// 		token,err:= authSerice.ValidateToken(tokenString)
// 		if err != nil {
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
// 				"status":"erros",
// 				"message": err.Error(),
// 			})
// 			return
// 		}
// 		claim,ok:= token.Claims.(jwt.MapClaims)
// 		if ! ok || !token.Valid {
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
// 				"status": "error",
// 				"message": errors.New("Error"),
// 			})
// 			return
// 		}

// 		userID:= int(claim["user_id"].(float64))
// 		user,err:= userService.GetUserByID(userID)
// 		if err != nil {
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
// 				"status":"error",
// 				"message":err.Error(),
// 			})
// 			return
// 		}

// 		ctx.Set("currentUser",user)
// 	}
// }

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzAxNzM1OTM1MDEsImlkIjo5LCJ1c2VyX2lkIjo5fQ.yjsn-VnZuyuVEwGgfwRqGjElxdNn3F4OHDsbJh9Tmhw

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzAxNzM3MzI3NzQsImlkIjo5LCJ1c2VyX2lkIjo5fQ.fA8nQ5OYkrp6ve-1vrZdmlQf9YeZTxM0ay_PrmJ0uX0