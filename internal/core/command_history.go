package gocli

import (
	"fmt"
)

type CommandHistory struct {
	Commands      []string
	CurrentIndex  int
	Cache         string
	IsCacheActive bool
}

func (c *CommandHistory) Append(command string) {
	c.Commands = append(c.Commands, command)
	c.ResetIndex()
}

func (c *CommandHistory) Clear() {
	c.Commands = []string{}
	c.ResetIndex()
}

func (c *CommandHistory) GetPrev(currentCommand string) (string, error) {
	if c.CurrentIndex <= 0 {
		return "", fmt.Errorf("no previous commands")
	}

	if c.CurrentIndex == len(c.Commands) {
		c.Cache = currentCommand
		c.IsCacheActive = true
	}

	c.CurrentIndex -= 1
	return c.Commands[c.CurrentIndex], nil
}

func (c *CommandHistory) GetNext() (string, error) {
	if c.CurrentIndex >= len(c.Commands)-1 {
		if c.IsCacheActive {
			c.IsCacheActive = false
			c.CurrentIndex += 1
			return c.Cache, nil
		}
		return "", fmt.Errorf("no more commands")
	}

	c.CurrentIndex += 1
	if c.CurrentIndex < len(c.Commands) {
		return c.Commands[c.CurrentIndex], nil
	}

	return "", fmt.Errorf("already at the most recent command")
}

func (c *CommandHistory) ResetIndex() {
	c.CurrentIndex = len(c.Commands)
}

func (c *CommandHistory) Count() int {
	return len(c.Commands)
}

func (c *CommandHistory) PrintHistory(limit int) {
	if limit == 0 || limit > len(c.Commands) {
		limit = len(c.Commands)
	}

	start := len(c.Commands) - limit
	for i := start; i < len(c.Commands); i++ {
		fmt.Println(c.Commands[i])
	}
}
