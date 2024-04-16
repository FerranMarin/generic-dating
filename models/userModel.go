package models

import (
	"math"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string  `gorm:"type:varchar(255);unique"`
	Password  string  `gorm:"type:varchar(255);not null"`
	Name      string  `gorm:"type:varchar(255);not null"`
	Gender    string  `gorm:"type:varchar(255);not null"`
	Age       int8    `gorm:"not null"`
	Latitude  float64 `gorm:"not null"`
	Longitude float64 `gorm:"not null"`
}

type ShortUser struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Gender         string  `json:"gender"`
	Age            int8    `json:"age"`
	Distance       float64 `json:"distanceFromMe"`
	Attractiveness int     `json:"-"`
}

// Calculate distance between two users using the Haversine formula
func (u *User) Distance(otherUser *User) float64 {
	const earthRadius = 6371 // Earth radius in kilometers

	// Convert latitude and longitude from degrees to radians
	lat1Rad := u.Latitude * math.Pi / 180
	lon1Rad := u.Longitude * math.Pi / 180
	lat2Rad := otherUser.Latitude * math.Pi / 180
	lon2Rad := otherUser.Longitude * math.Pi / 180

	// Calculate differences
	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	// Haversine formula
	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Distance in kilometers
	distance := earthRadius * c
	return distance
}

// Calculate attractiveness elo
func (u *User) Attractiveness(swipes []Swipe) int {
	const defaultScore = 1000
	var attractiveness = defaultScore

	// Calculate score
	for _, swipe := range swipes {
		if swipe.Preference {
			attractiveness += 10
		} else {
			attractiveness -= 10
		}
	}
	return attractiveness
}
