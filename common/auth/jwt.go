package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Ayobami6/common/config"
	pbRider "github.com/Ayobami6/common/proto/riders"
	pbUser "github.com/Ayobami6/common/proto/users"
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
	log.Println(riderID, ok)
    if!ok {
        return -1
    }
    return riderID
}

func RiderAuth(handlerFunc http.HandlerFunc, riderClient pbRider.RiderServiceClient) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		// get token frome request
		tokenString, err := utils.GetTokenFromRequest(r)
		if err != nil {
			log.Println("TokenString error: ", err)
			utils.WriteError(w, http.StatusInternalServerError, "Don't Panic This is From Us!")
			return
		}
		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC);!ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(config.GetEnv("JWT_SECRET", "secret")), nil
        })
		if err!= nil ||!token.Valid {
			log.Println("TokenValid error: ", err)
            Forbidden(w)
            return
        }
		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok {
			log.Println("Claims error: ", err)
			Forbidden(w)
			return
		}
		userIDStr, ok := (*claims)["UserID"].(string)
		if!ok {
			log.Println("UserId error: ", err)
            Forbidden(w)
            return
        }
		userID, err := strconv.Atoi(userIDStr)
		if err!= nil {
			log.Println("Atoi Convert error: ", err)
            utils.WriteError(w, http.StatusInternalServerError, "Don't Panic This is From Us!")
            return
        }
		var ID uint = uint(userID)
		log.Println(ID)
		// get rider by the user ID
		rider, err := riderClient.GetRiderByUserID(r.Context(), &pbRider.RiderUserID{
			UserId: int64(ID),
		})
		log.Println(rider)
		if err != nil {
			log.Println(err)
			Forbidden(w)
			return
		}
		if rider.UserId == 0 {
			log.Println("Rider not found",)
			Forbidden(w)
			return
		}
		riderId, err := strconv.Atoi(rider.RiderId)
		if err!= nil {
            log.Println("Atoi Convert error: ", err)
            utils.WriteError(w, http.StatusInternalServerError, "Don't Panic This is From Us!")
            return
        }
		// save Rider Id to the request context
		ctx := context.WithValue(r.Context(), UserKey, userID)
        ctx = context.WithValue(ctx, RiderKey, riderId)
        r = r.WithContext(ctx)
        handlerFunc(w, r)
	}
}

func UserAuth(handlerFunc http.HandlerFunc, riderClient pbRider.RiderServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := utils.GetTokenFromRequest(r)
		if err!= nil {
            http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
            return
        }
		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC);!ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(config.GetEnv("JWT_SECRET", "secret")), nil
        })
		if err!= nil || !token.Valid {
			log.Println("This is sign token error",err)
            Forbidden(w)
			return
        }
		// get claims
		claims, ok := token.Claims.(*jwt.MapClaims)
		fmt.Println(claims)
		if !ok {
			log.Println("This token claims error", ok)
			Forbidden(w)
            return
		}
		userIDStr, ok := (*claims)["UserID"].(string)
		if !ok {
			log.Println("this userIdstr extract error", ok)
			Forbidden(w)
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err!= nil {
			log.Println(err)
            Forbidden(w)
            return
        }
		var ID uint = uint(userID)
		// get rider by the user ID
		_, err = riderClient.GetRiderByUserID(r.Context(), &pbRider.RiderUserID{
			UserId: int64(ID),
		})
		log.Println("This rider fetch error", err)
		if err == nil {
			Forbidden(w)
            return
		}
		// save User Id to the request context
		ctx := context.WithValue(r.Context(), UserKey, userID)
        handlerFunc(w, r.WithContext(ctx))
	}
}

func Auth(handlerFunc http.HandlerFunc, userClient pbUser.UserServiceClient ) http.HandlerFunc {
	// return the http.Handler function
	return func(w http.ResponseWriter, r *http.Request) {
        tokenString, err := utils.GetTokenFromRequest(r)
		log.Println(tokenString)
		if err!= nil {
            http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
            return
        }
        token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC);!ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(config.GetEnv("JWT_SECRET", "secret")), nil
        })
		log.Println(err)
		log.Println(token)
        if err!= nil ||!token.Valid {
			log.Println("This is sign token error",err)
            Forbidden(w)
            return
        }
        // get claims from jwt
		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok {
			log.Println("This token claims error", ok)
			Forbidden(w)
			return
		}
		userIDStr, ok := (*claims)["UserID"].(string)
		if!ok {
			log.Println("This token Ok",)
            Forbidden(w)
            return
        }
		userID, err := strconv.Atoi(userIDStr)
		if err!= nil {
			log.Println("This is string convert")
            Forbidden(w)
            return
        }
		log.Println(userID)
		// verify user from the database
		_, err = userClient.GetUserByID(r.Context(), &pbUser.UserIDMessage{
			UserId: int64(userID),
		})
        if err!= nil {
			log.Println("This is user get error", err)
            Forbidden(w)
            return
        }
        // set user as context value
        ctx := context.WithValue(r.Context(), UserKey, userID)
        // call the original handler with the updated context
        handlerFunc(w, r.WithContext(ctx))

	}
}
