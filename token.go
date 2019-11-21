package gocloud

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/getlantern/errors"
	"gopkg.in/macaron.v1"
	"net/http"
	"strings"
	"time"
)

/*type SSOToken struct {
	Id	string
	Data	string
}*/

func getToken(c *macaron.Context) string {
	if !CloudConf.Token.Enable || CloudConf.Token.Name == "" {
		return ""
	}
	tkc, err := c.Req.Cookie(CloudConf.Token.Name)
	if err != nil {
		return ""
	}
	return tkc.Value
}
func getTokenAuth(c *macaron.Context) string {
	aths := c.Req.Header.Get("Authorization")
	if strings.HasPrefix(aths, "TOKEN ") {
		return strings.Replace(aths, "TOKEN ", "", 1)
	}
	return ""
}

var secret = func(token *jwt.Token) (interface{}, error) {
	return []byte(CloudConf.Token.Key), nil
}

func GetTokens(s string) *jwt.MapClaims {
	if s == "" {
		return nil
	}
	token, err := jwt.Parse(s, secret)
	if err == nil {
		claim, ok := token.Claims.(jwt.MapClaims)
		if ok {
			return &claim
		}
	}
	return nil
}
func GetToken(c *macaron.Context) *jwt.MapClaims {
	tk := getTokenAuth(c)
	if tk == "" {
		tk = getToken(c)
	}
	return GetTokens(tk)
}
func CreateToken(p *jwt.MapClaims, tmout time.Duration) (string, error) {
	if !CloudConf.Token.Enable || CloudConf.Token.Name == "" {
		return "", errors.New("token not enable")
	}
	claims := *p
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
func SetToken(c *macaron.Context, p *jwt.MapClaims, rem bool, doman ...string) (string, error) {
	tmout := time.Hour * 5
	if rem {
		tmout = time.Hour * 24 * 5
	}
	tokens, err := CreateToken(p, tmout)
	if err != nil {
		return "", err
	}
	cke := http.Cookie{Name: CloudConf.Token.Name, Value: tokens, HttpOnly: CloudConf.Token.Httponly}
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
	c.Resp.Header().Add("Set-Cookie", cke.String())
	return tokens, nil
}

func ClearToken(c *macaron.Context, doman ...string) error {
	if !CloudConf.Token.Enable || CloudConf.Token.Name == "" {
		return errors.New("token not enable")
	}

	cke := http.Cookie{Name: CloudConf.Token.Name, Value: "", HttpOnly: CloudConf.Token.Httponly}
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
	c.Resp.Header().Set("Set-Cookie", cke.String())
	return nil
}
