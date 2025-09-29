package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// createConfig is the first function called in main, it checks the existence of a config file.
// if a config exists, it simply returns and if it doesn't it creates a template and exits program execution.
func createConfig() {
	if _, err := os.Stat("config.json"); err != nil {
		logStatus("config.json does not exist, creating now\n")

		c := Config{
			AuthToken:    "CREATE TOKEN ON CARMEN AND PLACE HERE",
			CourseID:     "FIND COURSE ON CARMEN AND PLACE HERE",
			AssignmentID: "FIND ASSIGNMENT ID ON CARMEN AND PLACE HERE",
			LabTitle:     "TITLE OF LAB (HOW MAKEFILE TARGET APPEARS)",
			Students:     []string{"First Last", "First Last", "First Last"},
		}
		data, err := json.MarshalIndent(c, "", "    ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		// Write to file
		err = os.WriteFile("config.json", data, 0644)
		if err != nil {
			fmt.Println("Error writing file:", err)
			return
		}
		logStatus("successfully created config.json, ready to be filled in\n")
		logStatus("exiting, do not run with auto generated file\n")
		os.Exit(0)
	} else {
		logStatus("found config.json\n")
		return
	}
}

// downloadFile downloads a lab submission for a student and unzips it
func downloadFile(url string, s *Student) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	} else if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error with download request\n")
	}
	token := fmt.Sprintf("Bearer %s", auth)
	req, err := http.NewRequest("GET", s.LabDownloadURL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", token)

	resp, err2 := client.Do(req)
	if err2 != nil {
		panic(err2)
	}
	defer resp.Body.Close()

	namePieces := strings.Split(s.Name, " ")
	name := fmt.Sprintf("%s_%s", namePieces[0], namePieces[1])

	fileName := fmt.Sprintf("labs/%s.zip", name)
	f, err3 := os.Create(fileName)
	if err3 != nil {
		panic(err3)
	}
	defer f.Close()

	io.Copy(f, resp.Body)
	dest := fmt.Sprintf("labs/%s", name)
	cmd := exec.Command("unzip", fileName, "-d", dest)
	execErr := cmd.Run()
	if execErr != nil {
		logError(fmt.Sprintf("couldn't unzip (%s)\n", fileName))
	} else {
		logStatus(fmt.Sprintf("unzipped        (%s)\n", fileName))
	}

	return dest
}

func compileLab(path string, labTitle string) {
	fmt.Printf("path to compile at: %s", path)

	d, _ := os.Getwd()
	fmt.Printf("in dir (%s) before changing to (%s)\n", d, path)

	err := os.Chdir(path)
	if err != nil {
		fmt.Printf("couldnt change dir\n")
		panic(err)
	}

	dir, _ := os.Getwd()
	fmt.Printf("in dir (%s) before compile\n", dir)

	compileCmd := exec.Command("make", "-r", labTitle)
	compileCmd.Env = os.Environ()
	compileCmd.Stdout = os.Stdout
	compileCmd.Stderr = os.Stderr
	compErr := compileCmd.Run()
	if compErr != nil {
		fmt.Printf("couldnt compile\n")
	}
	backErr := os.Chdir("../../")
	if backErr != nil {
		fmt.Printf("couldnt hcange back\n")
		panic(backErr)
	}

}

func writeOutFiles(students []Student) []string {
	os.Mkdir("labs", 0755)
	var dirNames []string
	for i, s := range students {
		student := &students[i]
		namePieces := strings.Split(student.Name, " ")
		if _, err := os.Stat(fmt.Sprintf("labs/%s_%s.zip", namePieces[0], namePieces[1])); err == nil {
			logError(fmt.Sprintf("zip download for %s already exists, not downloading\n", s.Name))
		} else {
			if student.LabDownloadURL != "" {
				logStatus(fmt.Sprintf("downloading submission for %s\n", student.Name))
				dirName := downloadFile(s.LabDownloadURL, student)
				dirNames = append(dirNames, dirName)
			}
		}
	}
	return dirNames
}
