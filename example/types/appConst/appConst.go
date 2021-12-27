package appConst

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	Token string
	jwt.StandardClaims
}

var ApiMap = map[string][]string{

}
