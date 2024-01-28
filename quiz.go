package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)


func main() {
	csvFilename := flag.String("csv","problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "The time limit for the quiz in seconds")
	flag.Parse()

	//Open csv file
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Error opening the CSV file: %s", *csvFilename))
	}

	defer file.Close()

	//Create new csv reader
	reader := csv.NewReader(file)

	//Read all records into a slice
	records, err := reader.ReadAll()
	if err != nil {
		exit("Failed to parse the provided csv file")
	}

	problems := parseRecords(records)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	//Initialize score keeper
	score := 0

	//Process each record as a quiz question
	problemloop:
	for i, p := range problems{
		// Ask the user the question
		fmt.Printf("Problem #%d: %s = \n", i+1,p.q)
		answerChan := make(chan string)
		go func() {
			var answer string
		    fmt.Scanf("%s\n", &answer)
			answerChan <- answer
		}()
		select {
		case <- timer.C:
			fmt.Println()
			break problemloop
		case answer := <- answerChan:
			if answer == p.a {
				score++
			}
		}
}
       fmt.Printf("You scored %d out of %d.\n", score, len(problems))
}

func parseRecords(records [][]string) []problem {
	ret := make([]problem, len(records))
	for i, record := range records {
		ret[i] = problem{
			q: record[0],
			a: strings.TrimSpace(record[1]),
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

//Error handling
func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}


