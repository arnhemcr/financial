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
MCSV2lent translates financial transactions from this module's [CSV] records (mcsv) to [Ledger] journal entries (lent).

It:
  - read CSV records from standard input
  - parses a transaction from each record following this module's format
  - writes transactions to standard output as Ledger journal entries
  - encloses the credit entry for a transaction between journalled accounts with comments

Usage:

	mcsv2lent [flag]

The flag is:

	-f string
	      name of XML file listing Ledger account names with journals

See also [this module's README].

[CSV]: https://en.wikipedia.org/wiki/Comma-separated_values
[Ledger]: https://en.wikipedia.org/wiki/Ledger_(software)
[this module's README]: https://github.com/arnhemcr/financial/tree/main
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

	var (
		jas []string // The list of Ledger journalled account names.
		err error
	)

	fileName := parseFlags()
	if fileName != "" {
		jas, err = aft.LoadLedgerAccountNames(fileName)
		if err != nil {
			log.Fatal(err)
		}
	}

	mcsv := aft.NewModuleCSVRecordFormat()
	r := csv.NewReader(os.Stdin)
	r.FieldsPerRecord, r.ReuseRecord = -1, true

	for {
		fs, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			log.Fatal(err)
		}

		var t aft.Transaction

		err = t.ParseCSV(fs, mcsv)
		if err != nil {
			n, _ := r.FieldPos(0)
			log.Fatalf("line %v: %v", n, err)
		}

		ent := t.StringLedger()

		if 0 < t.Amount &&
			slices.Contains(jas, t.ThisAccount) &&
			slices.Contains(jas, t.OtherAccount) {
			ent = aft.StartMirrorEntry + ent + aft.EndMirrorEntry
		}

		fmt.Fprint(os.Stdout, ent)
	}
}

/*
ParseFlags parses this program's configuration parsed from command line flags.
If the flags are invalid, this program exits with a non-zero status.
ParseFlags returns the name of a file or empty string if that flag was not set.
*/
func parseFlags() (fileName string) {
	flag.StringVar(&fileName, "f", "", "name of XML file listing Ledger account names with journals")

	flag.Usage = usage
	flag.Parse()

	return fileName
}

// Usage writes the help text for this program.
func usage() {
	fmt.Fprint(os.Stderr, `
MCSV2lent translates financial transactions from this module's CSV records (mcsv) to Ledger journal entries (lent).

It:
  - read CSV records from standard input
  - parses a transaction from each record following this module's format
  - writes transactions to standard output as Ledger journal entries
  - encloses the credit entry for a transaction between journalled accounts with comments

Usage:

	mcsv2lent [flag]

The flag is:

`)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}
