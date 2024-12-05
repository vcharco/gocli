package gocliutils

import "fmt"

type commandHistory struct {
	commands      []string
	currentIndex  int
	cache         string
	isCacheActive bool
}

var commandHistoryStore map[string]commandHistory = make(map[string]commandHistory)

func GetCommandHistory(cli string) *commandHistory {

	if history, exists := commandHistoryStore[cli]; exists {
		return &history
	}

	newHistory := commandHistory{[]string{}, 0, "", false}
	newHistory.ResetIndex()
	commandHistoryStore[cli] = newHistory

	return &newHistory
}

func (c *commandHistory) Append(command string) {
	c.commands = append(c.commands, command)
	c.ResetIndex()
}

func (c *commandHistory) Clear() {
	c.commands = []string{}
	c.ResetIndex()
}

func (c *commandHistory) GetPrev(currentCommand string) (string, error) {
	if c.currentIndex <= 0 {
		return "", fmt.Errorf("no previous commands")
	}

	if c.currentIndex == len(c.commands) {
		c.cache = currentCommand
		c.isCacheActive = true
	}

	c.currentIndex -= 1
	return c.commands[c.currentIndex], nil
}

func (c *commandHistory) GetNext() (string, error) {
	if c.currentIndex >= len(c.commands)-1 {
		if c.isCacheActive {
			c.isCacheActive = false
			c.currentIndex += 1
			return c.cache, nil
		}
		return "", fmt.Errorf("no more commands")
	}

	c.currentIndex += 1
	if c.currentIndex < len(c.commands) {
		return c.commands[c.currentIndex], nil
	}

	return "", fmt.Errorf("already at the most recent command")
}

func (c *commandHistory) ResetIndex() {
	c.currentIndex = len(c.commands)
}

func (c *commandHistory) Count() int {
	return len(c.commands)
}

func (c *commandHistory) PrintHistory(limit int) {
	if limit == 0 || limit > len(c.commands) {
		limit = len(c.commands)
	}

	start := len(c.commands) - limit
	for i := start; i < len(c.commands); i++ {
		fmt.Println(c.commands[i])
	}
}
