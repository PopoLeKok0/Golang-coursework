// Name : Mouad Ben lahbib
// Student number : 300259705
// Project CSI2120/CSI2520
// Winter 2025
// Robert Laganiere, uottawa.ca

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"
)

// movies with rating greater or equal are considered 'liked'
const iLiked float64 = 3.5

// minimum number of users who liked a movie for reliable recommendations
const minUsers int = 10

// number of recommended movies to display
const numRecommendations int = 20

// Define the Recommendation type
type Recommendation struct {
	userID     int     // recommendation for this user
	movieID    int     // recommended movie ID
	movieTitle string  // recommended movie title
	score      float32 // probability that the user will like this movie
	nUsers     int     // number of users who likes this movie
}

// get the probability that this user will like this movie
func (r Recommendation) getProbLike() float32 {
	return r.score / (float32)(r.nUsers)
}

// Define the User type
// and its list of liked items
type User struct {
	userID   int
	liked    []int // list of movies with ratings >= iLiked
	notLiked []int // list of movies with ratings < iLiked
}

func (u User) getUser() int {
	return u.userID
}

func (u *User) setUser(id int) {
	u.userID = id
}

func (u *User) addLiked(id int) {
	u.liked = append(u.liked, id)
}

func (u *User) addNotLiked(id int) {
	u.notLiked = append(u.notLiked, id)
}

// Function to read the ratings CSV file and process each row.
// The output is a map in which user ID is used as key
func readRatingsCSV(fileName string) (map[int]*User, error) {
	// Open the CSV file.
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a CSV reader.
	reader := csv.NewReader(file)

	// Read first line and skip
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	// creates the map
	users := make(map[int]*User, 1000)

	// Read all records from the CSV.
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Iterate over each record and convert the strings into integers or float.
	for _, record := range records {
		if len(record) != 4 {
			return nil, fmt.Errorf("each line must contain exactly 4 integers, but found %d", len(record))
		}

		// Parse user ID integer
		uID, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("error converting '%s' to userID integer: %v", record[0], err)
		}

		// Parse movie ID integer
		mID, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, fmt.Errorf("error converting '%s' to movieID integer: %v", record[1], err)
		}

		// Parse rating float
		r, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("error converting '%s' to rating: %v", record[2], err)
		}

		// checks if it is a new user
		u, ok := users[uID]
		if !ok {
			u = &User{uID, nil, nil}
			users[uID] = u
		}

		// add movie in user list
		if r >= iLiked {
			u.addLiked(mID)
		} else {
			u.addNotLiked(mID)
		}
	}

	return users, nil
}

// Function to read the movies CSV file and process each row.
// The output is a map in which movie ID is used as key
func readMoviesCSV(fileName string) (map[int]string, error) {
	// Open the CSV file.
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a CSV reader.
	reader := csv.NewReader(file)

	// Read first line and skip
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	// creates the map
	movies := make(map[int]string, 1000)

	// Read all records from the CSV.
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Iterate over each record and convert the strings into integers or float.
	for _, record := range records {
		if len(record) != 3 {
			return nil, fmt.Errorf("each line must contain exactly 3 entries, but found %d", len(record))
		}

		// Parse movie ID integer
		mID, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("error converting '%s' to movieID integer: %v", record[0], err)
		}

		// record 1 is the title
		movies[mID] = record[1]
	}

	return movies, nil
}

// checks if value is in the set
func member(value int, set []int) bool {
	for _, v := range set {
		if value == v {
			return true
		}
	}
	return false
}

