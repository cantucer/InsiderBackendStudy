package main

import (
	"context"
	"encoding/json"
	"fmt"
	"insiderbackendstudy/db"
	"insiderbackendstudy/types"
	"insiderbackendstudy/utils"
	"os"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var conn *pgx.Conn
var weekNo int // Week that is so far played.
var simulationDone bool = false
var port int
var frontendEnabled bool

func init() {
	fmt.Println("Hello world!")

	// Loading .env to get database credentials and team details.
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error while loading .env:", err)
		return
	}

	// Connecting to database.
	databaseURL := os.Getenv("DATABASE_URL")
	port, _ = strconv.Atoi(os.Getenv("PORT"))

	conn, err = db.Connect(databaseURL)
	if err != nil {
		fmt.Println("Error while connecting to database:", err)
		return
	}

	if conn == nil || conn.IsClosed() {
		fmt.Println("Failed to connect to the database.")
		return
	}
	fmt.Println("Connected to the database.")

	// Resetting the database.
	reset()

	// Setting up the front end.
	frontendEnabled, _ = strconv.ParseBool(os.Getenv("FRONTEND_ENABLED"))
}

func reset() {
	// Clearing the database.
	fmt.Println("Clearing the database.")

	err := db.CreateTables(conn)
	if err != nil {
		fmt.Println("Error while creating tables:", err)
		return
	}
	err = db.ClearTables(conn)
	if err != nil {
		fmt.Println("Error while clearing tables:", err)
		return
	}

	// Initializing teams and weeks.

	// Loading the team data from .env variable.
	teamData := os.Getenv("TEAMS")
	var parsedTeams []map[string]interface{}
	err = json.Unmarshal([]byte(teamData), &parsedTeams)
	if err != nil {
		fmt.Println("Error while parsing team data:", err)
		return
	}

	for _, team := range parsedTeams {
		teamName, ok := team["name"].(string)
		if !ok {
			fmt.Println("Invalid team name format.")
			continue
		}
		teamStrengthFloat, ok := team["strength"].(float64)
		if !ok {
			fmt.Println("Invalid team strength format.")
			continue
		}
		teamStrength := int(teamStrengthFloat)

		db.AddTeam(conn, teamName, teamStrength)
		fmt.Printf("Added team %s with strength %d.\n", teamName, teamStrength)
	}

	// Create matches.
	teams, _ := db.GetTeams(conn)
	matches := utils.CreateFixture(teams)

	for _, match := range matches {
		db.AddMatch(conn, match.HomeTeam.Name, match.AwayTeam.Name, match.Week)
		fmt.Printf("Added match %s vs %s for week %d.\n", match.HomeTeam.Name, match.AwayTeam.Name, match.Week)
	}

	weekNo = 0
	simulationDone = false
	fmt.Println("Cleared and prepared the database.")
}

func simulateMatch(match types.Match) {
	homeGoals, awayGoals := utils.SimulateMatch(match)

	homeWin := homeGoals > awayGoals
	awayWin := awayGoals > homeGoals

	fmt.Printf("Match: %s vs %s, Predicted Score: %d - %d\n", match.HomeTeam.Name, match.AwayTeam.Name, homeGoals, awayGoals)

	// Update match result in the database.
	db.UpdateMatchResult(conn, match.HomeTeam.Name, match.AwayTeam.Name, homeGoals, awayGoals)

	// Update team stats in the database.
	if homeWin {
		db.UpdateTeamStats(conn, match.HomeTeam.Name, homeGoals, awayGoals, 1, 0, 0)
		db.UpdateTeamStats(conn, match.AwayTeam.Name, awayGoals, homeGoals, 0, 0, 1)
	} else if awayWin {
		db.UpdateTeamStats(conn, match.AwayTeam.Name, awayGoals, homeGoals, 1, 0, 0)
		db.UpdateTeamStats(conn, match.HomeTeam.Name, homeGoals, awayGoals, 0, 0, 1)
	} else {
		db.UpdateTeamStats(conn, match.HomeTeam.Name, homeGoals, awayGoals, 0, 1, 0)
		db.UpdateTeamStats(conn, match.AwayTeam.Name, awayGoals, homeGoals, 0, 1, 0)
	}
}

