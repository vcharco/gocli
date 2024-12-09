# Gocli

This project consists of a Go module aimed at simplifying the development of applications based on interactive command-line interfaces. The module provides a function that opens a CLI which returns the entered command as a result.
The CLI includes features such as autocompletion, suggestions, command history, and shortcuts like `CTRL+C`, `CTRL+V` (see the complete list in the Special Commands section at the end of the doc).

## Learn by example

### Initial config

Here is a basic configuration, but covering most of the available options. We will define commands, options for those commands and some modifiers to make a param required or default. Also we configure the prompt, bypass character and control keys to be triggered.

```go
package main

import (
 "fmt"
 "strconv"

 gc "github.com/vcharco/gocli"
)

func main() {

  // Here we declare the commands and their params
  commands := []gc.Command {
    // We may set a description to commands and params for the help
    {Name: "foo", Description: "Perform foo operations", Params: []gc.Param {
      // If no type defined, it will be a flag (Ej: -f, --help, --get-history)
      {Name: "-f", Description: "foo flag"},
      // We set types for param validation (Ej: --foo foobar, --limit 20, -l 20)
      {Name: "--foo", Type: gc.Text},
      // This is how we specify a default value and set it to required
      {Name: "fooDefault", Type: gc.Number, Modifier: gc.DEFAULT | gc.REQUIRED},
    }},
    // This is how we set a hidden command. This is a valid command with autcompletion
    // as the rest of the params, but will not be displayed in the suggestions when tab
    {Name: "exit", Hidden: true},
  }

  // Configuration
  cli := gc.Terminal {
    Styles:          TerminalStyles{}, // Default styles
    Commands:        commands,
    BypassCharacter: ":", // Allows to execute commands by the OS -> :ls -l
    CtrlKeys:        []byte{gc.Ctrl_A, gc.Ctrl_B}, // CRTL keys to caputure
  }

  for {
    // Gets the user input
    response := cli.Get()

    // Also we may pass a text what will be appended after the prompt
    // response := cli.Get("text sent to the CLI")

    // Here you handle the user input
  }
}
```

### Customizing the terminal

Styling is optional, as default values are applied. But you will get a better user experience changing this default values. Note that we use `Bg` prefixed color names for `BackgroundColor` and `SelBackgroundColor`.

```go
// First, define a new TerminalStyles objects and fill only what you want
styles := gc.TerminalStyles {
  Prompt:                 "Gocli> "           // Text to be displayed in the prompt
  PromptColor:            gc.Blue,            // Color of the prompt
  Cursor:                 gc.CursorBlock,     // CursorBlock, CursorBar or CursorUnderline
  ForegroundColor:        gc.White,           // Color of the text when not selected
  BackgroundColor:        gu.BgTransparent,   // Background color of the terminal
  SelForegroundColor:     gc.Blue,            // Color of the text when selected
  SelBackgroundColor:     gc.BgLightBlue,     // Color of the selection
  HelpTextForeground:     gu.LightGray        // Color of the help (?) text
	HelpTitlesForeground:   gu.Blue             // Color of the title sections in the help
	HelpCommandForeground:  gu.White            // Color of the command in the help display
	HelpParamsForeground:   gu.Yellow           // Color of the params in the help display
	HelpRequiredForeground: gu.Red              // Color of (REQUIRED) param flag in help
	HelpLineColor:          gu.Blue             // Color of the line in the help display
}

// Finally, add this configuration to you Terminal object
cli := gc.Terminal {
  Styles:          styles,                // Our custom styles
  Commands:        commands,
  BypassCharacter: ":",
  CtrlKeys:        []byte{gc.Ctrl_A, gc.Ctrl_B},
}
```

### Playing with the response

Here it's an example of how to use the response received.

```go
// First we check the response is a valid command
if response.Type == gc.Cmd {
  // Now we check what command was typed by the user
  switch response.Command {
    case "foo":
      // This is how we retrieve values. First is the name and then the default value
      fooDefault := response.GetParam("fooDefault", "").(string)
      fooVal := response.GetParam("--foo", "foo default value").(string)

      // For flag params, default value must be a boolean. Normally set to false
      fVal := response.GetParam("-f", false).(bool)

      fmt.Printf("default: %v; -f: %v; --foo: %v\n", fooDefault, fVal, fooVal)
  }
}
```

### Command history

The cli has a defautl command history. We use the UP/DOWN arrow keys to get the previuos command or the next command in the history as in any other cli.

```go
// This is how we print the history
cli.PrintHistory(20)

// This is how we clear the history
cli.ClearHistory()

// This is how we get the number of commands in the history
numCmds := cli.CountHistory()

// This is how we get all commands in history
historyCommands := cli.GetHistory()

// This is how we get an specific command in history
historyCommand := cli.GetHistoryAt(5)
```

### Checking response errors

There are several kind of errors, but we may trigger all of them by checking the value of the `Error` attribute. Then, we may check the type of error.

```go
if response.Error != nil {
  switch response.Type {
    // The command doesn't match to any of declared in the Options
    case gc.CmdError:
      fmt.Printf("Invalid command: %v\n", response.Error.Error())
    // Some parameter has an invalid value or the param doesn't exist
    case gc.ParamError:
      fmt.Printf("Invalid parameters: %v\n", response.Error.Error())
    // This is an internal error, should never happen
    case gc.ExecutionError:
      fmt.Printf("Internal error: %v\n", response.Error.Error())
   }
}
```

