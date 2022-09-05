package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/diff"
)

func format(input string) (string, error) {
	var sb strings.Builder
	var prev, prevNonSpace item
	lex := lex("", input, "", "", true, true, true)
	for current := lex.nextItem(); current.typ != itemEOF; current = lex.nextItem() {
		//sb.WriteString(fmt.Sprintf("[%d:%d:%d]", prevNonSpace.typ, prev.typ, current.typ))
		switch current.typ {
		case itemError:
			return "", errors.New(current.val)
		case itemComment:
			sb.WriteString(current.val)
		case itemLeftDelim:
			sb.WriteString(current.val)
		case itemLeftTrimDelim:
			sb.WriteString(current.val)
		case itemRightTrimDelim:
			sb.WriteString(current.val)
		case itemRightDelim:
			if prev.typ != itemRawString && prev.typ != itemRightTrimDelim {
				sb.WriteRune(' ')
			}
			sb.WriteString(current.val)
		case itemSpace:
			if strings.Contains(current.val, "\n") {
				sb.WriteString(strings.TrimLeft(current.val, " "))
			} else if prev.typ == itemLeftDelim || prev.typ == itemLeftTrimDelim {
				sb.WriteString(current.val)
			}
		case itemRawString:
			sb.WriteString(formatSafe(current.val))
		case itemRightParen: // no space before )
			sb.WriteString(current.val)
		case itemText: // "plain" text
			sb.WriteString(current.val)
		default:
			switch {
			case current.typ == itemField && (prev.typ == itemField || prev.typ == itemVariable || prev.typ == itemRightParen): // no space between .a & .b in .a.b or ).b
			case current.typ == itemChar && current.val == ",": // no space before "," in a, b := range ...
			case prevNonSpace.typ == itemLeftParen: // no space after (
			case prevNonSpace.typ == itemLeftDelim && prev.typ == itemSpace: // see line #39
			case prev.typ == itemLeftTrimDelim: // leftTrimDelim contains a space
			case prev.typ == itemSpace && strings.Contains(prev.val, "\n"):
			default:
				sb.WriteRune(' ')
			}
			//sb.WriteString(fmt.Sprintf("[%v:%s", current.typ, current.val))
			sb.WriteString(current.val)
		}
		prev = current
		if current.typ != itemSpace {
			prevNonSpace = current
		}
	}
	return sb.String(), nil
}

func formatSafe(input string) string {
	formatted, err := format(input)
	if err != nil {
		return input
	}
	return formatted
}

func formatFile(name string) error {
	data, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	formatted, err := format(string(data))
	if err != nil {
		return fmt.Errorf("format %s: %w", name, err)
	}
	return os.WriteFile(name, []byte(formatted), 0644)
}

func checkFile(name string) error {
	data, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	formatted, err := format(string(data))
	if err != nil {
		return fmt.Errorf("format %s: %w", name, err)
	}
	if !bytes.Equal(data, []byte(formatted)) {
		var b bytes.Buffer
		err = diff.Text(fmt.Sprintf("a/%s", name), fmt.Sprintf("b/%s", name), data, formatted, bufio.NewWriter(&b))
		if err != nil {
			return fmt.Errorf("diff %s: %w", name, err)
		}
		println(b.String())
		return fmt.Errorf("would reformat %s", name)
	}
	return nil
}
