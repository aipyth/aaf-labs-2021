package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aipyth/aaf-labs-2021/ruban_fi-91_zhytkevych_fi-91/domain"
)

type TokenType string

const TokenTypeNil TokenType = ""
const TokenTypeKeyword TokenType = "keyword"
const TokenTypeIdentificator TokenType = "identificator"
const TokenTypeValue TokenType = "value"

type Token struct {
    Raw string
    Type TokenType
}

var TokenIdentificatorRegexp = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]*$")

type CommandType string

const CommandTypeCreate     CommandType = "create"
const CommandTypeInsert     CommandType = "insert"
const CommandTypeSearch     CommandType = "search"
const CommandTypeQuit       CommandType = "quit"
const CommandTypePrintIndex CommandType = "print_index"

var TokenKeywords = []string{
    string(CommandTypeCreate),
    string(CommandTypeInsert),
    string(CommandTypeSearch),
    string(CommandTypeQuit),
    string(CommandTypePrintIndex),
    "where",
}

var emptySymbols = []rune{
    ' ',
    '\n',
    '\t',
    '\r',
}

type Command struct {
    Type CommandType
    Identificator *Token
    InsertDocument string
    SearchQuery *domain.SearchQuery
}

func NewToken(content string, t TokenType) *Token {
    return &Token{
        Raw:  content,
        Type: t,
    }    
}

func (t *Token) String() string {
    return fmt.Sprintf("{Token: %v|%v}", t.Type, t.Raw)
}

func isTokenKeyword(s string) bool {
    for _, c := range TokenKeywords {
        if s == c {
            return true
        }
    }
    return false
}

func isEmpty(s rune) bool {
    for _, c := range emptySymbols {
        if s == c {
            return true
        }
    }
    return false
}

func split(rawString string) []string {
    contents := make([]string, 0)
    wordStart := 0
    skipWhites := false
    for i, c := range rawString {
        // use skipWhites flag variable to skip whitespaces in value tokens
        if c == '"' {
            skipWhites = !skipWhites
        }
        if skipWhites && i != len(rawString)-1 {
            continue
        }

        if isEmpty(c) && wordStart == i {
            wordStart = i + 1
        } else if isEmpty(c) {
            contents = append(contents, rawString[wordStart:i])
            wordStart = i + 1
        } else if i == len(rawString)-1{
            contents = append(contents, rawString[wordStart:])
        } else {
           continue 
        }
    }
    return contents
}

func tokenizeCommand(s string) []*Token {
    contents := split(s) 
    tokens := make([]*Token, 0)
    for _, c := range contents {
        // check whether the token is keyword
        // the token is keyword while the last one wasn't
        if isTokenKeyword(strings.ToLower(c)) && (len(tokens) == 0 ||
            (len(tokens) != 0 && tokens[len(tokens)-1].Type != TokenTypeKeyword)) {
            tokens = append(tokens, &Token{
                Raw: strings.ToLower(c),
                Type: TokenTypeKeyword,
            })
        // the token is not first and the previous was a keyword
        } else if len(tokens) != 0 && tokens[len(tokens)-1].Type == TokenTypeKeyword &&
            tokens[len(tokens)-1].Raw != "where" {
            tokens = append(tokens, &Token{
                Raw: c,
                Type: TokenTypeIdentificator,
            })
        // everything else is supposed to be a value
        } else {
            tokens = append(tokens, &Token{
                Raw: c,
                Type: TokenTypeValue,
            })
        }
    }
    return tokens
}

// areAllIdentificatorsValid checks all identificators to the predefined regexp
func areAllIdentificatorsValid(tokens []*Token) bool {
    for _, token := range tokens {
        if token.Type == TokenTypeIdentificator &&
            !TokenIdentificatorRegexp.MatchString(token.Raw) {
            return false
        }
    }
    return true
}

func analyzeCreate(tokens []*Token, command *Command) error {
    // there must be 2 tokens
    if len(tokens) != 2 {
        return errors.New("wrong number of tokens")
    }
    // Second token must be identificator.
    switch tokens[1].Type {
    case TokenTypeKeyword:
    case TokenTypeValue:
    case TokenTypeNil:
        return errors.New("no identificator")
    case TokenTypeIdentificator:
        command.Identificator = tokens[1] 
    }
    return nil
}