// compute Jaccard similarity index between two users
func computeSimilarity(u1 *User, u2 *User) float32 {
	// Count common liked and not liked movies
	commonLiked := 0
	commonNotLiked := 0

	// Count likes in common
	for _, m := range u1.liked {
		if member(m, u2.liked) {
			commonLiked++
		}
	}

	// Count dislikes in common
	for _, m := range u1.notLiked {
		if member(m, u2.notLiked) {
			commonNotLiked++
		}
	}

	// Get union size
	// First, collect all movies viewed by both users (removing duplicates)
	allMovies := make(map[int]bool)

	// Add all movies from u1
	for _, m := range u1.liked {
		allMovies[m] = true
	}
	for _, m := range u1.notLiked {
		allMovies[m] = true
	}

	// Add all movies from u2
	for _, m := range u2.liked {
		allMovies[m] = true
	}
	for _, m := range u2.notLiked {
		allMovies[m] = true
	}

	// Calculate Jaccard index
	if len(allMovies) == 0 {
		return 0
	}

	return float32(commonLiked+commonNotLiked) / float32(len(allMovies))
}

// Generator producing Recommendation instances from movie list
func generateMovieRec(wg *sync.WaitGroup, stop <-chan bool, userID int, titles map[int]string) <-chan Recommendation {
	outputStream := make(chan Recommendation)

	go func() {
		defer wg.Done()
		defer close(outputStream)

		for k, v := range titles {
			select {
			case <-stop:
				return
			case outputStream <- Recommendation{userID, k, v, 0.0, 0}:
			}
		}
	}()

	return outputStream
}

// Filter out movies already seen by the user
func filterSeenMovies(wg *sync.WaitGroup, stop <-chan bool, in <-chan Recommendation, user *User) <-chan Recommendation {
	outputStream := make(chan Recommendation)

	go func() {
		defer wg.Done()
		defer close(outputStream)

		for rec := range in {
			// Skip movies the user has already seen (liked or not liked)
			if member(rec.movieID, user.liked) || member(rec.movieID, user.notLiked) {
				continue
			}

			select {
			case <-stop:
				return
			case outputStream <- rec:
			}
		}
	}()

	return outputStream
}

// Filter out movies that have not been liked by at least K users
func filterMinimumUsers(wg *sync.WaitGroup, stop <-chan bool, in <-chan Recommendation,
	users map[int]*User, minK int) <-chan Recommendation {
	outputStream := make(chan Recommendation)

	go func() {
		defer wg.Done()
		defer close(outputStream)

		for rec := range in {
			// Count users who liked this movie
			likedCount := 0
			for _, u := range users {
				if member(rec.movieID, u.liked) {
					likedCount++
				}
			}

			// Skip if less than K users liked it
			if likedCount < minK {
				continue
			}

			// Update the recommendation with the count
			rec.nUsers = likedCount

			select {
			case <-stop:
				return
			case outputStream <- rec:
			}
		}
	}()

	return outputStream
}

// Compute similarity scores for recommendations
func computeScores(wg *sync.WaitGroup, stop <-chan bool, in <-chan Recommendation,
	users map[int]*User, currentUser *User) <-chan Recommendation {
	outputStream := make(chan Recommendation)

	go func() {
		defer wg.Done() // Fixed: This should only be called once per goroutine
		defer close(outputStream)

		for rec := range in {
			score := float32(0)
			count := 0

			// Iterate through all users
			for _, u := range users {
				// Skip current user
				if u.userID == currentUser.userID {
					continue
				}

				// Check if user u liked this movie
				if member(rec.movieID, u.liked) {
					// Compute similarity with current user
					similarity := computeSimilarity(currentUser, u)
					score += similarity
					count++
				}
			}

			// Update recommendation score
			rec.score = score
			rec.nUsers = count

			select {
			case <-stop:
				return
			case outputStream <- rec:
			}
		}
	}()

	return outputStream
}

// Collect and sort recommendations
func collectRecommendations(recs <-chan Recommendation, n int) []Recommendation {
	var allRecs []Recommendation

	// Collect all recommendations
	for rec := range recs {
		if rec.nUsers > 0 { // Ensure we don't divide by zero
			allRecs = append(allRecs, rec)
		}
	}

	// Sort by probability (score/nUsers) in descending order
	sort.Slice(allRecs, func(i, j int) bool {
		return allRecs[i].getProbLike() > allRecs[j].getProbLike()
	})

	// Return top N recommendations
	if len(allRecs) > n {
		return allRecs[:n]
	}
	return allRecs
}

