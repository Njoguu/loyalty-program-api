/*
This Package contains middleware functions
*/

package middlewares

import (
	"net/http"
	"api/utils/token"
	"github.com/gin-gonic/gin"
)

/*
JwtAuthMiddleware is a middleware function that ensures JWT token 
authentication for incoming requests.
It returns a gin.HandlerFunc that can be added to the middleware stack 
to protect routes from unauthorized access.
*/
func JwtAuthMiddleware() gin.HandlerFunc{

	// Return a new handler function that takes in a Gin context
	//  and checks for a valid JWT token.
	return func(c *gin.Context){
		err := token.TokenValid(c)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized!",
			})
		}
	}
}