package jsonwrap

import (
	"fmt"
	"regexp"
	"strings"
)

// eg: `error found in #1 byte of ...|\"120\"|..., bigger context ...|\"120\"|..., error found in #10 byte of ...|on\": \"120\"``.
// NOTE: '(?s)' permits to include new line characters.
var rxJSONiterParsingMsg = regexp.MustCompile(`byte of \.\.\.\|((?s).+)\|\.\.\.`)

// eg: `fairjungle.WebhookReq.Data: fairjungle.WebhookPatch.Events: ReadObject: found unknown field: caca, `.
var rxJSONiterUnknowFieldMsg = regexp.MustCompile(`found unknown field: (.+), `)

// Error is a json parsing error.
type Error struct {
	Kind ErrorKind
	msg  string
	err  error
}

// ErrorKind represent a json parsing error kind.
type ErrorKind string

// Possible error kinds.
const (
	ErrorDecodingFailed  ErrorKind = "decodingFailed"
	ErrorParsingFailed   ErrorKind = "parsingFailed"
	ErrorUnexpectedField ErrorKind = "unexpectedField"
)

// Error implements error interface.
func (e Error) Error() string {
	return e.msg
}

// Unwrap returns wrapped error.
func (e Error) Unwrap() error {
	return e.err
}

func handleDecodingError(err error) error {
	if err != nil {
		parts := strings.Split(err.Error(), "error found in")

		// detect jsoniter 'unknown field' error
		if len(parts) > 0 {
			matches := rxJSONiterUnknowFieldMsg.FindStringSubmatch(parts[0])
			if len(matches) > 1 {
				return Error{
					Kind: ErrorUnexpectedField,
					msg:  fmt.Sprintf("unexpected field '%s'", matches[1]),
					err:  err,
				}
			}
		}

		// detect jsoniter 'parsing' error message
		if len(parts) > 1 {
			subparts := strings.Split(parts[1], "bigger context")
			matches := rxJSONiterParsingMsg.FindStringSubmatch(subparts[0])
			if len(matches) > 1 {
				return Error{
					Kind: ErrorParsingFailed,
					msg:  fmt.Sprintf("parsing failed at '%s'", matches[1]),
					err:  err,
				}
			}
		}

		return Error{
			Kind: ErrorDecodingFailed,
			msg:  "failed to decode json",
			err:  err,
		}
	}

	return err
}
