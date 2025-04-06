package output

import "fmt"

type ConsoleOutput struct{}

func NewConsoleOutput() *ConsoleOutput {
	return &ConsoleOutput{}
}

func (c *ConsoleOutput) Write(msg string) error {
	fmt.Println(msg)
	return nil
}