// Simulate the next week.
func simulateWeek() bool {
	weekNo++

	matches, _ := db.GetMatches(conn, weekNo) // Get matches.
	if len(matches) == 0 {
		// If there are no matches for the week, it means the tournament is done.
		simulationDone = true
		weekNo--     // Decrement week number since no matches were played.
		return false // No more matches to simulate, tournament is done.
	}

	// Simulate each match for the week.
	for _, match := range matches {
		simulateMatch(match)
	}

	return true // Successfully simulated the week.
}

func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	// Reset endpoint to clear the database and reset the simulation.
	router.POST("/reset", func(c *gin.Context) {
		reset()

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Simulation has been reset. Played matches are deleted and week is reset to week 1.",
		})
	})

	// Simulate week and tournament endpoints.
	router.POST("/simulate_week", func(c *gin.Context) {
		if simulationDone || !simulateWeek() { // Simulate a weel.
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Simulation for tournament is already done. Reset to simulate again.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": fmt.Sprintf("Simulation for week %d is done.", weekNo),
		})
	})
	router.POST("/simulate_tournament", func(c *gin.Context) {
		if simulationDone {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Simulation for tournament is already done. Reset to simulate again.",
			})
			return
		}
		for !simulationDone { // Simulate weeks until the tournament is done.
			if !simulateWeek() {
				break
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Simulation for tournament is done.",
		})
	})

	// Endpoints to fetch teams and their statistics.
	router.GET("/teams", func(c *gin.Context) {
		teams, err := db.GetTeams(conn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to fetch teams.",
			})
			return
		}

		c.JSON(http.StatusOK, teams)
	})

	/// Endpoint to fetch all matches, matches for a specific week, last matches, next matches.
	router.GET("/all_matches", func(c *gin.Context) {
		matches, err := db.GetAllMatches(conn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to fetch matches.",
			})
			return
		}

		c.JSON(http.StatusOK, matches)
	})
	router.GET("/matches/:week", func(c *gin.Context) {
		week := c.Param("week")
		if week == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Week parameter is required.",
			})
			return
		}
		weekInt, err := strconv.Atoi(week)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Week parameter must be an integer.",
			})
			return
		}
		matches, err := db.GetMatches(conn, weekInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to fetch matches for the specified week.",
			})
			return
		}
		c.JSON(http.StatusOK, matches)
	})

	// Endpoint to get last matches.
	router.GET("/last_matches", func(c *gin.Context) {
		if weekNo == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "No matches have been played yet.",
			})
			return
		}

		matches, err := db.GetMatches(conn, weekNo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to fetch last matches.",
			})
			return
		}

		c.JSON(http.StatusOK, matches)
	})

	// Endpoint to get next matches.
	router.GET("/next_matches", func(c *gin.Context) {
		if simulationDone {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Simulation for tournament is already done. Reset to simulate again.",
			})
			return
		}

		matches, err := db.GetMatches(conn, weekNo+1)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to fetch next matches.",
			})
			return
		}

		c.JSON(http.StatusOK, matches)
	})

	// Endpoint to predict winning chances.
	router.GET("/predict_chances", func(c *gin.Context) {
		if simulationDone {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Simulation for tournament is already done. Reset to simulate again.",
			})
			return
		}

		teams, _ := db.GetTeams(conn)

		matches, _ := db.GetAllUnplayedMatches(conn)
		if len(matches) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Simulation for tournament is already done. Reset to simulate again.",
			})
			return
		}

		chances := utils.PredictWinningChances(teams, matches)
		c.JSON(http.StatusOK, chances)
	})

	if frontendEnabled {
		router.Static("/frontend", "./frontend")
	}

	// Opening the server.
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run server: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
}
