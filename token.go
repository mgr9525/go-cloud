package gocloud

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func CreateToken(claims jwt.MapClaims, tmout time.Duration) (string, error) {
	claims["times"] = time.Now()
	if tmout > 0 {
		claims["timeout"] = time.Now().Add(tmout)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokens, err := token.SignedString([]byte(CloudConf.Token.Key))
	if err != nil {
		return "", err
	}
	return tokens, nil
}
func SetToken(c *gin.Context, p jwt.MapClaims, rem bool, doman ...string) (string, error) {
	if !CloudConf.Token.Enable || CloudConf.Token.Name == "" {
		return "", errors.New("token not enable")
	}
	tmout := time.Hour * 5
	if rem {
		tmout = time.Hour * 24 * 5
	}
	tokens, err := CreateToken(p, tmout)
	if err != nil {
		return "", err
	}
	cke := http.Cookie{
		Value:    tokens,
		Name:     CloudConf.Token.Name,
		HttpOnly: CloudConf.Token.Httponly,
	}
	if CloudConf.Token.Path != "" {
		cke.Path = CloudConf.Token.Path
	}
	if CloudConf.Token.Domain != "" {
		cke.Domain = CloudConf.Token.Domain
	}
	if len(doman) > 0 {
		cke.Domain = doman[0]
	}

	cke.MaxAge = 60 * 60 * 5
	if rem {
		cke.MaxAge = 60 * 60 * 24 * 5
	}
	c.Writer.Header().Add("Set-Cookie", cke.String())
	return tokens, nil
}

func ClearToken(c *gin.Context, doman ...string) error {
	if !CloudConf.Token.Enable || CloudConf.Token.Name == "" {
		return errors.New("token not enable")
	}
	cke := http.Cookie{
		Name:     CloudConf.Token.Name,
		HttpOnly: CloudConf.Token.Httponly,
	}
	if CloudConf.Token.Path != "" {
		cke.Path = CloudConf.Token.Path
	}
	if CloudConf.Token.Domain != "" {
		cke.Domain = CloudConf.Token.Domain
	}
	if len(doman) > 0 {
		cke.Domain = doman[0]
	}
	cke.MaxAge = -1
	c.Writer.Header().Set("Set-Cookie", cke.String())
	return nil
}

func getToken(c *gin.Context) string {
	tkc, err := c.Request.Cookie("gokinstk")
	if err != nil {
		return ""
	}
	return tkc.Value
}
func getTokenAuth(c *gin.Context) string {
	ats := c.GetHeader("Authorization")
	if ats == "" {
		return ""
	}
	aths, err := url.PathUnescape(ats)
	if err != nil {
		return ""
	}
	strings.Replace(aths, "TOKEN ", "", 1)
	return aths
}
func GetTokens(s string) jwt.MapClaims {
	if s == "" {
		return nil
	}
	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		return []byte(CloudConf.Token.Key), nil
	})
	if err == nil {
		claim, ok := token.Claims.(jwt.MapClaims)
		if ok {
			return claim
		}
	}
	return nil
}
func GetToken(c *gin.Context) jwt.MapClaims {
	tk := getTokenAuth(c)
	if tk == "" {
		tk = getToken(c)
	}
	return GetTokens(tk)
}
