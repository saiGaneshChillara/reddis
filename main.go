package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var ARITY = map[string]struct{ Min, Max int} {
	"PING": {1, 2},
	"ECHO": {2, 2},
	"COMMAND": {1, 2},
}

func checkArity(cmd string, args []string) (bool, string) {
	if arity, ok := ARITY[cmd]; ok {
		if len(args) < arity.Min || len(args) > arity.Max {
			return false, fmt.Sprintf("ERR wrong number of arguments for '%s' command", cmd)
		}
	}
	return true, ""
}

func eb(s string, ok bool) string {
	if !ok { return "$-1\r\n" }
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}
func es(s string) string { return fmt.Sprintf("+%s\r\n", s) }
func ee(m string) string { return fmt.Sprintf("-%s\r\n", m) }

func handle(args []string) string {
	cmd := strings.ToUpper(args[0])
	switch cmd {
	case "PING":
		// TODO: If len(args) > 2, return ee("ERR wrong number of arguments for 'PING' command")
		if ok, err := checkArity(cmd, args); !ok {
			return ee(err)
		}
		if len(args) == 1 { return es("PONG") }
		return eb(args[1], true)
	case "ECHO":
		// TODO: If len(args) != 2, return ee("ERR wrong number of arguments for 'ECHO' command")
		if ok, err := checkArity(cmd, args); !ok {
			return ee(err)
		}
		return eb(args[1], true)
	case "COMMAND":
		if len(args) > 1 && strings.ToUpper(args[1]) == "DOCS" { return es("OK") }
		return es("OK")
	}
	return ee(fmt.Sprintf("ERR unknown command '%s'", args[0]))
}

func parseArgs(line string) []string {
	var args []string
	var cur strings.Builder
	inQ := false
	for _, ch := range line {
		switch {
		case ch == '"' && !inQ: inQ = true
		case ch == '"' && inQ: inQ = false
		case ch == ' ' && !inQ:
			if cur.Len() > 0 { args = append(args, cur.String()); cur.Reset() }
		default: cur.WriteRune(ch)
		}
	}
	if cur.Len() > 0 { args = append(args, cur.String()) }
	return args
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" { continue }
		fmt.Print(handle(parseArgs(line)))
	}
}
