/*
Copyright (C) 2025-2026 Andrew Flint.

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
Mrglent [filters] financial transaction entries from multiple [Ledger] journals into one general journal.

For further information, see [this package's README].

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

	parseFlags()

	s := bufio.NewScanner(os.Stdin)

	es, err := parseEntries(s)
	if err != nil {
		log.Fatal(err)
	}

	oes := sortEntries(es)
	for _, oe := range oes {
		fmt.Fprint(os.Stdout, oe)
	}
}

/*
IsIndent reports whether the rune is white space used to indent a Ledger entry's postings.
See "Transactions and Comments" in the [Ledger 3 manual].
*/
func isIndent(r rune) bool {
	switch r {
	case ' ', '\t':
		return true
	default:
		return false
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
func parseEntries(s *bufio.Scanner) ([]entry, error) {
	var (
		es                            []entry
		e                             entry
		inBlockComment, inMirrorEntry bool
		lnN                           int
	)

	for s.Scan() {
		ln := s.Text() + "\n"
		lnN++

		switch {
		case ln == aft.StartBlockComment:
			inBlockComment = true

			continue
		case ln == aft.EndBlockComment:
			inBlockComment = false

			continue
		case inBlockComment:
			continue
		}

		switch {
		case ln == aft.StartMirrorEntry:
			inMirrorEntry = true

			continue
		case ln == aft.EndMirrorEntry:
			inMirrorEntry = false

			continue
		case inMirrorEntry:
			continue
		}

		r0 := rune(ln[0])

		switch {
		case unicode.IsDigit(r0):
			if e.Date != "" {
				es = append(es, e)
			}

			d, err := parseDate(ln)
			if err != nil {
				return es, fmt.Errorf("input line %v: %w", lnN, err)
			}

			e.Date, e.Text = d, ln
		case isIndent(r0):
			e.Text += ln
		}
	}

	if e.Date != "" {
		es = append(es, e)
	}

	return es, nil
}

/*
ParseDate returns the date in "2006-01-02" layout from the start of a Ledger entry's text.
If parseDate fails to parse a date, it returns the error.
*/
func parseDate(text string) (string, error) {
	const dLen = len(time.DateOnly)

	var (
		tLen = len(text)
		d    string
	)

	if dLen <= tLen {
		d = text[0:dLen]
	} else {
		d = text[0:tLen]
	}

	date, err := time.Parse(time.DateOnly, d)
	if err != nil {
		return "", fmt.Errorf("parseDate: %w", err)
	}

	return date.Format(time.DateOnly), nil
}

/*
ParseFlags parses this program's command line flags.
If help was requested, parseFlags writes this program's help text then exits.
If the flags are invalid, this program exits with a non-zero status.
*/
func parseFlags() {
	var help bool

	flag.BoolVar(&help, "h", false, "write this help text then exit")

	flag.Usage = usage
	flag.Parse()

	if help {
		usage()
		os.Exit(0)
	}
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
Mrglent filters financial transaction entries
from multiple Ledger journals into one general journal.
Entries with dates are copied from input to output.
But dated entries marked as mirrors (between "mirror entry" comments),
automatic transactions, comments and command directives are discarded.
Mrglent sorts the general ledger entries by date ascending.
See also program mcsv2lent.

The only flag is:

`)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}
