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
MCSV2lent:
  - [filters] financial transactions from [comma-separated values (CSV)] records
    in this module's CSV format to [Ledger] journal entries
  - optionally, marks credits between Ledger accounts with journals as mirrored

[comma-separated values (CSV)]: https://en.wikipedia.org/wiki/Comma-separated_values
[filters]: https://en.wikipedia.org/wiki/Filter_(software)
[Ledger]: https://en.wikipedia.org/wiki/Ledger_(software)
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
		jas, err = aft.LedgerAccounts(jafn)
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
			log.Fatal(err)
		}

		ent := t.StringLedger()

		if 0 < t.Amount && slices.Contains(jas, t.ThisAccount) &&
			slices.Contains(jas, t.OtherAccount) {
			ent = aft.StartMirror + ent + aft.EndMirror
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
		"name of file listing Ledger accounts with journals")

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
MCSV2lent

Usage:

	mcsv2lent [flags]

The flags are:

`)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}
