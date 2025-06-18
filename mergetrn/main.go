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
Mergetrn [filters] financial transactions in the [Ledger] journal entry format:
  - discarding mirrored transactions that have the code "(MT)"
  - ordering the remaining transactions by date ascending

It assumes the date layout of journal entries is "YYYY-MM-DD".

Mirrored transactions are an issue when merging journals of accounts
that have transfers between those accounts.
Each transfer has a debit transaction in one journal
mirrored by a credit transaction in another.
If one of the mirrored transactions is not discarded, the result is a double transfer.

# Example

The following example shows how to use mergetrn.
It assumes mergetrn is installed in a Unix-like environment
and is being run from its source directory.

	cat LCU.journal NB_current.journal NB_emergency.journal | mergetrn

Credit mirrored transactions have been marked with code "(MT)"
in the National Bank emergency and Local Credit Union journals.
Mergetrn discards those transactions, and orders the remainder by date ascending.

[filters]: https://en.wikipedia.org/wiki/Filter_(software)
[Ledger]: https://en.wikipedia.org/wiki/Ledger_(software)

[Ledger 3 Manual]: https://ledger-cli.org/doc/ledger3.html
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	aft "github.com/arnhemcr/financial/transaction"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

const mirrorCode = "(MT)"

func main() {
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)

	parseFlags()

	var date string // The date of the current transaction.

	var discard bool // Whether to discard the current transaction.

	// Map each date to the lines of transactions made on that date.
	date2lines := make(map[string]string)

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		ln := s.Text()

		fs := strings.Fields(ln)

		n := len(fs)
		if n == 0 {
			continue // Discard this blank line.
		}

		const sp = ' '

		if ln[0] == sp && !discard {
			/*
				This indented line belongs to the current transaction.
				It could be an account line or a comment
				(see "The Most Basic Entry" in [Ledger 3 Manual].
			*/
			date2lines[date] += ln + "\n"

			continue
		}

		/*
			This line is not indented,
			so it does not belong to any previous transaction.
		*/
		discard = false

		d, _ := aft.ParseDate(fs[0], time.DateOnly)
		if d == "" {
			/*
				Discard this line, which is probably a global comment
				(see "Commenting on your Journal" in [Ledger 3 Manual]).
			*/
			continue
		}

		// This is the first line in the next transaction.
		date = d

		if 2 <= n && fs[1] == mirrorCode {
			// This mirrored transaction will be discarded.
			discard = true

			continue
		}

		// Ensure there is a lines string for this date.
		_, found := date2lines[date]
		if !found {
			date2lines[date] = ""
		}

		date2lines[date] += ln + "\n"
	}

	err := s.Err()
	if err != nil {
		log.Fatalf("s.err: %v", err)
	}

	outputInOrder(date2lines)
}

// OutputInOrder prints the lines of the remaining transactions ordered by date ascending.
func outputInOrder(date2lines map[string]string) {
	var ds []string

	for d := range date2lines {
		ds = append(ds, d)
	}

	sort.Strings(ds)

	for _, d := range ds {
		fmt.Fprint(os.Stdout, date2lines[d])
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
Mergetrn filters financial transactions in the Ledger journal format:
 - discarding mirrored transactions that have the code "(MT)"
 - ordering the remaining transactions by date ascending

It assumes the transaction's date layout is "2006-01-02" also known as time.DateOnly.
The only flag is:

`)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}
