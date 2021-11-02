package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/aipyth/aaf-labs-2021/ruban_fi-91_zhytkevych_fi-91/domain"
)

var dom = domain.NewDomain()

func executeCommand(command *Command) {
    var err error
    switch command.Type {
    case CommandTypeCreate:
        err = dom.CreateCollection(command.Identificator.Raw)
        if err != nil {
            os.Stderr.WriteString("[ERROR]:" + err.Error())
            return
        }
        os.Stdout.WriteString("Collection " + command.Identificator.Raw + " created!")
    case CommandTypeInsert:
        err = dom.InsertDocument(command.Identificator.Raw, command.InsertDocument)
        if err != nil {
            os.Stderr.WriteString("[ERROR]:" + err.Error())
            return
        }
        os.Stdout.WriteString("Document added to " + command.Identificator.Raw + ".")
    case CommandTypeSearch:
        documents := dom.Search(*command.SearchQuery)
        for _, doc := range documents {
            // os.Stdout.WriteString(doc.String())
            fmt.Println(doc)
        }
    default:
        os.Stderr.WriteString("[ERROR]: unknown command")
    }

}

func main() {
    rbuff := bufio.NewReader(os.Stdin)

    var command *Command
    for {
        s, err := rbuff.ReadString(';')
        if err != nil {
            os.Stderr.WriteString("[ERROR]: " + err.Error())
            rbuff.Reset(os.Stdin)
            continue
        }

        command, err = NewCommand(s)
        if err != nil {
            os.Stderr.WriteString("[ERROR]: " + err.Error())
            continue
        }

        executeCommand(command)
    }
}
