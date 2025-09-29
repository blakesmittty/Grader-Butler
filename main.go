package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// important information needed to make api calls and a list of names we
// can use to compare information against
type Config struct {
	AuthToken    string   `json:"auth_token"`
	CourseID     string   `json:"course_id"`
	AssignmentID string   `json:"assignment_id"`
	LabTitle     string   `json:"lab_title"`
	Students     []string `json:"students"`
}

// certain api enpoints and http info we need
type ApiInfo struct {
	BaseURL      string
	StudentIdURL string
}

// context struct that will get thrown around to the minions
type Context struct {
	Cfg Config
	Api ApiInfo
}

// http client we will use for reaching out to canvas
var client = &http.Client{}

// base url we use in every api call (no trailing slash)
var baseURL = "https://osu.instructure.com/api/v1/courses"

var auth string

// loadConfig builds out and returns the config struct from the config.json file
func loadConfig() Config {
	var c Config
	data, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	jsonErr := json.Unmarshal(data, &c)
	if jsonErr != nil {
		panic(jsonErr)
	}
	return c
}

// buildContext builds out and returns the context struct that will be passed around
func buildContext() Context {
	var ctx Context
	// grab the config.json file
	ctx.Cfg = loadConfig()
	// base url in every api call
	ctx.Api.BaseURL = "https://osu.instructure.com/api/v1/courses" // no trailing slash
	// api endpoint we will use to gather ids for each student
	ctx.Api.StudentIdURL = fmt.Sprintf("%s/%s/users?enrollment_type[]=student&per_page=100", ctx.Api.BaseURL, ctx.Cfg.CourseID)

	auth = ctx.Cfg.AuthToken

	return ctx
}

func main() {
	createConfig()
	ctx := buildContext()
	// studentIDs := getStudents(&ctx)
	// students := getDownloadURLs(&ctx, studentIDs)
	students := getStudents(&ctx)
	getDownloadURLs(&ctx, students)
	writeOutFiles(students)
	writeConclusion(students)

}
