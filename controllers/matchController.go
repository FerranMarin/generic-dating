package controllers

import (
	"strconv"

	"net/http"

	"github.com/FerranMarin/generic-dating/initializers"
	"github.com/FerranMarin/generic-dating/models"
	"github.com/gin-gonic/gin"
)

func SwipeOnUser(c *gin.Context) {
	user, _ := c.Get("user")
	swipeIdParam := c.Param("id")
	_, swipeId := strconv.ParseUint(swipeIdParam, 10, 32)

	var body struct {
		Preference bool `json:"preference"`
	}

	if c.Bind(&body) != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Make sure swipeId is a valid user
	var swipeUser models.User
	swipeUserResult := initializers.DB.First(&swipeUser, swipeId)
	if swipeUserResult.Error != nil && swipeUser.ID == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Swipe id provided is invalid"})
		return
	}

	var existingSwipe models.Swipe
	// Make sure swipe does not already exist
	swipeResult := initializers.DB.Where("user_id = ? and swiped_user_id = ?", user.(models.User).ID, swipeUser.ID).First(existingSwipe)
	if swipeResult.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Checking if swipe existed found an error"})
		return
	}
	if swipeResult.RowsAffected != 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Already swipped on this user"})
	}

	swipe := models.Swipe{
		UserID:       user.(models.User).ID,
		SwipedUserID: swipeUser.ID,
		Preference:   body.Preference,
	}
	result := initializers.DB.Create(&swipe)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed create swipe"})
		return
	}

	// Figure out if reverse swipe exists
	var reverseSwipe models.Swipe
	reverseResult := initializers.DB.Where("user_id = ? and swiped_user_id = ?", swipeUser.ID, user.(models.User).ID).First(reverseSwipe)
	if reverseResult.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Checking if swipe existed found an error"})
		return
	}
	if reverseResult.RowsAffected == 0 {
		c.IndentedJSON(http.StatusOK, gin.H{"matched": false})
	} else {
		match := models.Match{
			UserID1: user.(models.User).ID,
			UserID2: swipeUser.ID,
		}
		matchResult := initializers.DB.Create(&match)
		if matchResult.Error != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed create match"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"matched": true, "matchID": match.ID})
	}
}
