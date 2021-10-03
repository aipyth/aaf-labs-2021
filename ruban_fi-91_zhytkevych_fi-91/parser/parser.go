package parser

import (
	"errors"
	"strings"
)

type Command struct {
	cmdName  string
	args     []string
	extraCmd *Command
}

func createParser(args []string) (*Command, error) {
	if len(args) > 1 {
		return _, errors.New("Error, invalid CREATE command input, use: `CREATE collection_name`")
	}
	return &Command{
		cmdName:  "CREATE",
		args:     args,
		extraCmd: nil,
	}, nil
}

func insertParser(args []string) (*Command, error) {
	for i := range args[1:] {
		args[i] = strings.Split(v, "\"")
		if len(args[i]) != 3 || args[1] == "" {
			return errors.New("Error, invalid document name")
		}
	}
	return &Command{
		cmdName:  "INSERT",
		args:     args,
		extraCmd: nil,
	}, nil
}

func searchParser(args []string) (*Command, error) {

	for i := range args[1:] {
		args[i] = strings.Split(v, "\"")
		if len(args[i]) != 3 || args[1] == "" {
			return errors.New("Error, invalid document name")
		}
	}
	return &Command{
		cmdName:  "INSERT",
		args:     args,
		extraCmd: nil,
	}, nil
}

func selectParser(name string) func(args []string) (*Command, error) {
	switch name {
	case "CREATE":
		return createParser
	case "INSERT":
		return insertParser
	case "SEARCH":
		return searchParser
	}
	return error.New("Invalid command " + name)
}

func deleteWhiteSpaces(arr []string) []string {
	for i := range arr {
		if arr[i] == " " {
			new_arr = (make([]string), arr[:i]...)
			new_arr = append(new_arr, arr[i+1:]...)
			arr[i] = new_arr
		}
	}
	return new_arr
}

func Parse(command string) *Command {
	splitted_string := strings.Split(command, " ")
	splitted_string = deleteWhiteSpaces(splitted_string)
	commandName := splitted_string[0]

}
