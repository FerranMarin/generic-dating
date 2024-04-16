package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/FerranMarin/generic-dating/helpers"
	"github.com/FerranMarin/generic-dating/initializers"
	"github.com/FerranMarin/generic-dating/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)
	if user.ID == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		// Token expires after 24h
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to generate JWT"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"token": tokenString})
}

func CreateRandom(c *gin.Context) {

	type Coordinates struct {
		Lat float64
		Lng float64
	}

	var (
		GENDERS      = []string{"male", "female"}
		MALE_NAMES   = []string{"Noah", "Theo", "Oliver", "George", "Leo", "Freddie", "Arthur", "Archie", "Alfie", "Charlie", "Oscar", "Henry", "Harry", "Jack", "Teddy", "Finley", "Arlo", "Luca", "Jacob", "Tommy", "Lucas", "Theodore", "Max", "Isaac", "Albie", "James", "Mason", "Rory", "Thomas", "Rueben", "Roman", "Logan", "Harrison", "William ", "Elijah", "Ethan", "Joshua", "Hudson", "Jude", "Louie", "Jaxon", "Reggie", "Oakley", "Hunter", "Alexander", "Toby", "Adam", "Sebastian", "Daniel", "Ezra", "Rowan", "Alex", "Dylan", "Ronnie", "Kai", "Hugo", "Louis", "Riley", "Edward", "Finn", "Grayson", "Elliot", "Caleb", "Benjamin", "Bobby", "Frankie", "Zachary", "Brody", "Jackson", "Ollie", "Jasper", "Liam", "Stanley", "Sonny", "Blake", "Albert", "Joseph", "Chester", "Carter", "David", "Milo", "Ellis", "Jenson", "Samuel", "Gabriel", "Eddie", "It is", "Rupert", "Eli", "Myles", "Brodie", "Parker", "Ralph", "Miles", "Jayden", "Billy", "Elliott", "Jax", "Ryan", "Joey"}
		FEMALE_NAMES = []string{"Olivia", "Amelia", "Isla", "Lily", "Ava", "Freya", "Ivy", "Milly", "Mia", "Sophia", "Poppy", "Emily ", "Willow", "Grace", "Evie", "Elsie", "Rosie", "Isabella", "Daisy", "Sienna", "Florence", "Ella", "Charlotte", "Harper ", "Phoebe", "Ruby", "Sofia", "Sophie", "Evelyn", "Maisie", "Emilia", "Aria", "Matilda", "Maya", "Luna", "Esme", "Hallie", "Alice", "Lottie", "Mila", "Isabelle", "Violet", "Ellie", "Aurora", "Maeve", "Scarlett", "Delilah", "Ada", "Bonnie", "Penelope", "Ayla", "Erin", "Layla", "Chloe", "Arabella", "Eva", "Mabel", "Eliza", "Rose", "Thea", "Robyn", "Molly ", "Imogen", "Nancy", "Zara", "Harriet", "Bella", "Gracie", "Emma", "Eleanor", "Lucy", "Eden", "Lyla", "Jessica", "Lyra", "Darcie", "Heidi", "Elizabeth", "Orla", "Iris", "Elodie", "Margot", "Hannah", "Lola", "Ophelia", "Eloise", "Maddison", "Lara", "Hazel", "Frankie", "Pippa", "Lilly", "Nova", "Autumn", "Clara", "Nellie", "Myla", "Athena", "Amelie", "Niamh"}
		CITIES       = []Coordinates{
			{Lat: 51.509865, Lng: -0.118092}, // London
			{Lat: 51.752022, Lng: -1.257677}, // Oxford
			{Lat: 54.607868, Lng: -5.926437}, // Belfast
			{Lat: 55.860916, Lng: -4.251433}, // Glasgow
			{Lat: 51.568535, Lng: -1.772232}, // Swindon
			{Lat: 53.801277, Lng: -1.548567}, // Leeds
			{Lat: 53.383331, Lng: -1.466667}, // Sheffield
		}
	)

	plain_password := helpers.GenerateRandomString(8)

	hash, err := bcrypt.GenerateFromPassword(plain_password, 10)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
		return
	}

	gender := GENDERS[rand.Intn(len(GENDERS))]
	var name string
	if gender == "male" {
		name = MALE_NAMES[rand.Intn(len(MALE_NAMES))]
	} else {
		name = FEMALE_NAMES[rand.Intn(len(FEMALE_NAMES))]
	}
	email := strings.ToLower(name) + string(helpers.GenerateRandomString(4)) + "@fakemail.com"
	city := CITIES[rand.Intn(len(CITIES))]

	user := models.User{
		Email:     email,
		Password:  string(hash),
		Name:      name,
		Gender:    gender,
		Age:       int8(rand.Intn(100-18) + 18),
		Latitude:  city.Lat,
		Longitude: city.Lng,
	}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed create user"})
		return
	}

	// Change user password for the plain one to return it
	user.Password = string(plain_password)

	c.IndentedJSON(http.StatusCreated, user)
}

