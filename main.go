package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	csvFileName := flag.String("csv", "questions.csv", "A csv file with format 'Question,answer'")
	timeLimit := flag.Int("limit", 30, "time limit for the quiz in seconds")
	shuffle := flag.Bool("Shuffle", true, "choose if you want to shuffle your questions before quiz")

	flag.Parse()

	// Reading the CSV and Parsing

	file := openFile(*csvFileName)

	r := csv.NewReader(file)

	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to parse the provided csv file."))
	}

	problems := parseLines(lines)

	if *shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(problems), func(i, j int) { problems[i], problems[j] = problems[j], problems[i] })
	}

	// Timer
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correctAnswers := 0
	ended := showProblems(&correctAnswers, problems, timer)

	if ended {
		showResults(correctAnswers, len(problems))
	}

}

func openFile(s string) *os.File {
	file, err := os.Open(s)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the csv file: %s", s))
	}
	return file
}

func showResults(correct, nProblems int) {
	fmt.Printf("\nYou scored %d out of %d\n", correct, nProblems)
}

func showProblems(correct *int, problems []problem, t *time.Timer) bool {
	finish := false

	for i, p := range problems {
		fmt.Printf("Trivia #%d: %s\n", i+1, p.question)

		fmt.Printf("\tOption %v) %v\n", "A", p.option1)
		fmt.Printf("\tOption %v) %v\n", "B", p.option2)
		fmt.Printf("\tOption %v) %v\n", "C", p.option3)
		fmt.Printf("\tOption %v) %v\n", "D", p.option4)

		ansChannel := make(chan string)
		go func() {
			var answer string
			fmt.Printf("Choose your option: ")
			fmt.Scanf("%s\n", &answer)
			ansChannel <- answer
		}()

		select {
		case <-t.C:
			fmt.Println("Time's up!")
			showResults(*correct, len(problems))
			return finish

		case answer := <-ansChannel:
			if answer == p.answer {
				*correct++
			}
		}

		if i+1 == len(problems) {
			finish = true
		}
	}
	return finish
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
			option1:  line[2],
			option2:  line[3],
			option3:  line[4],
			option4:  line[5],
		}
	}

	return ret
}

type problem struct {
	question string
	answer   string
	option1  string
	option2  string
	option3  string
	option4  string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
