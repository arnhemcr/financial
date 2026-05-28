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
MCSV2lent [filters] financial transactions
from this module's [comma-separated values (CSV)] records to [Ledger] journal entries.
If given a list of Ledger account names with journals,
it also marks the credit entry of transfers between those accounts.
Marked entries are discarded when those journals are merged by this module's program mrglent.

MCSV2lent reads lines from standard input.
It parses each line as a transaction CSV record in this module's format (mcsv).
If a line cannot be parsed, mcsv2lent writes a message to standard error and exits with a non-zero status.

The list of Ledger account names with journals is empty by default, but it can be loaded from an XML file.
For example:

	<LedgerAccountsWithJournals>
	    <Account>Assets:Current</Account>    <!-- NB.journal -->
	    <Account>Assets:Emergency</Account>  <!-- LCU.journal -->
	</LedgerAccountsWithJournals>

A transaction whose amount is positive and whose this and other accounts are both on the list is marked.

MCSV2lent writes transactions to standard output in Ledger journal entry format.
The entry for a marked transaction
is preceded by Ledger global comment line "# mirror entry" and followed by "# end mirror entry".

Usage:

	mcsv2lent [flags]

The flags are:

	-f string
	      name of file containing list of Ledger accounts with journals in XML
	-h    write this help text then exit

See also [this package's README].

[comma-separated values (CSV)]: https://en.wikipedia.org/wiki/Comma-separated_values
[filters]: https://en.wikipedia.org/wiki/Filter_(software)
[Ledger]: https://en.wikipedia.org/wiki/Ledger_(software)
[this package's README]: https://github.com/arnhemcr/financial/tree/main
*/
package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	aft "github.com/arnhemcr/financial/transaction"
	"io"
	"log"
	"os"
	"slices"
)

func main() {
	log.SetPrefix("mcsv2lent: ")
	log.SetFlags(0)

	jafn := parseFlags() // The name of the file listing Ledger accounts with journals.

	var jas []string // The list of Ledger accounts with journals.

	var err error

	if jafn != "" {
		jas, err = aft.LoadLedgerAccountNames(jafn)
		if err != nil {
			log.Fatal(err)
		}
	}

	r := csv.NewReader(os.Stdin)
	r.FieldsPerRecord, r.ReuseRecord = -1, true

	mcsv := aft.NewModuleCSVRecordFormat()

	for {
		fs, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		var t aft.Transaction

		err = t.ParseCSV(fs, mcsv)
		if err != nil {
			n, _ := r.FieldPos(0)
			log.Fatalf("%v on line %v", err, n)
		}

		ent := t.StringLedger()

		if 0 < t.Amount && slices.Contains(jas, t.ThisAccount) &&
			slices.Contains(jas, t.OtherAccount) {
			ent = aft.StartMirrorEntry + ent + aft.EndMirrorEntry
		}

		fmt.Fprint(os.Stdout, ent)
	}
}

/*
ParseFlags returns this program's configuration parsed from command line flags.
If help was requested, parseFlags writes help text then exits.
If the flags are invalid, this program exits with a non-zero status.
*/
func parseFlags() string {
	var fileName string

	flag.StringVar(&fileName, "f", "",
		"name of file containing list of Ledger accounts with journals in XML")

	var help bool

	flag.BoolVar(&help, "h", false, "write this help text then exit")

	flag.Usage = usage
	flag.Parse()

	if help {
		usage()
		os.Exit(0)
	}

	return fileName
}

// Usage writes the help text for this program.
func usage() {
	fmt.Fprint(os.Stderr, `
MCSV2lent filters financial transactions
from this module's comma-separated values (CSV) records to Ledger journal entries.
If given a list of Ledger account names with journals,
it also marks the credit entry of transfers between those accounts.
Marked entries are discarded when those journals are merged by this module's program mrglent.

MCSV2lent reads lines from standard input.
It parses each line as a transaction CSV record in this module's format (mcsv).
If a line cannot be parsed, mcsv2lent writes a message to standard error and exits with a non-zero status.

The list of Ledger account names with journals is empty by default, but it can be loaded from an XML file.
For example:

    <LedgerAccountsWithJournals>
        <Account>Assets:Current</Account>    <!-- NB.journal -->
        <Account>Assets:Emergency</Account>  <!-- LCU.journal -->
    </LedgerAccountsWithJournals>

A transaction whose amount is positive and whose this and other accounts are both on the list is marked.

MCSV2lent writes transactions to standard output in Ledger journal entry format.
The entry for a marked transaction
is preceded by Ledger global comment line "# mirror entry" and followed by "# end mirror entry".

Usage:

	mcsv2lent [flags]

The flags are:

`)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}
