# Gocli

This project consists of a Go module aimed at simplifying the development of applications based on interactive command-line interfaces. The module provides a function that opens a CLI which returns the entered command as a result.
The CLI includes features such as autocompletion, suggestions, command history, and shortcuts like CTRL+C for exit the program safety and CTRL+L for clearing the screen.

## Example

```go
package main

import (
  "fmt"

  "github.com/vcharco/gocli"
)

func main() {

  // Configuration
  cli := gocli.Terminal{
    Prompt:          "GOH> ",
    Options:         []string{"bar", "foo", "clear-history", "print-history", "exit"},
    HistoryId:       "main",
    ExitMessage:     "Good bye baby!",
    BypassCharacter: ":",
  }

  loop := true
  for loop {
    cmd, err := cli.Get()

    if err != nil {
      fmt.Printf("something went wrong: %v\n", err)
      break
    }

    switch cmd {
    case "bar":
      // somebarLogicHere()
    case "foo":
      // someFooLogicHere()
    case "clear-history":
      cli.ClearHistory()
    case "print-history":
      cli.PrintHistory(20)
    case "exit":
      loop = false
    case "":
      // Special case: Executed OS command
    default:
      fmt.Println("Invalid command")
    }
  }
}
```

## Configuration options

- `Prompt`: This is the text at the beggining of the line.
- `Options`: This list of options is used for autocompletion and suggestions.
- `HistoryId`: This can be any identifier. CLIs that share the same identifier also share the same command history.
- `ExitMessage`: Message prompted when user press CTRL+C
- `BypassCharacter`: Gocli checks if the input starts with this character, and in that case, instead of processing it, it sends it directly to the operating system's console. This allows you to execute OS commands without leaving Gocli.
  - Example for BypassCharacter `!`: `Prompt> !ls -l`

## Special outputs

Gocli always returns the command typed by the user (with autocompletion applied) or an error if something goes wrong (though this should never happen). The only exception is when an OS command is executed (indicated by the BypassCharacter). In this case, it returns an empty string.
