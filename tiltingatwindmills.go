package main

import "github.com/ChimeraCoder/anaconda"
import "gopkg.in/gin-gonic/gin.v1"
import "net/http"
import "net/url"
import "os"
import "log"
import "strconv"

func setup() {
	anaconda.SetConsumerKey(os.Getenv("TW_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TW_CONSUMER_SECRET"))
}

func extractTweets(c *gin.Context, timeline []anaconda.Tweet) []string {

	var tweets []string
	for _, tweet := range timeline {
		tweets = append(tweets, tweet.Text)
	}

	return tweets
}

// Get a list of friends from Twitter, put them in a map if setting,
// else add to a list if found in map
func getFriends(api *anaconda.TwitterApi, friends map[string]bool, userId string, setOrCheck bool) []string {

	var mutual []string

	values := url.Values{}
	values.Set("user_id", userId)

	ch := api.GetFriendsListAll(values)

	for friendPage := range ch {
		if friendPage.Error == nil {
			for _, user := range friendPage.Friends {
				id := strconv.FormatInt(user.Id, 10)
				log.Printf("name = %s, id = %s", user.Name, id)
				if setOrCheck {
					friends[id] = true
				} else {
					_, found := friends[id]
					if found {
						mutual = append(mutual, user.Name)
					}
				}
			}
		}
	}

	return mutual
}

func main() {
	setup()

	router := gin.Default()

	//Swagger docs
	router.Static("/docs", "../doc/dist")

	// This handler will match /tweets/john but will not match either /tweets/ or /tweets
	router.GET("/tweets/:name", func(c *gin.Context) {
		name := c.Param("name")
		api := anaconda.NewTwitterApi(os.Getenv("TW_ACCESS_TOKEN"), os.Getenv("TW_ACCESS_TOKEN_SECRET"))
		values := url.Values{}
		values.Add("screen_name", name)
		timeline, err := api.GetUserTimeline(values)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
			return
		}
		tweets := extractTweets(c, timeline)
		c.JSON(http.StatusOK, tweets)
	})

	// This route will match /common/abby/boris
	router.GET("/common/:name/:other", func(c *gin.Context) {
		names := c.Param("name")
		other := c.Param("other")
		names = names + ","
		names = names + other
		api := anaconda.NewTwitterApi(os.Getenv("TW_ACCESS_TOKEN"), os.Getenv("TW_ACCESS_TOKEN_SECRET"))
		values := url.Values{} // Would like to stringify ids but not sure if anaconda will handle that ...
		users, err := api.GetUsersLookup(names, values)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error}) // TODO Wrong status code!
			return
		}
		if len(users) < 2 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Could not find both users by name"})
			return
		}
		if len(users) != 2 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Too many users with those names"})
			return
		}

		// This map contains the names of the friends
		//
		var friends map[string]bool
		friends = make(map[string]bool)

		// Get the friends for the smaller of the two
		var getUserIndex = 1
		if users[0].FriendsCount < users[1].FriendsCount {
			getUserIndex = 0
		}

		userId := strconv.FormatInt(users[getUserIndex].Id, 10)
		log.Printf("Getting friends for index %d, user %s\n", getUserIndex, userId)
		_ = getFriends(api, friends, userId, true) // true = set value

		// Now get the other user's friends
		getUserIndex = (getUserIndex + 1) % 2

		userId = strconv.FormatInt(users[getUserIndex].Id, 10)
		log.Printf("Getting friends for index %d, user %s\n", getUserIndex, userId)
		mutual := getFriends(api, friends, userId, false) // false = check values

		c.JSON(http.StatusOK, mutual)
	})

	router.Run(":8000")
}
