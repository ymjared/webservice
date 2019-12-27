package example

import (
	"bufio"
	"fmt"
	"os"
)

//微型版聊天机器人
func ChatRobot() {
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("Please input your name:")

	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Printf("An error occurred:%s\n", err)
		os.Exit(1)
	} else {
		name := input[:len(input)-1]
		fmt.Printf("Hello, %s! What can i do for you?\n", name)
	}

	for {
		input, err := inputReader.ReadString('\n')
		input = input[:len(input)-1]
		if err != nil {
			fmt.Printf("An error occurred %s\n", err)
			continue
		}

		switch input {
		case "":
			continue
		case "who":
			fmt.Println("My name is damei.")
		case "function":
			fmt.Println("I can song,movie,play...")
		case "bye":
			fmt.Println("Bye man.")
			os.Exit(0)
		default:
			fmt.Println("I will try to be better.")
			continue
		}
	}
}
