package auth

import (
	"os"

	"github.com/Gommunity/GoWithWith/services/auth"
)

func InitConfig() {
	auth.AbuseDetected = &auth.AbuseDetected{
		MaxIP:            os.Getenv("AuthAttemptsForIp"),
		MaxIPAndUsername: os.Getenv("AuthAttemptsForIpAndUser"),
	}
}
