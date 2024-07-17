package auth

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Ayobami6/common/config"
	"github.com/Ayobami6/common/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

var UserKey contextKey= "UserID"

const RiderKey contextKey = "RiderID"


func CreateJWT(secret []byte, userId int) (string, error) {
	exp, err := strconv.ParseInt(config.GetEnv("JWT_EXPIRATION", "25000"), 10, 64)
	if err != nil {
		return "", err
	}
	expiration := time.Second * time.Duration(exp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserID": strconv.Itoa(userId),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenStr, nil

}

func Forbidden(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, "Unauthorized",)
}


func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
    if !ok {
        return -1
    }
    return userID
}

func GetRiderIDFromContext(ctx context.Context) int {
	riderID, ok := ctx.Value(RiderKey).(int)
    if!ok {
        return -1
    }
    return riderID
}
