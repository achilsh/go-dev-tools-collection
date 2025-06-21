package jwtwrapperlib

import (
	"fmt"

	"github.com/gin-gonic/gin"
	jwt_lib "github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor, func(token *jwt_lib.Token) (interface{}, error) {
			b := ([]byte(secret))
			return b, nil
		})


		if err != nil || token == nil {
			c.AbortWithError(401, err)
			return 
		}

		if claims, ok := token.Claims.(jwt_lib.MapClaims); ok{
			fmt.Println(claims)
			// fmt.Printf("claims: %v\n", claims)
			// 	fmt.Printf("%v\n", claims["iss"] )
			// 	fmt.Printf("%v\n", claims["sub"])
			// 	fmt.Printf("%v\n",claims["aaa"])
			// 	fmt.Printf("%v\n",claims["exp"])
		}
	
	}
}