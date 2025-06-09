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
here are some examples to try.

## Run with default input and output formats

	echo "1982-10-08,Assets:Saving,Assets:Current,DB,Daily allowance,-30,ALD" | \
		./translate

By default, transactions are translated
from this module's CSV records to [Ledger] journal entries.
The fields in the record are date,
this account (the one this transaction belongs to), other account,
code (or transaction type), memo (or description), amount and currency.

## Configure input format

	translate -f national_bank.xml <national_bank.csv

A CSVRecordFormat, read from the [XML] file (-f flag),
maps the fields in an input CSV record to the fields in a transaction within translate.
If the other account field in the record is an empty string,
translate sets that field in the transaction to "Imbalance".

## Set this account and output format

	translate -f local_CU.xml -t Assets:Saving -o modcsv <local_CU.csv

Another financial institution with a different CSV record format and its own XML file.
This input record does not contain the this account field,
so it has to be set (-t flag).
This statement orders transactions by date descending.

Translate outputs the transactions as this module's CSV records (-o flag),
and reverses their order to date ascending.

## Merge CSV statements to Ledger journal

	translate -f national_bank.xml -t Assets:Current -o modcsv <national_bank.csv >all.csv
	translate -f local_CU.xml -t Assets:Saving -o modcsv <local_CU.csv >>all.csv
	sed -E -f adjust.sed all.csv | sort -t , -k 1 | translate >all.journal
	ledger -f all.journal balance

Both statements are translated to this module's CSV records.
The stream editor replaces account numbers with names and removes mirrored transactions.
The records are sorted to date ascending then translated to entries in a Ledger journal.
Finally, Ledger reports balances for the journal.

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

var (
	errOutFormatName = errors.New("no such output format name")
	errThisAccount   = errors.New("cannot get this account: " +
		"CSV records do not contain that field and its flag is empty string")
)

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

	switch cfg.outFormatName {
	case aft.Ledger:
	case aft.ModuleCSV:
	default:
		log.Fatalf("%q: %v", cfg.outFormatName, errOutFormatName)
	}

	if cfg.thisAccount == "" && inFormat.ThisAccountI == 0 {
		log.Fatal(errThisAccount)
	}

	r := csv.NewReader(os.Stdin)
	/*
		The number of fields in a record is checked by function aft.ParseCSV,
		so disable the reader's check.
		Set the reader to reuse its record, instead of reallocating,
		to improve performance.
	*/
	r.FieldsPerRecord = -1
	r.ReuseRecord = true

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
		fmt.Sprintf("output format name: %q or %q for this module's CSV record", aft.Ledger, aft.ModuleCSV))
	flag.StringVar(&cfg.thisAccount, "t", "",
		"name of this account: the one that this statement and its transactions belong to")

	flag.Usage = usage
	flag.Parse()

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
			return ts, fmt.Errorf("r.Read: %w", err)
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

/*
StringTransactions writes the transactions in the named format.
It assumes the transactions are in date order either ascending or descending.
If the first transaction is later than the last one,
stringTransactions reverses the order.
*/
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