func DiscoverUsers(c *gin.Context) {
	user, _ := c.Get("user")

	// Get query parameters
	minAgeStr := c.Query("min_age")
	maxAgeStr := c.Query("max_age")
	gender := c.Query("gender")

	// Parse minAge parameter
	var minAge int8
	if minAgeStr != "" {
		age64, err := strconv.ParseInt(minAgeStr, 10, 8)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid max_age parameter"})
			return
		}
		minAge = int8(age64)
	}

	// Parse maxAge parameter
	var maxAge int8
	if maxAgeStr != "" {
		age64, err := strconv.ParseInt(maxAgeStr, 10, 8)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid min_age parameter"})
			return
		}
		maxAge = int8(age64)
	}

	if maxAge != 0 && minAge != 0 && maxAge <= minAge {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "max_age parameter has to be greater than min_age"})
		return
	}

	// Get all users user has already swiped
	var userSwipes []models.Swipe

	swipesResult := initializers.DB.Where("user_id = ?", user.(models.User).ID).Find(&userSwipes)
	if swipesResult.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed retreive user swipes"})
		return
	}
	// Extract the IDs of the users the current user has swiped on
	var swipedUserIDs []uint
	for _, swipe := range userSwipes {
		swipedUserIDs = append(swipedUserIDs, swipe.SwipedUserID)
	}

	// Build query conditions based on optional filters
	query := initializers.DB.Not("ID", append(swipedUserIDs, user.(models.User).ID))
	switch {
	case maxAge != 0 && minAge != 0:
		query = query.Where("age >= ? AND age <= ?", minAge, maxAge)
	case maxAge != 0:
		query = query.Where("age >= ?", maxAge)
	case minAge != 0:
		query = query.Where("age <= ?", minAge)
	}
	if gender != "" {
		query = query.Where("gender = ?", gender)
	}

	var users []models.User
	result := query.Find(&users)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed retreive users"})
		return
	}

	requestingUser := user.(models.User)

	// Transform Users to ShortUser
	var shortUsers []models.ShortUser
	for _, u := range users {
		// Get all swipes given to this user
		var swipes []models.Swipe
		result := initializers.DB.Where("swiped_user_id = ?", u.ID).Find(&swipes)
		if result.Error != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed retreive user swipes"})
			return
		}
		shortUsers = append(shortUsers, models.ShortUser{ID: u.ID, Name: u.Name, Gender: u.Gender, Age: u.Age, Distance: requestingUser.Distance(&u), Attractiveness: u.Attractiveness(swipes)})
	}

	// sort by distance and attractiveness
	sort.Slice(shortUsers, ByDistanceAndAttractiveness(shortUsers))

	c.IndentedJSON(http.StatusOK, gin.H{"results": shortUsers})
}

func ByDistanceAndAttractiveness(users []models.ShortUser) func(i, j int) bool {
	return func(i, j int) bool {
		// First, compare by distance
		if users[i].Distance < users[j].Distance {
			return true
		} else if users[i].Distance > users[j].Distance {
			return false
		}

		// If distance is equal, compare by attractiveness
		return users[i].Attractiveness > users[j].Attractiveness
	}
}
