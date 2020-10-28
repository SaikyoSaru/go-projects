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

var grade_map = map[int]string{
	0: "Inte godkänd, försök igen!\n",
	3: "Godkänd, men bättre kan du!\n",
	4: "Merparten rätt, börjar bli vass\n",
	5: "Alla rätt! Ace! Vi köper lite donuts för att fira!\n",
}

func main() {
	// csvFilename := flag.String("csv", "problems.csv", "a csv in the format of question/answer")
	difficulty := flag.String("difficulty", "Easy", "The desired difficulty")
	timeLimit := flag.Int("limit", 60, "time limit in seconds")
	flag.Parse()

	problems := generateProblems(*difficulty, 10)
	wrong_com, correct_com := collectComments()
	correct := 0
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	tot_tries := 3
problemLoop:
	for i, p := range problems {
		try_cnt := 0
		fmt.Printf("problem #%d: %s \n", i+1, p.q)
	tryLoop:
		for try_cnt < tot_tries {

			answerCh := make(chan string)
			go func() {
				var answer string
				fmt.Printf("Try no:%d: max_tries:%d\n", try_cnt, tot_tries)
				fmt.Scanf("%s\n", &answer)
				answerCh <- answer
			}()

			select {
			case <-timer.C:
				fmt.Println()
				break problemLoop
			case answer := <-answerCh:
				try_cnt++
				if answer == p.a {
					correct++
					fmt.Printf("%s\n", correct_com[rand.Intn(len(correct_com))][0])
					break tryLoop
				} else {
					if try_cnt == tot_tries {
						fmt.Printf("Failure..., the correct answer is: %s \n", p.a)
					} else {
						fmt.Printf("%s\n", wrong_com[rand.Intn(len(wrong_com))][0])
					}
				}
			}
		}

	}
	grade := calculateGrade(correct, len(problems))
	fmt.Printf("Du fick %d av %d rätt. betyg: %s\n", correct, len(problems), grade)
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func generateProblems(difficulty string, noProblems int) []problem {
	i := 0
	problemSet := make([]problem, noProblems)
	for i < noProblems {
		problemSet[i] = generateArimethicProblem(difficulty)
		i++
	}
	return problemSet
}

func collectComments() ([][]string, [][]string) {
	file, err := os.Open("correct_answer.csv")
	if err != nil {
		exit(fmt.Sprintf("Failed to open correct com file\n"))
	}
	r := csv.NewReader(file)
	corr_com, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to parse positive comments csv file.\n"))
	}

	file, err = os.Open("wrong_answer.csv")
	if err != nil {
		exit(fmt.Sprintf("Failed to open csv file\n"))
	}
	r = csv.NewReader(file)
	wron_com, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to parse wrong comments file."))
	}

	return wron_com, corr_com
}

func calculateGrade(correct int, max int) string {
	percent := float32(correct) / float32(max)
	grade := 0
	switch {
	case percent >= 0.85:
		grade = 5
		break
	case percent >= 0.7:
		grade = 4
		break
	case percent >= 0.5:
		grade = 3
		break
	default:
		grade = 0
		break
	}
	return grade_map[grade]
}

func generateArimethicProblem(difficulty string) problem {
	var var1, var2, size int
	var op string
	switch difficulty {

	case "Easy":
		size = 20
		op = string([]byte{"+-"[rand.Intn((2))]})
		break
	case "Medium":
		op = string([]byte{"+-*/"[rand.Intn((4))]})
		size = 10
		break
	case "Hard":
		op = string([]byte{"+-*/"[rand.Intn((4))]})
		size = 20
		break
	}

	var1 = rand.Intn(size)
	var2 = rand.Intn(size)
	return problem{
		q: fmt.Sprintf("%d %s %d =", var1, op, var2),
		a: calculateAnswer(var1, var2, op),
	}
}

func calculateAnswer(var1 int, var2 int, op string) string {
	var ans string
	switch op {
	case "+":
		ans = fmt.Sprintf("%d", var1+var2)
		break
	case "-":
		ans = fmt.Sprintf("%d", var1-var2)
		break
	case "*":
		ans = fmt.Sprintf("%d", var1*var2)
		break
	case "/":
		ans = fmt.Sprintf("%f", float32(var1)/float32(var2))
		break
	}
	return ans
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
