/*
Copyright (C) 2026 Andrew Flint.

This file is part of arnhemcr/financial.

Arnhemcr/financial is free software:
you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Arnhemcr/financial is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with arnhemcr/financial.
If not, see <https://www.gnu.org/licenses/>.
*/

/*
Mrglent [filters] financial transactions in [Ledger] entry format from multiple journals into a general journal.

Mrglent reads Ledger journals from standard input.
It extracts dated journal entries.
If an entry's date cannot be parsed according to the layout,
mrglent writes a message to standard error and exits with a non-zero status.
Dated entries marked as mirrors (between "# mirror entry" and "# end mirror entry" comment lines) are discarded.
See this module's program mcsv2lent for more on marked entries.
All other journal content is also discarded including automatic transactions, global comments and command directives.

Mrglent orders the entries by date ascending and writes them to standard output.

Usage:

	mrglent [flags]

The flags are:

	-d string
	  	Go date layout of Ledger journal entries (default "2006-01-02")
	-h	write this help text then exit

See also [this package's README].

[filters]: https://en.wikipedia.org/wiki/Filter_(software)
[Ledger]: https://ledger-cli.org
[this package's README]: https://github.com/arnhemcr/financial/tree/main
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	aft "github.com/arnhemcr/financial/transaction"
	"log"
	"os"
	"slices"
	"time"
	"unicode"
)

func main() {
	log.SetPrefix("mrglent: ")
	log.SetFlags(0)

	dateLayout := parseFlags()
	if !aft.IsDateLayout(dateLayout) {
		log.Fatalf("date layout must be Go-style e.g. %q", time.DateOnly)
	}

	s := bufio.NewScanner(os.Stdin)

	es, err := parseEntries(s, dateLayout)
	if err != nil {
		log.Fatal(err)
	}

	oes := sortEntries(es)
	for _, oe := range oes {
		fmt.Fprint(os.Stdout, oe)
	}
}

// An entry represents a dated Ledger journal entry.
type entry struct {
	Date string
	Text string
}

/*
ParseEntries reads a stream of Ledger journals and returns entries with dates.
Other content is discarded, including dated entries marked as mirrors and Ledger block comments.
If it fails to parse the date of an entry, parseEntries returns the error.

For further information on dated entries (or transactions) and block comments,
see "Transactions and Comments" and "Commenting on your journal" in the [Ledger 3 manual].

[Ledger 3 manual]: https://ledger-cli.org/doc/ledger3.html
*/
func parseEntries(s *bufio.Scanner, dateLayout string) ([]entry, error) {
	var (
		es                            []entry
		e                             entry
		inBlockComment, inMirrorEntry bool
		lnN                           int
	)

	for s.Scan() {
		ln := s.Text() + "\n"
		lnN++

		if inBlock(&inBlockComment, ln, aft.StartBlockComment, aft.EndBlockComment) {
			continue
		}

		if inBlock(&inMirrorEntry, ln, aft.StartMirrorEntry, aft.EndMirrorEntry) {
			continue
		}

		switch {
		case unicode.IsDigit(rune(ln[0])):
			if e.Date != "" {
				es = append(es, e)
			}

			d, err := aft.ParseDate(ln, dateLayout)
			if err != nil {
				return es, fmt.Errorf("line %v: %w", lnN, err)
			}

			// This line starts with a date and is the first line in the next entry.
			e.Date, e.Text = d, ln
		case aft.IsLedgerIndented(ln):
			// This line is indented and belongs to the current entry.
			e.Text += ln
		}
	}

	if e.Date != "" {
		es = append(es, e)
	}

	return es, nil
}

/*
InBlock reports whether the line from a Ledger journal is in a block delimited by start and end lines.
It also updates the in block state.
*/
func inBlock(state *bool, line, startLine, endLine string) bool {
	switch {
	case line == startLine:
		*state = true
	case line == endLine:
		*state = false
	case *state:
		// The line is in-between start and end lines.
	default:
		return false
	}

	return true
}

/*
ParseFlags returns the date layout of Ledger journal entries parsed from command line flags.
If help was requested, parseFlags writes this program's help text then exits.
If the flags are invalid, this program exits with a non-zero status.
*/
func parseFlags() string {
	var dateLayout string

	flag.StringVar(&dateLayout, "d", time.DateOnly, "Go date layout of Ledger journal entries")

	var help bool

	flag.BoolVar(&help, "h", false, "write this help text then exit")

	flag.Usage = usage
	flag.Parse()

	if help {
		usage()
		os.Exit(0)
	}

	return dateLayout
}

// Sort orders the texts of a list of Ledger journal entries by date ascending.
func sortEntries(es []entry) []string {
	d2txts := make(map[string][]string)

	var ds []string

	for _, e := range es {
		d := e.Date

		_, found := d2txts[d]
		if !found {
			d2txts[d] = []string{}

			ds = append(ds, d)
		}

		d2txts[d] = append(d2txts[d], e.Text)
	}

	var oes []string

	slices.Sort(ds)

	for _, d := range ds {
		oes = append(oes, d2txts[d]...)
	}

	return oes
}

// Usage writes the help text for this program.
func usage() {
	fmt.Fprint(os.Stderr, `
Mrglent filters financial transactions in Ledger entry format from multiple journals into a general journal.

Mrglent reads Ledger journals from standard input.
It extracts dated journal entries.
If an entry's date cannot be parsed according to the layout, 
mrglent writes a message to standard error and exits with a non-zero status.
Dated entries marked as mirrors (between "# mirror entry" and "# end mirror entry" comment lines) are discarded.
See this module's program mcsv2lent for more on marked entries.
All other journal content is also discarded including automatic transactions, global comments and command directives.

Mrglent orders the entries by date ascending and writes them to standard output.

Usage:

	mrglent [flags]

The flags are:

`)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}
