package models

import (
	"errors"
	"strconv"
	"strings"
	"time"

	db "../database"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

// AuthAttempt Model
type AuthAttempt struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`

	IP       string `json:"ip" bson:"ip"`
	Username string `json:"username" bson:"username"`
}

// AbuseDetected ...
type AbuseDetected struct {
	MaxIP            string
	MaxIPAndUsername string
}

// Check attempts in login route
func (Config *AbuseDetected) Check(ip, username string) error {

	// Block spammer for one hour
	lastHour := time.Now().Add(-1 * time.Hour)

	getAuthAttempt := db.Connection.Model("AuthAttempt")

	// Check for IPs
	numIP, _ := getAuthAttempt.Find(
		bson.M{
			"ip": ip,
			"createdAt": bson.M{
				"$gt": lastHour,
			},
		},
	).Count()

	// Check for UsernameAndIPs
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

// AttemptCreate Create attempt from unsuccessful auth
func AttemptCreate(ip, username string) error {

	getAuthAttempt := db.Connection.Model("AuthAttempt")

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
