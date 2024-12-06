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
    Prompt:                "GOH> ",
    Options:               []string{"bar", "foo", "clear-history", "print-history", "exit"},
    ExitMessage:           "Have a nice day!",
    BypassCharacter:       ":",
    AllowInvalidCommands:  false,
    InvalidCommandMessage: "Invalid command!",
  }

  loop := true
  for loop {
    cmd, err := cli.Get()

    if err != nil {
      fmt.Printf("something went wrong: %v\n", err)
      break
    }

    switch cmd {
    case "foo":
      // someFooLogicHere()
    case "bar":
      // someBarLogicHere()
    case "clear-history":
      cli.ClearHistory()
    case "print-history":
      cli.PrintHistory(20)
    case "exit":
      loop = false
    case "":
      // Command by OS or invalid command with AllowInvalidCommands=false
    default:
      // Not reachable because AllowInvalidCommands is set to false
    }
  }
}
```

## Configuration options

- `Prompt`: This is the text at the beggining of the line.
- `Options`: This list of options is used for autocompletion and suggestions.
- `ExitMessage`: Message prompted when user press CTRL+C
- `BypassCharacter`: Gocli checks if the input starts with this character, and in that case, instead of processing it, it sends it directly to the operating system's console. This allows you to execute OS commands without leaving Gocli.
  - Example for BypassCharacter `!`: `Prompt> !ls -l`
- `AllowInvalidCommands`: If this property is set to false (default), only commands included in the Options property are accepted. If the command is followed by parameters separated by spaces, it is considered valid, as only the command itself is checked. Otherwise, if the property is set to true, the command is returned regardless of whether it is valid.
- `InvalidCommandMessage`: This message is displayed when an invalid command is executed and the AllowInvalidCommands property is set to false.

## Special outputs

Gocli always returns the command typed by the user (with autocompletion applied) or an error if something goes wrong (though this should never happen). The only exceptions are when an OS command is executed (indicated by the BypassCharacter), or when a command is invalid and the AllowInvalidCommands property is set to false (default). In these cases, it returns an empty string.
