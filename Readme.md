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
  {Name: "foo", Options: []gocli.CandidateOption{
   {Name: "-f"},
   {Name: "--foo", Type: gocli.Text},
   {Name: "default", Type: gocli.Number, Modifier: gocli.DEFAULT},
  }},
  {Name: "exit"},
  {Name: "clear-history"},
  {Name: "print-history", Options: []gocli.CandidateOption{
   {Name: "default", Type: gocli.Number, Modifier: gocli.DEFAULT | gocli.REQUIRED},
  }},
 }

 // Configuration
 cli := gocli.Terminal{
  Prompt:          "GOH> ",
  Options:         options,
  ExitMessage:     "Have a nice day!",
  BypassCharacter: ":",
 }

 loop := true
 for loop {

  // Gets the user input
  response := cli.Get()

  // Checking response errors
  if response.Error != nil {
   fmt.Println()
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
   fmt.Printf("Comman executed by the OS: %v\n", userInput)
  }

 }
}
```

## Configuration options

- `Prompt`: This is the text at the beggining of the line.
- `Options`: This list of options is used for autocompletion and suggestions. It contains a sublist of valid parameters for each command.
- `BypassCharacter`: Gocli checks if the input starts with this character, and in that case, instead of processing it, it sends it directly to the operating system's console. This allows you to execute OS commands without leaving Gocli.
  - Example for BypassCharacter `:`: `Prompt> :ls -l`
- `ExitMessage`: Prints a nice message when user exits pressing CTRL+C

### Options

**Options** (Candidate) are the commands available for your custom cli. Each command must be provided with a Name. Optionally, you may provide a list of parameters (CandidateOption).

**Parameters** should be provided with a Name. You may provide a Type, this will validate if the value provided next to the parameter match the type or not. Several types are supported right now, see below. If no Type is specified, then it will be a boolean flag, which means that it cannot receive any value. If the property is present, value is true, else, false. Finally, you may add a modifier as a binary flag (that means that you hav to provide this values separated by a `|`).

### Parameter Types

- `None`: No validations will be performed (default)
- `Date`: Must match the pattern YYYY-MM-DD
- `Domain`: Domain name. Ej: some.example.com
- `Email`: Ej: some@example.com
- `Ipv4`: Ej: 192.168.0.12
- `Ipv6`: Ej: 2001:0db8:85a3:0000:0000:8a2e:0370:7334
- `Number`: Only integer numbers. Ej: 14, 43, 22, 17
- `Phone`: Phone number. May start with +. Ej: +34 612345678
- `Text`: Not empty text
- `Time`: Must match the pattern HH:mm
- `Url`: Url, including schema (http/https), hostname, path and params
- `UUID`: UUID version 4

### Modifiers

- `DEFAULT`: If this modifier is set, you don't need to type the name of the parameter, you only have to write a value whitout paramter name and it will be automatically binded. You only are able to set one parameter with this flag.
- `REQUIRED`: If this flag is set, the parameter must be supplied, in other case, an error will be prompted and the command will fail.

## Response

```go
type TerminalResponse struct {
  Command  string                // The command executed by Gocli
  Options  map[string]string     // Options that follow the command (validated)
  RawInput string                // The user input without validations neither splits
  Type     TerminalResponseType  // It tells you what happened, see below
  Error    error                 // Nil or the error ocurred
}
```

### TerminalResponseType

- `Cmd`: The command was successfully executed
- `OsCmd`: The command was executed by the OS terminal
- `CmdError`: Error validating the command
- `ParamError`: Error validating some parameter
- `ExecutionError`: Internal error, should not happen
