package main

import "fmt"

func writeConclusion(students []Student) {
	logSummary()
	for _, s := range students {
		logStudent(&s)
	}
	fmt.Printf("Wrote %d bytes to ( labs/ )\n", TotalSize)
}

func logStudent(s *Student) {
	unsubmitted := "\033[31mUNSUBMITTED\033[0m"
	submitted := "\033[32mSUBMITTED\033[0m"
	if s.Submitted {
		fmt.Printf("[%s]   %s\n", submitted, s.Name)
	} else {
		fmt.Printf("[%s] %s\n", unsubmitted, s.Name)
	}
}

func logStatus(str string) {
	status := "\x1b\033[33mSTATUS\033[0m"
	fmt.Printf("[%s] %s", status, str)
}

func logSummary() {
	fmt.Printf("\033[1;97;43m[=========== Summary ===========]\033[0m\n")
}

func logError(str string) {
	err := "\033[31mERROR\033[0m"
	fmt.Printf("[%s] %s", err, str)
}
