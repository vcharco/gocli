# Gocli

This project consists of a Go module aimed at simplifying the development of applications based on interactive command-line interfaces. The module provides a function that opens a CLI which returns the entered command as a result.
The CLI includes features such as autocompletion, suggestions, command history, and shortcuts like CTRL+C for exit the program safety and CTRL+L for clearing the screen.

## Example

```go
package main

import (
 "fmt"
 "strconv"

 "github.com/vcharco/gocli"
)

func main() {

 options := []gocli.Candidate{
  {Name: "foo", DefaultOptionType: gocli.Text, Options: []gocli.CandidateOption{
   {Name: "-f"},
   {Name: "--foo", Type: gocli.Text},
  }},
  {Name: "exit"},
  {Name: "clear-history"},
  {Name: "print-history", DefaultOptionType: gocli.Number},
 }

 // Configuration
 cli := gocli.Terminal{
  Prompt:               "GOH> ",
  Options:              options,
  ExitMessage:          "Have a nice day!",
  BypassCharacter:      ":",
 }

 loop := true
 for loop {
  response := cli.Get()

  // Checking response errors
  if response.Error != nil {
   switch response.Type {
   case gocli.CmdError:
    fmt.Printf("Invalid command: %v\n", response.Error.Error())
   case gocli.ParamError:
    fmt.Printf("Invalid parameters: %v\n", response.Error.Error())
   case gocli.ExecutionError:
    fmt.Printf("Internal error: %v\n", response.Error.Error())
   }
  }

  // Command was successfully executed by cli
  if response.Type == gocli.Cmd {
   switch response.Command {
   case "foo":
    // This is how we get default value if defined a DefaultOptionType
    fooDefault, existsFooDefault := response.Options["default"]

    // This is how we get a flag param (without type)
    _, existsF := response.Options["-f"]

    // This is how we get a param in the rest of the cases
    fooVal, existsFoo := response.Options["--foo"]

    // Default param is always retrieved if we set a DefaultOptionType
    // else, we would get an error, so feel free to avoid this check
    if existsFooDefault {
     fmt.Println("The default value is " + fooDefault)
    }

    // For non default params, we must check if they are retrieved
    if existsF {
     fmt.Println("-f param is set")
    }

    if existsFoo {
     fmt.Printf("The value of --foo is %v\n", fooVal)
    }

   case "clear-history":
    cli.ClearHistory()
   case "print-history":
    // For non Text type params, we may need a cast, but the format will
    // be valid as they were been already checked by the cli
    limit := 0
    value, exists := response.Options["default"]
    if exists {
     limit, _ = strconv.Atoi(value)
    }
    cli.PrintHistory(limit)
   case "exit":
    loop = false
   }
  }

  // Command was executed by OS, but we can perform additional actions
  if response.Type == gocli.OsCmd {
   userInput := response.RawInput
   fmt.Printf("Comman executed by the OS: %v", userInput)
  }

 }
}
```

## Configuration options

- `Prompt`: This is the text at the beggining of the line.
- `Options`: This list of options is used for autocompletion and suggestions. It contains a sublist of valid parameters for each command.
- `BypassCharacter`: Gocli checks if the input starts with this character, and in that case, instead of processing it, it sends it directly to the operating system's console. This allows you to execute OS commands without leaving Gocli.
  - Example for BypassCharacter `!`: `Prompt> !ls -l`
- `ExitMessage`: Prints a nice message when user exits pressing CTRL+C

## Special outputs

Gocli always returns the command typed by the user (with autocompletion applied) or an error if something goes wrong (though this should never happen). The only exceptions are when an OS command is executed (indicated by the BypassCharacter), or when a command is invalid and the AllowInvalidCommands property is set to false (default). In these cases, it returns an empty string.