func analyzeInsert(tokens []*Token, command *Command) error {
    // there must be 3 tokens
    if len(tokens) != 3 {
        return errors.New("wrong number of tokens")
    }
    // Second token must be identificator
    switch tokens[1].Type {
    case TokenTypeKeyword:
    case TokenTypeValue:
    case TokenTypeNil:
        return errors.New("no identificator")
    case TokenTypeIdentificator:
        command.Identificator = tokens[1]
    }
    // Third parameter must be value
    switch tokens[2].Type {
    case TokenTypeKeyword:
    case TokenTypeIdentificator:
    case TokenTypeNil:
        return errors.New("no value")
    case TokenTypeValue:
        doc := tokens[2].Raw
        if doc[0] != '"' || doc[len(doc)-1] != '"' {
            return errors.New("value is not enclosed in double quotation marks")
        }
        command.InsertDocument = doc[1:len(doc)-1]
    }
    return nil 
}

// WARNING! do not look at realization of this function, that is awful!
// better to look at those raccoons https://www.youtube.com/watch?v=6Sq2uRVYmE4
func analyzeSearch(tokens []*Token, command *Command) error {
    command.SearchQuery = new(domain.SearchQuery)

    if len(tokens) < 2 {
        return errors.New("wrong number of tokens")
    }

    switch tokens[1].Type {
    case TokenTypeKeyword:
    case TokenTypeValue:
    case TokenTypeNil:
        return errors.New("no identificator")
    case TokenTypeIdentificator:
        command.Identificator = tokens[1]
    }

    // `SEARCH collection WHERE;` is not acceptable
    if len(tokens) == 3 {
        return errors.New("wrong number of tokens")
    }

    // analyze `WHERE ...;` part
    if len(tokens) == 4 || len(tokens) == 6 {
        if tokens[2].Type != TokenTypeKeyword ||
            tokens[2].Raw != "where" {
            return errors.New("query keyword invalid")
        }
        // next up is query: ("keyword" | "prefix"* | "k1" <n> "k2")
        if len(tokens) == 4 { // "keyword" and "prefix"*
            value := tokens[3].Raw
            if value[0] != '"' {
                return errors.New("value is not enclosed in double quotation marks")
            }
            switch value[len(value)-1] {
            case '*':
                if value[len(value)-2] != '"' {
                    return errors.New("invalid value")
                }
                command.SearchQuery.Prefix = value[1:len(value)-2]
            case '"':
                command.SearchQuery.Keyword = value[1:len(value)-1]
            default:
                return errors.New("invalid value")
            }
        } else if len(tokens) == 6 { // "k1" <n> "k2"
            var value string
            // "k1"
            value = tokens[3].Raw
            if value[0] != '"' || value[len(value)-1] != '"' {
                return errors.New("invalid value")
            }
            command.SearchQuery.Keyword = value[1:len(value)-1]
            // <n>
            value = tokens[4].Raw
            if value[0] != '<' || value[len(value)-1] != '>' {
                return errors.New("invalid value")
            }
            num, err := strconv.Atoi(value[1:len(value)-1])
            if err != nil || num < 0 {
                return errors.New("invalid number as value")
            }
            command.SearchQuery.N = uint(num)
            // "k2"
            value = tokens[5].Raw
            if value[0] != '"' || value[len(value)-1] != '"' {
                return errors.New("invalid value")
            }
            command.SearchQuery.KeywordE = value[1:len(value)-1]
        }
    } else if len(tokens) >= 5 {
        return errors.New("wrong number of tokens")
    }
    return nil
}

func analyzeTokens(tokens []*Token) (*Command, error) {
    if len(tokens) < 1 {
        return nil, errors.New("unsufficient number of tokens")
    }

    command := new(Command)
    // Possible first tokens:
    // (create | insert | search)
    switch tokens[0].Type {
    case TokenTypeIdentificator:
    case TokenTypeValue:
    case TokenTypeNil:
        return nil, errors.New("command not identified")
    case TokenTypeKeyword:
       command.Type = CommandType(tokens[0].Raw) 
    }

    switch command.Type {
    case CommandTypeCreate:
        if err := analyzeCreate(tokens, command); err != nil {
            return nil, err
        }
    case CommandTypeInsert:
        if err := analyzeInsert(tokens, command); err != nil {
            return nil, err
        }
    case CommandTypeSearch:
        if err := analyzeSearch(tokens, command); err != nil {
            return nil, err
        }
    case CommandTypePrintIndex:
    case CommandTypeQuit:
        break
    default:
        return nil, errors.New("unsupported command type")
    }

    return command, nil
}

func NewCommand(input string) (*Command, error) {
    tokens := tokenizeCommand(input)
    if !areAllIdentificatorsValid(tokens) {
        return nil, errors.New("invalid identificators")
    }
    command, err := analyzeTokens(tokens)
    if err != nil {
        return nil, err
    }
    return command, nil
}

