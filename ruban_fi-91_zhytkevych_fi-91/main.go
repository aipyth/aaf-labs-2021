package main

import (
	"bufio"
	"fmt"
	"log"
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
        documents := dom.Search(command.Identificator.Raw, *command.SearchQuery)
        for _, doc := range documents {
            if doc.Collection.Name == command.Identificator.Raw {
                // os.Stdout.WriteString(doc.String())
                fmt.Println(doc)
            }
        }
    case CommandTypePrintIndex:
        os.Stdout.WriteString(dom.IndexerRepresentationString() + "\n")
    case CommandTypeQuit:
        os.Exit(0)
    default:
        os.Stderr.WriteString("[ERROR]: unknown command\n")
    }

}

func pairedQuotationMarks(s string) (paired bool) {
    paired = true
    for _, c := range s {
        if c == '"' { paired = !paired }
    }
    return
}

func splitIntoCmd(payload string) []string {
    cmds := make([]string, 0, 1)
    lastCmdStart := 0
    inStrIdentifier := false
    for i, c := range payload {
        switch c {
        case '"':
            inStrIdentifier = !inStrIdentifier
        case ';':
            if inStrIdentifier { continue }
            cmds = append(cmds, payload[lastCmdStart:i])
            lastCmdStart = i + 1
        }
    }
    return cmds
}

func main() {
    rbuff := bufio.NewReader(os.Stdin)

    os.Stdout.WriteString(`# DDDB - Documents database with full-text search
Available commands:
    QUIT;
    CREATE collection_name;
    INSERT collection_name “value”;
    PRINT_INDEX;
    SEARCH collection_name [WHERE query];
        query := “keyword” 
 		 | “prefix”*
 		 | “keyword_1” <N> “keyword_2”

`)

    var payload string
    var command *Command
    for {
        os.Stdout.WriteString(">")
        s, err := rbuff.ReadString('\n')
        if err != nil {
            os.Stderr.WriteString("[ERROR]: " + err.Error() + "\n")
            rbuff.Reset(os.Stdin)
            continue
        }

        trimmed := strings.TrimSpace(s)
        if len(trimmed) == 0 {
            continue
        }
        if !pairedQuotationMarks(trimmed) {
            continue
        }

        if trimmed[len(trimmed)-1] == ';' {
            payload += trimmed
            cmds := splitIntoCmd(payload)
            log.Println("cmds: ", cmds)
            for _, cmd := range cmds {
                cmd := strings.TrimSpace(cmd)
                if cmd == "" { continue }
                command, err = NewCommand(cmd)
                if err != nil {
                    os.Stderr.WriteString("[ERROR]: " + err.Error() + "\n")
                } else {
                    executeCommand(command)
                }
            }
            payload = ""
        } else {
            payload += s
        }
    }
}
