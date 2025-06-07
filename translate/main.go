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
Translate [filters] transaction records from a
[comma-separated values (CSV)] account statement to a standard format.

# Examples

Assuming translate is installed in a Unix-like environment and
being run from the translate source code directory,
here are some examples of its use.

## Run with default input and output formats

	echo "1982-10-08,12-3456-7890123-45,01-2345-6789012-34,Daily allowance,-30,ALD" | \
		translate

By default, transactions are translated
from this module's CSV records to [Ledger] journal entries.
The fields in the record are date,
this account (the one this transaction belongs to), other account,
memo, amount and currency.

## Configure input format

	translate -f national_bank.xml <national_bank.csv

The format of a National Bank CSV statement is configured by the [XML] file (-f flag).
Translate warns about the statement's header line, which cannot be translated.
If a transaction's other account field is empty string,
translate resets it to "Imbalance".

## Set this account and output format

	translate -f local_CU.xml -t Assets:Saving -o csv <local_CU.csv

In contrast to the last example, a Local Credit Union statement:

  - does not contain its own account number in a this account field
  - contains debit and credit fields instead of an amount field
  - orders transactions by date descending instead of ascending

This account is set to its Ledger name (-t flag).
The output format is set to this module's CSV records (-o flag).
Translate outputs transactions ordered by date ascending.

## Merge CSV statements to Ledger journal

	translate -f national_bank.xml -t Assets:Current -o csv <national_bank.csv >all.csv
	translate -f local_CU.xml -t Assets:Saving -o csv <local_CU.csv >>all.csv
	sed -E -f adjust.sed all.csv | sort -t , -k 1 | translate

Each CSV statement is translated to this module's CSV records.
The stream editor replaces account numbers with names and removes mirrored transactions.
The CSV records are sorted by date ascending then translated to Ledger journal entries.

[comma-separated values (CSV)]: https://en.wikipedia.org/wiki/Comma-separated_values
[filters]: https://en.wikipedia.org/wiki/Filter_(software)
[Ledger]: https://en.wikipedia.org/wiki/Ledger_(software)
[XML]: https://en.wikipedia.org/wiki/XML
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

// The configuration returned by parseFlags.
type config struct {
	formatFileName string
	help           bool
	outFormatName  string
	thisAccount    string
}

var errThisAccount = errors.New("this account and its index cannot be empty string and zero")

func main() {
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)

	cfg := parseFlags()
	if cfg.help {
		usage()
		os.Exit(0)
	}

	var err error

	inFormat := aft.GetModuleCSVFormat()
	if cfg.formatFileName != "" {
		inFormat, err = aft.ReadCSVFormat(cfg.formatFileName)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = inFormat.Validate()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.thisAccount == "" && inFormat.ThisAccountI == 0 {
		log.Fatal(errThisAccount)
	}

	r := csv.NewReader(os.Stdin)
	/*
		The number of fields in a record is checked by function aft.ParseCSV,
		so disable the reader's check.
	*/
	r.FieldsPerRecord = -1

	ts, err := parseTransactions(r, cfg.thisAccount, inFormat)
	if err != nil {
		log.Fatal(err)
	}

	stringTransactions(ts, os.Stdout, cfg.outFormatName)
}

// ParseFlags returns the values of flags from the command that ran translate.
func parseFlags() config {
	var cfg config

	flag.StringVar(&cfg.formatFileName, "f", "",
		"input CSV record format file name (default input format is this module's CSV record)")
	flag.BoolVar(&cfg.help, "h", false, "write this help text then exit")
	flag.StringVar(&cfg.outFormatName, "o", aft.Ledger,
		fmt.Sprintf("output format name: %q or %q for this module's CSV record", aft.Ledger, aft.CSV))
	flag.StringVar(&cfg.thisAccount, "t", "",
		"name of this account: the one that this statement and its transactions belong to")

	flag.Usage = usage
	flag.Parse()

	if cfg.outFormatName == "" {
		cfg.outFormatName = aft.Ledger
	}

	return cfg
}

/*
ParseTransactions reads a CSV account statement,
parses a transaction's fields from those in the CSV record on each line
then returns the transactions.
If it fails to read the statement, parseTransactions returns an error.
*/
func parseTransactions(r *csv.Reader, thisAccount string, crf aft.CSVRecordFormat) ([]aft.Transaction, error) {
	var ts []aft.Transaction

	for {
		fs, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return ts, err
		}

		var t aft.Transaction

		t.ThisAccount = thisAccount

		err = t.ParseCSV(fs, crf)
		if err != nil {
			i, _ := r.FieldPos(0)
			log.Printf("%v on line %v", err, i)

			continue
		}

		ts = append(ts, t)
	}

	return ts, nil
}

// StringTransactions writes the transactions in the named format ordered by date ascending.
func stringTransactions(ts []aft.Transaction, w *os.File, name string) {
	n := len(ts)

	tSeq := slices.All(ts)
	if 2 <= n && ts[0].Date > ts[n-1].Date {
		tSeq = slices.Backward(ts)
	}

	for _, t := range tSeq {
		fmt.Fprint(w, t.StringFormat(name))
	}
}

// Usage prints the help text for translate.
func usage() {
	fmt.Fprint(os.Stderr, `
Translate filters financial transaction records
from a comma-separated values (CSV) account statement to a standard format.

Usage:

	translate [flags]

The flags are:

`)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}
