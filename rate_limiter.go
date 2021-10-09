package RateLimiter

import (
	"errors"
	"time"
)

type RateLimiter interface {
	CheckLimit(user *User, t int64) error
}

type User struct {
	UserId     int
	UserName   string
	WindowSize int
	Log        []int64
}

type RequestReceiver struct {
	rateLimiter RateLimiter
}

func (rr RequestReceiver) ProcessRequest(user *User) error {
	t := int64(time.Now().Nanosecond())
	err := rr.rateLimiter.CheckLimit(user, t)
	if err != nil {
		println("Error occured for this user", user.UserId, " :timestamp ", t)
		return err
	} else {
		time.Sleep(1 * time.Nanosecond)
		println("Request went through for this user", user.UserId, " :timestamp ", t)
	}
	return nil
}

type SlidingLogRateLimiter struct {
	Users []*User
}

func (slr SlidingLogRateLimiter) CheckLimit(user *User, t int64) error {
	var userFound bool
	var i int
	for _, u := range slr.Users {
		if user.UserId == u.UserId {
			userFound = true
		}
	}

	// Return if user is not enabled for sliding log rate limit
	if !userFound {
		return errors.New("rate limiting not enabled for user")
	}

	// Remove expired requests from the user log queue
	for i = 0; i < len(user.Log); i++ {
		if user.Log[i]+int64(1*time.Millisecond) < int64(time.Now().Nanosecond()) {
			user.Log = append(user.Log[:i], user.Log[i+1:]...)
		}
	}

	// Check if request is being throttled or allowed
	if len(user.Log) >= user.WindowSize {
		return errors.New("request throttled")
	} else {
		user.Log = append(user.Log, t)
	}
	return nil
}
