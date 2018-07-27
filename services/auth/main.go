package auth

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Gommunity/GoWithWith/config/database"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

type AuthAttempt struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`

	IP       string `json:"ip" bson:"ip"`
	Username string `json:"username" bson:"username"`
}

type AbuseDetected struct {
	MaxIP            string
	MaxIPAndUsername string
}

type (
	JWTUserEmbed struct {
		ID       string `json:"userId"`
		Username string `json:"username"`
	}
	JWTClaims struct {
		Session string `json:"session"`
		SID     string `json:"sid"`
		JWTUserEmbed
		jwt.StandardClaims
	}
)

func (Config *AbuseDetected) Check(ip, username string) error {

	lastHour := time.Now().Add(-1 * time.Hour)
	getAuthAttempt := database.Connection.Model("AuthAttempt")

	numIP, _ := getAuthAttempt.Find(
		bson.M{
			"ip": ip,
			"createdAt": bson.M{
				"$gt": lastHour,
			},
		},
	).Count()
	numIPAndUsername, _ := getAuthAttempt.Find(
		bson.M{
			"ip":       ip,
			"username": username,
			"createdAt": bson.M{
				"$gt": lastHour,
			},
		},
	).Count()

	maxIP, _ := strconv.Atoi(Config.MaxIP)
	maxIPAndUsername, _ := strconv.Atoi(Config.MaxIPAndUsername)

	ipLimitReached := numIP >= maxIP
	ipUserLimitReached := numIPAndUsername >= maxIPAndUsername

	if ipLimitReached || ipUserLimitReached {
		return errors.New("Maximum number of auth attempts reached")
	}

	return nil
}

func AttemptCreate(ip, username string) error {

	getAuthAttempt := database.Connection.Model("AuthAttempt")
	attempt := &AuthAttempt{}
	getAuthAttempt.New(attempt)

	attempt.IP = ip
	attempt.Username = strings.ToLower(username)

	err := attempt.Save()

	if err != nil {
		panic(err)
	}

	return nil
}

func ParseJWT(token string, secret []byte) (*jwt.Token, error) {

	var err error
	token = strings.Replace(token, "Bearer ", "", 1)
	tokenParsed := new(jwt.Token)
	tokenParsed, err = jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return secret, nil
	})

	if err == nil && tokenParsed.Valid {
		return tokenParsed, nil
	}

	return &jwt.Token{}, err
}

func CreateJWToken(session, SID, username, userID string, signingKey []byte) string {

	var SignedToken string
	var err error

	claims := JWTClaims{
		session,
		SID,
		JWTUserEmbed{
			userID,
			username,
		},
		jwt.StandardClaims{},
	}

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if SignedToken, err = tokenJWT.SignedString(signingKey); err != nil {
		panic(err)
	}

	SignedToken = "Bearer " + SignedToken
	return SignedToken
}
