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
	  	Go-style date layout of Ledger journal entries (default "2006-01-02")
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
		log.Fatalf("%v not %q", ErrDateLayout, dateLayout)
	}

	s := bufio.NewScanner(os.Stdin)

	es, err := parseEntries(s, dateLayout)
	if err != nil {
		log.Fatal(err)
	}

	etexts := sortEntries(es)
	for _, etext := range etexts {
		fmt.Fprint(os.Stdout, etext)
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
func parseEntries(s *bufio.Scanner, dateLayout string) (es []entry, err error) {
	var (
		inBlockComment, inMirrorEntry bool
		e                             entry
		n                             int // The number of the line being parsed.
	)

	for s.Scan() {
		line := s.Text() + "\n"
		n++

		if inBlock(line, aft.StartBlockComment, aft.EndBlockComment, &inBlockComment) ||
			inBlock(line, aft.StartMirrorEntry, aft.EndMirrorEntry, &inMirrorEntry) {
			continue
		}

		switch {
		case unicode.IsDigit(rune(line[0])):
			if e.Date != "" {
				es = append(es, e)
			}

			d, err := aft.ParseDate(line, dateLayout)
			if err != nil {
				return es, fmt.Errorf("line %v: %w", n, err)
			}

			// This line starts with a date and is the first line in the next entry.
			e.Date, e.Text = d, line
		case aft.IsLedgerIndented(line):
			// This indented line is a continuation of the current entry.
			e.Text += line
		}
	}

	if e.Date != "" {
		es = append(es, e)
	}

	return es, nil
}

/*
InBlock reports whether the line from a Ledger journal is in a block delimited by start and end lines.
It may also update the in block state.
*/
func inBlock(line, start, end string, state *bool) bool {
	switch {
	case line == start:
		*state = true
	case line == end:
		*state = false
	case *state:
		// This line is in a block between start and end lines.
	default:
		return false
	}

	return true
}

// Sort orders the texts of the list of Ledger journal entries by date ascending.
func sortEntries(es []entry) (etexts []string) {
	d2ets := make(map[string][]string) // The map of entry dates to entry texts.
	ds := []string{}                   // The list of entry dates.

	for _, e := range es {
		d := e.Date

		_, found := d2ets[d]
		if !found {
			d2ets[d] = []string{}

			ds = append(ds, d)
		}

		d2ets[d] = append(d2ets[d], e.Text)
	}

	slices.Sort(ds)

	for _, d := range ds {
		etexts = append(etexts, d2ets[d]...)
	}

	return etexts
}

/*
ParseFlags returns the date layout of Ledger journal entries parsed from command line flags.
If help was requested, parseFlags writes this program's help text then exits.
If the flags are invalid, this program exits with a non-zero status.
*/
func parseFlags() string {
	var dateLayout string

	flag.StringVar(&dateLayout, "d", time.DateOnly, "Go-style date layout of Ledger journal entries")

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
