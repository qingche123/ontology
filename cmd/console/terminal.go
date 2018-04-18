/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package console

import (
	"fmt"
	colorable "gx/ipfs/QmdvecVcFhbo5x4f3arqmfxyE3NzqwWyp77KzA68EKXJeX/go-colorable"
	"io"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/peterh/liner"
	sdk "github.com/ontio/ontology-go-sdk"
)

var (
	passwordRegexp = regexp.MustCompile(`personal.[nus]`)
	onlyWhitespace = regexp.MustCompile(`^\s*$`)
	exit           = regexp.MustCompile(`^\s*exit\s*;*\s*$`)
)

// DefaultPrompt is the default prompt line prefix to use for user input querying.
const DefaultPrompt = "ontology > "

// Config is the collection of configurations to fine tune the behavior of the
// console.
type Config struct {
	OntSDK   *sdk.OntologySdk   // Ontology go-sdk
	Prompt   string       // Input prompt prefix string (defaults to DefaultPrompt)
	Prompter UserPrompter // Input prompter to allow interactive user feedback (defaults to TerminalPrompter)
	Printer  io.Writer    // Output writer to serialize any display strings to (defaults to os.Stdout)
}

// Console is a JavaScript interpreted runtime environment. It is a fully fleged
// console attached to a running node via an external or in-process RPC
// client.
type Console struct {
	OntSDK   *sdk.OntologySdk   // Ontology go-sdk
	prompt   string       // Input prompt prefix string
	prompter UserPrompter // Input prompter to allow interactive user feedback
	printer  io.Writer    // Output writer to serialize any display strings to
}

func New(ontSdk *sdk.OntologySdk) (*Console, error) {
	// Initialize the console and return
	console := &Console{
		OntSDK:   ontSdk,
		prompt:   DefaultPrompt,
		prompter: Stdin,
		printer:  colorable.NewColorableStdout(),
	}
	if err := console.Init(); err != nil {
		return nil, err
	}

	return console, nil
}

// init retrieves the available APIs from the remote RPC provider and initializes
// the console's namespaces based on the exposed modules.
func (c *Console) Init() error {
	process = NewProcess(c)

	if c.prompter != nil {
		c.prompter.SetWordCompleter(c.AutoCompleteInput)
	}
	return nil
}

// AutoCompleteInput is a pre-assembled word completer to be used by the user
// input prompter to provide hints to the user about the methods available.
func (c *Console) AutoCompleteInput(line string, pos int) (string, []string, string) {
	if len(line) == 0 || pos == 0 {
		return "", nil, ""
	}
	// Chunk data to relevant part for autoCompletion
	start := pos - 1
	for ; start > 0; start-- {
		// Skip all methods and namespaces (i.e. including the dot)
		if line[start] == '.' || (line[start] >= 'a' && line[start] <= 'z') || (line[start] >= 'A' && line[start] <= 'Z') {
			continue
		}
		// Handle wallet in a special way (i.e. other numbers aren't auto completed)
		if start >= 5 && line[start-5:start] == "wallet" {
			start -= 5
			continue
		}
		// We've hit an unexpected character, autoComplete form here
		start++
		break
	}
	return line[:start], CompleteKeywords(line[start:pos]), line[pos:]
}

// Welcome show summary of current Geth instance and some metadata about the
// console's available modules.
func (c *Console) Welcome() {
	welcome := "\n    Welcome to the Ontology terminal! If you want to leave it, Please enter <exit>.\n\n" +
				"    Modules:\n\n" +
				"    [wallet]    [rpc]   [contract]\n "
	fmt.Fprintln(c.printer, welcome)
}

// Evaluate executes code and pretty prints the result to the specified output
// stream.
func (c *Console) Evaluate(statement string) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(c.printer, "[native] error: %v\n", r)
		}
	}()
	return Call(statement, c)
}

// Interactive starts an interactive user session, where input is propted from
// the configured user prompter.
func (c *Console) Interactive() {
	var (
		prompt    = c.prompt          // Current prompt line (used for multi-line inputs)
		indents   = 0                 // Current number of input indents (used for multi-line inputs)
		input     = ""                // Current user input
		scheduler = make(chan string) // Channel to send the next prompt on and receive the input
	)
	// Start a goroutine to listen for promt requests and send back inputs
	go func() {
		for {
			// Read the next user input
			line, err := c.prompter.PromptInput(<-scheduler)
			if err != nil {
				// In case of an error, either clear the prompt or fail
				if err == liner.ErrPromptAborted { // ctrl-C
					prompt, indents, input = c.prompt, 0, ""
					scheduler <- ""
					continue
				}
				close(scheduler)
				return
			}
			// User input retrieved, send for interpretation and loop
			scheduler <- line
		}
	}()
	// Monitor Ctrl-C too in case the input is empty and we need to bail
	abort := make(chan os.Signal, 1)
	signal.Notify(abort, syscall.SIGINT, syscall.SIGTERM)

	// Start sending prompts to the user and reading back inputs
	for {
		// Send the next prompt, triggering an input read and process the result
		scheduler <- prompt
		select {
		case <-abort:
			fmt.Fprintln(c.printer, "caught interrupt, exiting")
			return

		case line, ok := <-scheduler:
			// User input was returned by the prompter, handle special cases
			if !ok || (indents <= 0 && exit.MatchString(line)) {
				return
			}
			if onlyWhitespace.MatchString(line) {
				continue
			}
			// Append the line to the input and check for multi-line interpretation
			input += line + "\n"

			indents = countIndents(input)
			if indents <= 0 {
				prompt = c.prompt
			} else {
				prompt = strings.Repeat(".", indents*3) + " "
			}
			// If all the needed lines are present, save the command and run
			if indents <= 0 {
				c.Evaluate(input)
				input = ""
			}
		}
	}
}

// countIndents returns the number of identations for the given input.
// In case of invalid input such as var a = } the result can be negative.
func countIndents(input string) int {
	var (
		indents     = 0
		inString    = false
		strOpenChar = ' '   // keep track of the string open char to allow var str = "I'm ....";
		charEscaped = false // keep track if the previous char was the '\' char, allow var str = "abc\"def";
	)

	for _, c := range input {
		switch c {
		case '\\':
			// indicate next char as escaped when in string and previous char isn't escaping this backslash
			if !charEscaped && inString {
				charEscaped = true
			}
		case '\'', '"':
			if inString && !charEscaped && strOpenChar == c { // end string
				inString = false
			} else if !inString && !charEscaped { // begin string
				inString = true
				strOpenChar = c
			}
			charEscaped = false
		case '{', '(':
			if !inString { // ignore brackets when in string, allow var str = "a{"; without indenting
				indents++
			}
			charEscaped = false
		case '}', ')':
			if !inString {
				indents--
			}
			charEscaped = false
		default:
			charEscaped = false
		}
	}
	return indents
}

func (c *Console) Stop(graceful bool) error {
	return nil
}