### Handle CTRL+Key Combinations

```go
// First, configure what keys you want to trigger
cli := gc.Terminal {
  // ...
  CtrlKeys:     []byte{gc.Ctrl_A, gc.Ctrl_B},
}

// Get the response
response := cli.Get()

// Now check what CTRL+Key combination was triggered
if response.Type == gc.CtrlKey {
  switch response.CtrlKey {
  case gc.Ctrl_a:
    fmt.Println("Captured CTRL+A")
  case gc.Ctrl_b:
    fmt.Println("Captured CTRL+B")
  }
}
```

### Commands bypassed to the OS

Commands with response type `OsCmd` were executed by the console of the operative system. Normally, we don't need to do anything with this responses, as its main pupose is just execute a command by the OS console, but we still perform some actions. In this case, we only have the `RawInput` attribute available.

```go
if response.Type == gc.OsCmd {
  userInput := response.RawInput
  fmt.Printf("Comman executed by the OS: %v\n", userInput)
}
```

## Types, values and other usefull information

### Configuration

- `Prompt`: This is the text at the beggining of the line.
- `PromptColor`: Set the color of the prompt.
- `Commands`: This list of commands is used for autocompletion and suggestions. It contains a sublist of valid parameters for each command.
- `BypassCharacter`: Gocli checks if the input starts with this character, and in that case, instead of processing it, it sends it directly to the operating system's console. This allows you to execute OS commands without leaving Gocli.
  - Example for BypassCharacter `:`: `Prompt> :ls -l`
- `CtrlKeys`: A list of CTRL+Key combinations you want to override. When one of these combinations is detected, gocli will respond with the Type `CtrlKey` and the value of the detected combination will be available in the reponse property `CtrlKey`.

### Commands

They are the commands available for your custom cli. Each command must be provided with a Name. Optionally, you may provide a list of parameters. If you set the Hidden attribute, the command still be valid, but won't be displayed in the suggestions when pressing tabulator.

**Parameters** should be provided with a Name. You may provide a Type, this will validate if the value provided next to the parameter match the type or not. Several types are supported right now, see below. For Number or FloatNumber types, the casting is performed automatically. As the returned type is an interface{} type, you must make the assertions `.(string)`, `.(int)`, `.(float64)` or `.(bool)` If no Type is specified, then it will be a boolean flag, which means that it cannot receive any value. If the property is present, value is true, else, false. Finally, you may add a modifier as a binary flag (that means that you hav to provide this values separated by a `|`).

### Parameter Types

- `None`: No validations will be performed (default)
- `Date`: Must match the pattern YYYY-MM-DD
- `Domain`: Domain name. Ej: some.example.com
- `Email`: Ej: some@example.com
- `Ipv4`: Ej: 192.168.0.12
- `Ipv6`: Ej: 2001:0db8:85a3:0000:0000:8a2e:0370:7334
- `Number`: Only integer numbers. Ej: 14, 43, 22, 17
- `FloatNumber`: Float numbers. Ej: 14.3, 43.234, 22.0
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
  Param  map[string]interface{}  // Parameters that follow the command (validated)
  RawInput string                // The user input without validations neither splits
  Type     TerminalResponseType  // It tells you what happened, see below
  CtrlKey  byte                  // If Type = CtrlKey, this is the CTRL+key combination
  Error    error                 // Nil or the error ocurred
}
```

### TerminalResponseType

- `Cmd`: The command was successfully executed
- `OsCmd`: The command was executed by the OS terminal
- `CmdHelp`: User has printed the help
- `CtrlKey`: A registered Ctrl key has been pressed
- `CmdError`: Error validating the command
- `ParamError`: Error validating some parameter
- `ExecutionError`: Internal error, should not happen

## Special commands and characters

There are several shortcuts listed down here to make more fluid your interaction (you may override this shortcuts, but is not recommended in most cases).

- `CTRL+C`: Copy the selected text to the clipboard. If no text is selected, it copies the full line.
- `CTRL+V`: Paste the clipboard.
- `CTRL+L`: Clear the screen.
- `CTRL+X`: Exit the program (safely)
- `CTRL+A`: Move the cursor at the beginning of the line
- `CTRL+E`: Move the cursor at the end of the line

There are two special characters.

- `BypassCharacter`: This character must be declared in order to bypass commands to the OS terminal. Let's say we set this charcter to `:`.
  - `CLI> :ls -l`
  - `CLI> :whoami`
  - `CLI> :grep root /etc/passwd`
- `?`: This is a special character for displaying help. Its usage is simple, type this character after a command or while you type the command and then press Enter. If this character is detected at the end of the input, the terminal will recognize the command and will display all available information like required and non required params, flags, default params and, of course, all the descriptions provided when we declared the Options.
  - `CLI> print-history?`: It displays the help for the command `print-history`.
  - `CLI> print-hi?` : We don't need to end the command if there are no conflicts with other commands.
  - `CLI> print-history 20 ?`: We may display the command help even if we have already type parameters.
