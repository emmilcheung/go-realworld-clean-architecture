package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	models "github.com/gothinkster/golang-gin-realworld-example-app/internal/models"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/utils"
	"github.com/jinzhu/gorm"
)

// Strips 'TOKEN ' prefix from token string
func stripBearerPrefixFromTokenString(tok string) (string, error) {
	// Should be a bearer token
	if len(tok) > 5 && strings.ToUpper(tok[0:6]) == "TOKEN " {
		return tok[6:], nil
	}
	return tok, nil
}

// Extract  token from Authorization header
// Uses PostExtractionFilter to strip "TOKEN " prefix from header
var AuthorizationHeaderExtractor = &request.PostExtractionFilter{
	request.HeaderExtractor{"Authorization"},
	stripBearerPrefixFromTokenString,
}

// Extractor for OAuth2 access tokens.  Looks in 'Authorization'
// header then 'access_token' argument for a token.
var MyAuth2Extractor = &request.MultiExtractor{
	AuthorizationHeaderExtractor,
	request.ArgumentExtractor{"access_token"},
}

// A helper to write user_id and user_model to the context
func (mv *MiddlewareManager) UpdateContextUserModel(c *gin.Context, my_user_id uint, sessionId string) {
	var myUserModel models.User
	db := mv.db
	if my_user_id != 0 {
		db.First(&myUserModel, my_user_id)
	}
	c.Set("my_user_id", my_user_id)
	c.Set("my_user_model", myUserModel)
	c.Set("my_session_id", sessionId)
}

// You can custom middlewares yourself as the doc: https://github.com/gin-gonic/gin#custom-middleware
//  r.Use(AuthMiddleware(true))
func (mv *MiddlewareManager) AuthMiddleware(db *gorm.DB, auto401 bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		mv.UpdateContextUserModel(c, 0, "")
		token, err := request.ParseFromRequest(c.Request, MyAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
			b := ([]byte(utils.NBSecretPassword))
			return b, nil
		})
		if err != nil {
			if auto401 {
				c.AbortWithError(http.StatusUnauthorized, err)
			}
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			session_id := string(claims["session_id"].(string))
			my_user_id := uint(claims["id"].(float64))
			// check token is active
			sess, err := mv.sessUC.GetSessionByID(c.Request.Context(), session_id)
			if err != nil {
				log.Printf("GetSessionByID RequestID: %s, CookieValue: %s, Error: %s",
					utils.GetRequestID(c),
					token.Raw,
					err.Error(),
				)
				c.AbortWithError(http.StatusUnauthorized, err)
				return
			}
			//fmt.Println(my_user_id,claims["id"])
			mv.UpdateContextUserModel(c, my_user_id, sess.SessionID)
		}
	}
}
