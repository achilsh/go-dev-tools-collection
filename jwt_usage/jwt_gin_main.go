package main

import (
	"time"

	jwtwrapperlib "github.com/achilsh/go-dev-tools-collection/jwt_usage/jwt_wrapper_lib"
	"github.com/gin-gonic/gin"
	jwt_lib "github.com/golang-jwt/jwt/v5"
)

var (
	mysupersecretpassword = "unicornsAreAwesome"
)

func main() {
	r := gin.Default()

	public := r.Group("/api")

	public.GET("/", func(c *gin.Context) {
		// Create the token
		token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
		// Set some claims
		token.Claims = jwt_lib.MapClaims{
			"iss":"11111", //颁发者
			"sub":"2222222", //主题
			// 自定义声明，只要不和已经注册的声明名冲突即可
			"aaa":99901,
			// TODO: add other items.
			"Id":  "Christopher",
			"exp": time.Now().Add(time.Second * 10).Unix(),
		}
		// Sign and get the complete encoded token as a string
		tokenString, err := token.SignedString([]byte(mysupersecretpassword))
		if err != nil {
			c.JSON(500, gin.H{"message": "Could not generate token"})
		}
		c.JSON(200, gin.H{"token": tokenString})
	})

	private := r.Group("/api/private")
	private.Use(jwtwrapperlib.Auth(mysupersecretpassword))

	/*
		Set this header in your request to get here.
		Authorization: Bearer `token`
	*/

	private.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from private"})
	})

	r.Run("localhost:8080")
}