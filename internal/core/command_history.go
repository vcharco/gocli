package gocli

import (
	"errors"
	"fmt"
)

type commandHistory struct {
	Commands      []string
	CurrentIndex  int
	Cache         string
	IsCacheActive bool
}

func (c *commandHistory) append(command string) {
	c.Commands = append(c.Commands, command)
	c.resetIndex()
}

func (c *commandHistory) clear() {
	c.Commands = []string{}
	c.resetIndex()
}

func (c *commandHistory) getPrev(currentCommand string) (string, error) {
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

func (c *commandHistory) getNext() (string, error) {
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

func (c *commandHistory) resetIndex() {
	c.CurrentIndex = len(c.Commands)
}

func (c *commandHistory) count() int {
	return len(c.Commands)
}

func (c *commandHistory) print(limit int) {
	if limit == 0 || limit > len(c.Commands) {
		limit = len(c.Commands)
	}

	start := len(c.Commands) - limit
	for i := start; i < len(c.Commands); i++ {
		fmt.Println(c.Commands[i])
	}
}

func (c *commandHistory) getAt(index int) (string, error) {
	if index < 0 || index >= len(c.Commands) {
		return "", errors.New("Index out of range")
	}
	return c.Commands[index], nil
}

func (c *commandHistory) getAll() []string {
	return c.Commands
}
