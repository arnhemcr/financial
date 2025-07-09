/*
Copyright (C) 2025 Andrew Flint.

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
Mrglent [filters] multiple [Ledger] financial journals containing entries,
also known as transactions.
It:
  - removes mirror entries that have been marked with the transaction code "MT"
  - sorts the remaining entries by date ascending

Assuming multiple accounts each with its own Ledger journal,
transfers between those accounts will lead to mirror entries.
A mirror entry is a debit in one journal mirrored by a credit in another.
When those journals are merged,
one side of each mirror entry must be removed to avoid making the transfer twice.

Mrglent only merges entries: transactions which have dates.
Automatic transactions, comments and ommand directives in the input journals
are not currently copied to the output journal.
It uses the same "YYYY-MM-DD" date layout, for input and output,
as used by other arnhemcr/financial programs for output.

# Example

The following example shows how to use mrglent.
It assumes mrglent is installed in a Unix-like environment
and is being run from its source directory.

Start by concatenating the three journals into one:

	cat *.journal >g01
	ledger -f g01 register emergency

The output has six lines, the last four are a pair of mirror entries,
and the second to last line is out of date order.

Now use mrglent:

	cat *.journal | mrglent >g02
	ledger -f g02 register emergency

Mrglent removes both credit mirror entries, which have been marked with code "MT",
then sorts the remaining four entries by date ascending.

The format of a Ledger journal entry is described in the
"Transactions and Comments" section of the [Ledger 3 manual].

[filters]: https://en.wikipedia.org/wiki/Filter_(software)
[Ledger]: https://en.wikipedia.org/wiki/Ledger_(software)
[Ledger 3 manual]: https://ledger-cli.org/doc/ledger3.html
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
	"unicode"
)

func main() {
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)

	parseFlags()

	s := bufio.NewScanner(os.Stdin)

	var j ledgerJournal

	err := j.parse(s)
	if err != nil {
		log.Fatal(err)
	}

	j.demirror()
	j.sort()

	for _, e := range j.Entries {
		fmt.Fprint(os.Stdout, e.string())
	}
}

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

func usage() {
	fmt.Fprint(os.Stderr, `
Mrglent filters multiple Ledger financial journals containing entries,
also known as transactions. It:
  - removes mirror entries that have been marked with the transaction code "MT"
  - sorts the remaining entries by date ascending

The only flag is:

`)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}

/*
A ledgerEntry represents an entry, also known as a transaction, in a Ledger journal.
See also "Transactions and Comments" in the [Ledger 3 manual].

[Ledger 3 manual]: https://ledger-cli.org/doc/ledger3.html
*/
type ledgerEntry struct {
	Lines []string // this entry's lines
	/*
		Trn contains copies of this entry's date, code and memo fields.
		See also aft.ParseLedger.
	*/
	Trn aft.Transaction
}

/*
A ledgerJournal represents the entries in a Ledger journal.

In future, fields may be added to represent other Ledger journal elements
including automatic transactions, comments and command directives.
*/
type ledgerJournal struct {
	Entries []ledgerEntry // this journal's entries
}

/*
Parse parses this Ledger entry from its lines in a journal.
If it fails to parse the entry, parse returns an error.
*/
func (e *ledgerEntry) parse(lines []string) error {
	err := e.Trn.ParseLedger(lines)
	if err != nil {
		return fmt.Errorf("e.parse: %w", err)
	}

	e.Lines = lines

	return nil
}

// String returns this Ledger entry as a string.
func (e *ledgerEntry) string() string {
	var s string

	for _, ln := range e.Lines {
		s += ln
	}

	return s
}

// Demirror removes any entries which are marked as mirrored from this Ledger journal.
func (j *ledgerJournal) demirror() {
	es := make([]ledgerEntry, len(j.Entries))

	for _, e := range j.Entries {
		if e.Trn.Code == "MT" {
			continue
		}

		es = append(es, e)
	}

	j.Entries = es
}

/*
Parse parses the entries in this Ledger journal from a stream of lines.
If it fails to parse the journal, parse returns the first error.
*/
func (j *ledgerJournal) parse(s *bufio.Scanner) error {
	var (
		e   ledgerEntry
		lns []string // lines of the current entry
		err error
	)

	for s.Scan() {
		ln := s.Text() + "\n"

		r0 := rune(ln[0])

		if 0 < len(lns) {
			if r0 == ' ' || r0 == '\t' {
				// This line continues the current entry.
				lns = append(lns, ln)

				continue
			}

			// This line is the first one after the current entry.
			err = e.parse(lns)
			if err != nil {
				return err
			}

			lns = nil

			j.Entries = append(j.Entries, e)
		}

		if unicode.IsDigit(r0) {
			// This line starts an entry.
			lns = append(lns, ln)

			continue
		}
	}

	err = s.Err()
	if err != nil {
		return fmt.Errorf("j.parse: %w", err)
	}

	if 0 < len(lns) {
		err = e.parse(lns)
		if err != nil {
			return err
		}

		j.Entries = append(j.Entries, e)
	}

	return nil
}

// Sort sorts the entries in this Ledger journal by date ascending.
func (j *ledgerJournal) sort() {
	d2es := make(map[string][]ledgerEntry)

	var ds []string

	for _, e := range j.Entries {
		d := e.Trn.Date

		_, found := d2es[d]
		if !found {
			d2es[d] = []ledgerEntry{}

			ds = append(ds, d)
		}

		d2es[d] = append(d2es[d], e)
	}

	slices.Sort(ds)

	es := make([]ledgerEntry, len(j.Entries))

	for _, d := range ds {
		es = append(es, d2es[d]...)
	}

	j.Entries = es
}
