package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type CanvasUser struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	NameDotNum string `json:"short_name"` // looks like "John Doe (doe.#)"
}

// what we will use to keep track of students for later file work
type Student struct {
	ID             int
	Name           string
	NameDotNum     string
	LabDownloadURL string
	Submitted      bool
}

// relevant download information for any given lab submission
type SubmissionAttachment struct {
	ContentType string `json:"content-type"`
	Size        int    `json:"size"`
	DownloadUrl string `json:"url"`
	DisplayName string `json:"display_name"`
}

type SubmissionHistory struct {
	WorkflowState string                 `json:"workflow_state"`
	Submissions   []SubmissionAttachment `json:"attachments"`
}

var TotalSize = 0

// getUserIDs unmarshals all students in the course_id into an array of CanvasUser,
// mapping each name to an id so we can pick whos id we want and return a slice of students with their ids
func getUserIDs(studentNames []string, data []byte) []Student {
	var users []CanvasUser
	var students []Student
	err := json.Unmarshal(data, &users)
	if err != nil {
		panic(err)
	}

	IDLookup := make(map[string]int, len(users))
	for _, user := range users {
		IDLookup[user.Name] = user.ID
	}

	for i := range studentNames {
		var student Student
		name := studentNames[i]
		id := IDLookup[name]
		student.Name = name
		student.ID = id
		students = append(students, student)
	}

	return students
}

// request is a wrapper function for making http get requests, returns the data if not a download or the read closer for the
// response if its a file download. if the read closer is returned, the caller must still close it.
func request(ctx *Context, url string, isDownload bool) ([]byte, io.ReadCloser) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ctx.Cfg.AuthToken))
	resp, respErr := client.Do(req)
	if respErr != nil {
		panic(respErr)
	}
	if isDownload {
		return nil, resp.Body
	}
	defer resp.Body.Close()
	data, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		panic(readErr)
	}
	return data, nil

}

// buildLabDownloadURL builds the fully qualified url to download a students lab submission and returns it
func buildLabDownloadURL(ctx *Context, s *Student) string {
	return fmt.Sprintf("%s/%s/assignments/%s/submissions/%s?include[]=user&include[]=submission_history",
		baseURL,
		ctx.Cfg.CourseID,
		ctx.Cfg.AssignmentID,
		strconv.Itoa(s.ID))
}

// getStudents makes a request to carmen to receive all student ids in the specified course_id,
// returning a slice of students that are specified in config.json, with their ids
func getStudents(ctx *Context) []Student {
	var users []CanvasUser
	var students []Student
	data, _ := request(ctx, ctx.Api.StudentIdURL, false)
	err := json.Unmarshal(data, &users)
	if err != nil {
		panic(err)
	}
	IDLookup := make(map[string]int, len(users))
	for _, user := range users {
		IDLookup[user.Name] = user.ID
	}

	studentNames := ctx.Cfg.Students
	for i := range studentNames {
		var student Student
		name := studentNames[i]
		id := IDLookup[name]
		student.Name = name
		student.ID = id
		students = append(students, student)
	}
	//ids := getUserIDs(ctx.Cfg.Students, data)
	//return ids
	return students
}

// writeDownloadURL writes a submission download url to a Student, with various checks done along the way
func writeDownloadURL(s *Student, data []byte) {
	// the api will not return an attachments array if the workflow_state is "unsubmitted"
	type State struct {
		WorkflowState string `json:"workflow_state"`
	}
	var state State
	err := json.Unmarshal(data, &state)
	if err != nil {
		panic(err)
	}
	if state.WorkflowState == "submitted" || state.WorkflowState == "graded" {
		var info SubmissionHistory
		json.Unmarshal(data, &info)
		logStatus(fmt.Sprintf("found lab submission for (%s)\n", s.Name))
		fmt.Printf("sub: %v\n", info)
		fmt.Printf("lensubs: %d\n", len(info.Submissions))

		// submission at position 0 will be the most recent
		if len(info.Submissions) != 0 {
			contentType := info.Submissions[0].ContentType
			if contentType == "application/zip" {
				s.LabDownloadURL = info.Submissions[0].DownloadUrl
				TotalSize += info.Submissions[0].Size
				s.Submitted = true
			} else {
				//fmt.Printf("ISSUE: %s did not supply a zip file in their submission\n", s.Name)
				logStatus(fmt.Sprintf("No submission for %s\n", s.Name))
			}
		}
	} else {
		// push student to unsubmitted queue for printing at the end
		s.Submitted = false
	}

}

// getDownloadURLs makes a request to carmen to get urls for students submissions
func getDownloadURLs(ctx *Context, students []Student) {
	//token := fmt.Sprintf("Bearer %s", c.AuthToken)
	for i := range students {
		s := &students[i] // in a 'for i, value' loop we get a copy so we need the address of the student
		if s.ID != 0 {    // id = 0 implies the student has dropped the class
			url := buildLabDownloadURL(ctx, s)
			data, _ := request(ctx, url, false)
			writeDownloadURL(s, data)
		}
	}
}
