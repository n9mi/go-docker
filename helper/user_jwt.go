package helper

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/n9mi/go-docker/exception"
	"github.com/n9mi/go-docker/model/web"
)

func GenerateLoginToken(userEmail string, roleID uint) (web.Token, error) {
	var token web.Token

	accessToken, errAccess := GenerateAccessToken(userEmail, roleID)
	if errAccess != nil {
		return token, &exception.BadRequestError{Message: errAccess.Error()}
	}

	refreshToken, errRefresh := GenerateRefreshToken()
	if errRefresh != nil {
		return token, &exception.BadRequestError{Message: errRefresh.Error()}
	}

	token = web.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return token, nil
}

func GenerateAccessToken(userEmail string, roleID uint) (string, error) {
	var signedToken string

	envMinutes, errConv := strconv.Atoi(os.Getenv("JWT_ACCESS_KEY_EXPIRE_MINUTES"))
	if errConv != nil {
		return signedToken, errConv
	}
	minutesDuration := time.Duration(envMinutes) * time.Minute

	claims := web.AccessClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer: os.Getenv("APP_NAME"),
			ExpiresAt: time.Now().
				Add(minutesDuration).Unix(),
		},
		Email:  userEmail,
		RoleID: roleID,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, errSign := token.SignedString([]byte(os.Getenv("JWT_ACCESS_KEY_SIGNATURE")))
	if errSign != nil {
		return signedToken, errSign
	}

	return signedToken, nil
}

func GenerateRefreshToken() (string, error) {
	var refreshToken string

	envMinutes, errConv := strconv.Atoi(os.Getenv("JWT_REFRESH_KEY_EXPIRE_MINUTES"))
	if errConv != nil {
		return refreshToken, errConv
	}
	minutesDuration := time.Duration(envMinutes) * time.Minute

	claims := web.RefreshClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    os.Getenv("APP_NAME"),
			ExpiresAt: time.Now().Add(minutesDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, errSign := token.SignedString([]byte(os.Getenv("JWT_REFRESH_KEY_SIGNATURE")))
	if errSign != nil {
		return refreshToken, errSign
	}

	return signedToken, nil
}

func ParseRefreshToken(refreshStr string) (web.RefreshClaims, error) {
	var parsingResult web.RefreshClaims

	// check signing method
	token, errParse := jwt.Parse(refreshStr, func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &exception.TokenError{}
		} else if method != jwt.SigningMethodHS256 {
			return nil, &exception.TokenError{}
		}

		return []byte(os.Getenv("JWT_REFRESH_KEY_SIGNATURE")), nil
	})

	if errParse != nil {
		return parsingResult, errParse
	}

	// check if claims are valid
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return parsingResult, token.Claims.Valid()
	}

	claimsExpiredConv, ok := claims["exp"].(float64)
	if !ok {
		return parsingResult, &exception.TokenError{}
	}

	claimsIssuerConv, ok := claims["iss"].(string)
	if (!ok) && (claimsIssuerConv != os.Getenv("APP_NAME")) {
		return parsingResult, &exception.TokenError{}
	}

	parsingResult.ExpiresAt = int64(claimsExpiredConv)
	parsingResult.Issuer = claimsIssuerConv

	return parsingResult, nil
}

func ParseAccessToken(accessStr string) (web.AccessClaims, error) {
	var parsingResult web.AccessClaims

	// check signing method
	token, errParse := jwt.Parse(accessStr, func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return parsingResult, &exception.TokenError{}
		} else if method != jwt.SigningMethodHS256 {
			return parsingResult, &exception.TokenError{}
		}

		return []byte(os.Getenv("JWT_ACCESS_KEY_SIGNATURE")), nil
	})

	if errParse != nil {
		return parsingResult, &exception.TokenError{}
	}

	// check if claims valid
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return parsingResult, token.Claims.Valid()
	}

	claimsEmailConv, ok := claims["email"].(string)
	if !ok {
		return parsingResult, &exception.TokenError{}
	}

	claimsRoleIDConv, ok := claims["role_id"].(float64)
	if !ok {
		return parsingResult, &exception.TokenError{}
	}

	claimsExpiredConv, ok := claims["exp"].(float64)
	if !ok {
		return parsingResult, &exception.TokenError{}
	}

	claimsIssuerConv, ok := claims["iss"].(string)
	if (!ok) && (claimsIssuerConv != os.Getenv("APP_NAME")) {
		return parsingResult, &exception.TokenError{}
	}

	parsingResult.Email = claimsEmailConv
	parsingResult.RoleID = uint(claimsRoleIDConv)
	parsingResult.ExpiresAt = int64(claimsExpiredConv)
	parsingResult.Issuer = claimsIssuerConv

	return parsingResult, nil
}
