// middlewares/authz.go

package middlewares

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/lokesh20018/iitk-coin/auth"

	"github.com/gin-gonic/gin"
)

//validates token and authorizes users
func Authz() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(403, "No Authorization header provided")
			c.Abort()
			return
		}

		extractedToken := clientToken

		clientToken = strings.TrimSpace(extractedToken)

		jwtWrapper := auth.JwtWrapper{
			SecretKey: "verysecretkey",
			Issuer:    "AuthService",
		}

		claims, err := jwtWrapper.ValidateToken(clientToken)
		if err != nil {
			c.JSON(401, err.Error())
			c.Abort()
			return
		}

		c.Set("roll_no", claims.Roll_no)

		c.Next()

	}
}

func Authz_Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(403, "No Authorization header provided")
			c.Abort()
			return
		}

		err := godotenv.Load(".env")

		if err != nil {
			log.Println("Error loading .env file")
		}
		pass := os.Getenv("admin")
		// pass is admin007
		if clientToken == pass {
			c.Next()
		} else {
			c.JSON(403, "invaid admin credentials")
			c.Abort()
			return

		}

	}
}