// For comparing with 1 score computing stage
func runWithSingleScoreStage(currentUserID int, titles map[int]string, users map[int]*User) time.Duration {
	stop := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(4) // 4 stages total: generate, filter seen, filter minimum, compute scores

	start := time.Now()

	recChannel := generateMovieRec(&wg, stop, currentUserID, titles)
	filteredChannel1 := filterSeenMovies(&wg, stop, recChannel, users[currentUserID])
	filteredChannel2 := filterMinimumUsers(&wg, stop, filteredChannel1, users, minUsers)
	scoreChannel := computeScores(&wg, stop, filteredChannel2, users, users[currentUserID])

	// Collect and discard results (just timing)
	collectRecommendations(scoreChannel, numRecommendations)

	close(stop)
	wg.Wait()

	return time.Since(start)
}

// Run with 2 parallel score computing stages
func runWithTwoScoreStages(currentUserID int, titles map[int]string, users map[int]*User) ([]Recommendation, time.Duration) {
	stop := make(chan bool)
	var wg sync.WaitGroup

	// We have 6 pipeline stages: generate, filter seen, filter minimum, splitter, 2x score compute, merger
	wg.Add(7)

	start := time.Now()

	// Pipeline stages
	recChannel := generateMovieRec(&wg, stop, currentUserID, titles)
	filteredChannel1 := filterSeenMovies(&wg, stop, recChannel, users[currentUserID])
	filteredChannel2 := filterMinimumUsers(&wg, stop, filteredChannel1, users, minUsers)

	// Split into two channels for parallel processing
	// Buffer to ensure no blocking
	splitChan1 := make(chan Recommendation, 100)
	splitChan2 := make(chan Recommendation, 100)

	// Fan-out splitter
	go func() {
		defer wg.Done()
		defer close(splitChan1)
		defer close(splitChan2)

		counter := 0
		for rec := range filteredChannel2 {
			if counter%2 == 0 {
				splitChan1 <- rec
			} else {
				splitChan2 <- rec
			}
			counter++
		}
	}()

	// Parallel score computation
	scoreChan1 := computeScores(&wg, stop, splitChan1, users, users[currentUserID])
	scoreChan2 := computeScores(&wg, stop, splitChan2, users, users[currentUserID])

	// Merge results channel
	mergedChan := make(chan Recommendation)

	// Fan-in merger
	go func() {
		defer wg.Done()
		defer close(mergedChan)

		var wg2 sync.WaitGroup
		wg2.Add(2)

		// Forward from channel 1
		go func() {
			defer wg2.Done()
			for rec := range scoreChan1 {
				mergedChan <- rec
			}
		}()

		// Forward from channel 2
		go func() {
			defer wg2.Done()
			for rec := range scoreChan2 {
				mergedChan <- rec
			}
		}()

		wg2.Wait()
	}()

	// Collect results
	recommendations := collectRecommendations(mergedChan, numRecommendations)

	close(stop)
	wg.Wait()

	return recommendations, time.Since(start)
}

func main() {
	fmt.Println("Number of CPUs:", runtime.NumCPU())

	// User to be considered
	var currentUser int
	fmt.Println("Recommendations for which user? ")
	fmt.Scanf("%d", &currentUser)

	// Call the function to read and parse the movies CSV file.
	titles, err := readMoviesCSV("movies.csv")
	if err != nil {
		log.Fatal(err)
	}

	// Call the function to read and parse the ratings CSV file.
	ratings, err := readRatingsCSV("ratings.csv")
	if err != nil {
		log.Fatal(err)
	}

	// Get recommendations
	recommendations, executionTime := runWithTwoScoreStages(currentUser, titles, ratings)

	// Print recommendations
	fmt.Printf("\n\nRecommendations for user # %d:\n", currentUser)
	for _, rec := range recommendations {
		fmt.Printf("%s at %.4f [%2d]\n", rec.movieTitle, rec.getProbLike(), rec.nUsers)
	}

	// Print execution time
	fmt.Printf("\n\nExecution time: %s\n", executionTime)
}
