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
Mergetrn [filters] multiple [Ledger] journals, containing financial transactions:
  - discarding mirrored transactions that have been marked with the code "(MT)"
  - ordering the remaining transactions by date ascending

It assumes the entries in the Ledger journals are valid, and their date layout is "YYYY-MM-DD".

Assuming multiple accounts each with its own Ledger journal,
transfers between those accounts will lead to mirrored transactions.
A mirrored transaction is a debit in one journal mirrored by a credit in another.
When those journals are merged,
one side of each mirrored transaction must be discarded
to avoid making the transfer twice.

# Example

The following example shows how to use mergetrn.
It assumes mergetrn is installed in a Unix-like environment
and is being run from its source directory.

	cat LCU.journal NB_current.journal NB_emergency.journal | mergetrn

Credit mirrored transactions have been marked with code "(MT)"
in the National Bank emergency and Local Credit Union journals.
Mergetrn discards those transactions, and orders the remainder by date ascending.

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
	"sort"
)

const (
	mirrorCode = "MT"
	sp         = ' '
)

func main() {
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)

	parseFlags()

	var date2lines = make(map[string]string) // lines of transactions made on a date

	date := "0000-00-00" // date of current transaction
	date2lines[date] = ""

	var discard bool // whether to discard current transaction

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		ln := s.Text()

		if len(ln) == 0 {
			// blank line
			discard = false
			date2lines[date] += "\n"

			continue
		}

		if ln[0] == sp {
			// indented line that continues current transaction
			if !discard {
				date2lines[date] += ln + "\n"
			}

			continue
		}

		// global line, which is not indented
		discard = false

		var trn aft.Transaction

		err := trn.ParseLedger([]string{ln})
		if err != nil {
			date2lines[date] += ln + "\n"

			continue
		}

		date = trn.Date

		// first line of the next transaction
		if trn.Code == mirrorCode {
			// this transaction is mirrored
			discard = true

			continue
		}

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
Merge filters financial transactions from multiple Ledger journals:
 - discarding mirrored transactions that have been marked with the code "(MT)"
 - ordering the remaining transactions by date ascending

It assumes the transaction's date layout is "2006-01-02" also known as time.DateOnly.
The only flag is:

`)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}
