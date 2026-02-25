package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// Input prompts for a text input with a message
func Input(message string, defaultValue string) (string, error) {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", message, defaultValue)
	} else {
		fmt.Printf("%s: ", message)
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)
	if input == "" && defaultValue != "" {
		return defaultValue, nil
	}

	return input, nil
}

// Password prompts for a password input (hidden)
func Password(message string) (string, error) {
	fmt.Printf("%s: ", message)

	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println() // Add newline after password input

	return string(bytePassword), nil
}

// Confirm prompts for yes/no confirmation
func Confirm(message string, defaultValue bool) (bool, error) {
	var defaultStr string
	if defaultValue {
		defaultStr = "Y/n"
	} else {
		defaultStr = "y/N"
	}

	fmt.Printf("%s [%s]: ", message, defaultStr)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	input = strings.ToLower(strings.TrimSpace(input))

	if input == "" {
		return defaultValue, nil
	}

	return input == "y" || input == "yes", nil
}

// Select prompts for selection from a list of options
func Select(message string, options []string, defaultIndex int) (int, error) {
	fmt.Println(message)
	for i, option := range options {
		marker := " "
		if i == defaultIndex {
			marker = "*"
		}
		fmt.Printf("%s %d) %s\n", marker, i+1, option)
	}

	fmt.Printf("Please select [%d]: ", defaultIndex+1)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return defaultIndex, nil
	}

	index, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("invalid input: %s", input)
	}

	if index < 1 || index > len(options) {
		return 0, fmt.Errorf("invalid selection: %d", index)
	}

	return index - 1, nil
}
