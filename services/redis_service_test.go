package services

import (
	"fmt"
	"testing"
	"time"
)

func TestSignedToday(t *testing.T) {
	r := GetSignInDateKey("0x1234567890")
	fmt.Println(r)


	location, _ := time.LoadLocation("UTC")
    	currentTime := time.Now().In(location)
    	fmt.Println("Current Time: ",currentTime)

    	activityStartTime := time.Date(2024, 12, 24, 0, 0, 0, 0, location)

    	activityEndTime := time.Date(2025, 1, 6, 23, 59, 59, 0, location)
    	fmt.Println("Activity Start Time: ",activityStartTime)
    	fmt.Println("Activity End Time: ",activityEndTime)
}

