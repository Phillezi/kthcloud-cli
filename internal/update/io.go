package update

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func PromptYesNo(prompt string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s (y/n): ", prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("error reading input: %v", err)
		}

		input = strings.TrimSpace(strings.ToLower(input))

		if input == "y" || input == "yes" {
			return true, nil
		} else if input == "n" || input == "no" {
			return false, nil
		} else {
			fmt.Println("Invalid input, please enter 'y' or 'n'.")
		}
	}
}
