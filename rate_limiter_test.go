package RateLimiter

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSlidingLogRateLimiter_CheckLimit_OneValidUser(t *testing.T) {
	var i int
	// Creating users
	userAmogh := &User{UserId: 1, UserName: "Amogh Mishra", WindowSize: 5, Log: []int64{}}

	// Register user for sliding log
	slrt := new(SlidingLogRateLimiter)
	slrt.Users = append(slrt.Users, userAmogh)
	requestReceiver := RequestReceiver{rateLimiter: slrt}
	for i = 0; i < 7; i++ {
		err := requestReceiver.ProcessRequest(userAmogh)
		if i >= 5 {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
		time.Sleep(1 * time.Millisecond / 10)
	}
}

func TestSlidingLogRateLimiter_CheckLimit_TwoValidUsers(t *testing.T) {
	// Creating users
	userAmogh := &User{UserId: 1, UserName: "Amogh Mishra", WindowSize: 5}
	userAkash := &User{UserId: 2, UserName: "Akash Agarwal", WindowSize: 4}

	// Register user for sliding log
	slrt := new(SlidingLogRateLimiter)
	slrt.Users = append(slrt.Users, userAmogh)
	slrt.Users = append(slrt.Users, userAkash)

	requestReceiver := RequestReceiver{rateLimiter: slrt}

	var i int
	for i = 0; i < 7; i++ {
		err := requestReceiver.ProcessRequest(userAmogh)
		if i >= 5 {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
		time.Sleep(1 * time.Millisecond / 10)
	}

	for i = 0; i < 7; i++ {
		err := requestReceiver.ProcessRequest(userAkash)
		if i >= 4 {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
		time.Sleep(1 * time.Millisecond / 10)
	}

}
