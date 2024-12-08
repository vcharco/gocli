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

func (t *Terminal) PrintHistory(limit int) {
	if limit == 0 || limit > len(t.commandHistory.Commands) {
		limit = len(t.commandHistory.Commands)
	}

	start := len(t.commandHistory.Commands) - limit
	for i := start; i < len(t.commandHistory.Commands); i++ {
		fmt.Println(t.commandHistory.Commands[i])
	}
}

func (t *Terminal) ClearHistory() {
	t.commandHistory.clear()
}

func (t *Terminal) CountHistory() int {
	return len(t.commandHistory.Commands)
}

func (t *Terminal) GetHistoryAt(index int) (string, error) {
	if index < 0 || index >= len(t.commandHistory.Commands) {
		return "", errors.New("index out of range")
	}
	return t.commandHistory.Commands[index], nil
}

func (t *Terminal) GetHistory(index int) []string {
	return t.commandHistory.Commands
}
