package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() { // coverage-ignore
	// print a welcome message
	intro()

	// create a channel to indicate when the user wants to quit
	doneChan := make(chan bool)

	// start a goroutine to read user input with os.Stdin explicitly passed in and run program
	go readUserInput(os.Stdin, doneChan)

	// block until the doneChan gets a value
	<-doneChan

	// close the channel
	close(doneChan)

	// say goodbye
	fmt.Println("Goodbye.")
}

// refactoring to pass in io.Reader which has Read method that satisfies the bufio.NewScanner
// This decouples us from hard coded os.Stdin for testing of this function later.
func readUserInput(in io.Reader, doneChan chan bool) {
	scanner := bufio.NewScanner(in)

	// **Exit Condition**: The `checkNumbers` function returns two values:
	// a result message (`res`) and a boolean (`done`). If `done` is `true`,
	// it means the user has entered "q" to quit.
	// The goroutine then sends a `true` signal to the `doneChan` channel and exits the loop,
	// effectively ending the goroutine.

	for {
		res, done := checkNumbers(scanner)
		if done {
			doneChan <- true
			return
		}
		// **Output and Prompt**: If `done` is `false`, it prints the result message (`res`)
		//  and prompts the user for more input.
		fmt.Println(res)
		prompt()
	}
}

func checkNumbers(scanner *bufio.Scanner) (string, bool) {
	// read user input
	scanner.Scan()

	// check to see if the user wants to quit
	if strings.EqualFold(scanner.Text(), "q") {
		return "", true
	}

	// try to convert what the user typed into an int
	numToCheck, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return "Please enter a whole number!", false
	}

	_, msg := isPrime(numToCheck)

	return msg, false
}

func intro() {
	fmt.Println("Is it Prime?")
	fmt.Println("------------")
	fmt.Println("Enter a whole number, and we'll tell you if it is a prime number or not. Enter q to quit.")
	prompt()
}

func prompt() {
	fmt.Print("-> ")
}

func isPrime(n int) (bool, string) {
	// 0 and 1 are not prime by definition
	if n == 0 || n == 1 {
		return false, fmt.Sprintf("%d is not prime, by definition!", n)
	}

	// negative numbers are not prime
	if n < 0 {
		return false, "Negative numbers are not prime, by definition!"
	}

	// use the modulus operator repeatedly to see if we have a prime number
	for i := 2; i <= n/2; i++ {
		if n%i == 0 {
			// not a prime number
			return false, fmt.Sprintf("%d is not a prime number because it is divisible by %d!", n, i)
		}
	}

	return true, fmt.Sprintf("%d is a prime number!", n)
}
