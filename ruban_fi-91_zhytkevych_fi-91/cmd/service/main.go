package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aipyth/aaf-labs-2021/ruban_fi-91_zhytkevych_fi-91/domain"
)

var dom = domain.NewDomain()

func executeCommand(command *Command) {
    var err error
    switch command.Type {
    case CommandTypeCreate:
        err = dom.CreateCollection(command.Identificator.Raw)
        if err != nil {
            os.Stderr.WriteString("[ERROR]:" + err.Error() + "\n")
            return
        }
        os.Stdout.WriteString("Collection " + command.Identificator.Raw + " created!\n")
    case CommandTypeInsert:
        err = dom.InsertDocument(command.Identificator.Raw, command.InsertDocument)
        if err != nil {
            os.Stderr.WriteString("[ERROR]:" + err.Error()+ "\n")
            return
        }
        os.Stdout.WriteString("Document added to " + command.Identificator.Raw + ".\n")
    case CommandTypeSearch:
        documents := dom.Search(*command.SearchQuery)
        for _, doc := range documents {
            // os.Stdout.WriteString(doc.String())
            fmt.Println(doc)
        }
    default:
        os.Stderr.WriteString("[ERROR]: unknown command\n")
    }

}

func main() {
    rbuff := bufio.NewReader(os.Stdin)

    var payload string
    var command *Command
    for {
        os.Stdout.WriteString(">")
        s, err := rbuff.ReadString(';')
        if err != nil {
            os.Stderr.WriteString("[ERROR]: " + err.Error() + "\n")
            rbuff.Reset(os.Stdin)
            continue
        }
        s = strings.TrimSpace(s)
        payload += s
        if payload[len(payload)-1] == ';' {
            cmds := strings.Split(payload, ";")
            for _, cmd := range cmds {
                if cmd == "" { continue }
                command, err = NewCommand(cmd)
                if err != nil {
                    os.Stderr.WriteString("[ERROR]: " + err.Error() + "\n")
                } else {
                    executeCommand(command)
                }
            }
            payload = ""
        }
    }
}
