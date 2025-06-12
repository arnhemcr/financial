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
CSV2trn [filters] financial transactions from [comma-separated values (CSV)] records,
in an account statement, to a standard format.

# Examples

The following examples explain how to use csv2trn.
They assume csv2trn is installed in a Unix-like environment
and is being run from its source directory.

## Default input and output format

Running:

	echo "1982-10-08,Assets:Saving,Assets:Current,DB,Daily allowance,-30,ALD" | csv2trn

should produce the [Ledger] journal entry:

	1982-10-08 (DB) Daily allowance
	 Assets:Saving  -30 ALD
	 Assets:Current

In csv2trn, a transaction has the following fields:

  - Date
  - This account, which the transaction belongs to e.g. "Assets:Saving"
  - That account e.g. "Assets:Current"
  - Code, also known as the type of transaction
  - Memo, also known as the description
  - Amount
  - Currency

This example shows the default input and output formats.
CSV2trn reads fields from the CSV record, in this module's format, into a transaction.
It then writes the transaction as a Ledger journal entry.

## Custom input format

	cat NB_current.csv | ./csv2trn -f NB.xml
	cat NB_emergency.csv | ./csv2trn -f NB.xml

This example translates two CSV account statements
from National Bank to Ledger journal entries.
The bank has its own CSV statement and record format.
An [XML] file (-f flag) configures csv2trn for the input format with
the number and position of fields in the record and the date layout.

The example highlights a couple of points.
First, if a line in a statement cannot be interpreted as a transaction,
for example a header line, csv2trn writes a warning and continues.
Second, the transfer from current to emergency accounts appears as two mirrored entries:
a debit from the current account and a credit to the emergency account.

## CSV records without this account

	cat LCU.csv | csv2trn -f LCU.xml -t Assets:Saving -o modcsv

This example translates a CSV account statement from Local Credit Union to
this module's CSV records (-o flag).
The credit union's CSV records do not contain this account,
so its value is set from the command line (-t flag).
Instead of an amount field, the records contain credit and debit fields
from which amount is calculated.

Records in Local Credit Union account statements are ordered by date descending.
CSV2trn always writes transactions ordered by date ascending.

## Ledger journal from CSV account statements

	cat NB_current.csv | ./csv2trn -f NB.xml -t Assets:Current -c GBP -o modcsv >all.csv
	cat NB_emergency.csv | ./csv2trn -f NB.xml -t Assets:Emergency -c GBP -o modcsv >>all.csv
	cat LCU.csv | ./csv2trn -f LCU.xml -t Assets:Saving -c GBP -o modcsv >>all.csv

	sed -E -f adjust.sed all.csv | sort -t , -k 1 | ./csv2trn >all.journal
	ledger -f all.journal balance

All the statements are translated into this module's CSV format.
None of them contain currency, so it is set from the command line (-c flag).
For National Bank, the statement this account numbers
are overridden with names from the command line.

The stream editor (sed) replaces other account numbers with names
and removes mirrored records.
The records are sorted to date ascending then translated to Ledger journal entries.
Finally, ledger reports balances for the journal.

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
	currency       string
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

	switch cfg.outFormatName {
	case aft.Ledger, aft.ModuleCSV:
		// This output format name is valid.
	default:
		log.Fatalf("%q: %v", cfg.outFormatName, errOutFormatName)
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
		Set the reader to reuse its record, instead of reallocating,
		to improve performance.
	*/
	r.FieldsPerRecord = -1
	r.ReuseRecord = true

	ts, err := parseTransactions(r, cfg.thisAccount, cfg.currency, inFormat)
	if err != nil {
		log.Fatal(err)
	}

	stringTransactions(ts, os.Stdout, cfg.outFormatName)
}

// ParseFlags returns the values of flags from the command that ran csv2trn.
func parseFlags() config {
	var cfg config

	flag.StringVar(&cfg.currency, "c", "", fmt.Sprintf(
		"currency for transaction amounts: %q or a currency code e.g. %q", "$", "GBP"))
	flag.StringVar(&cfg.formatFileName, "f", "", "file name containing input CSV record format")
	flag.BoolVar(&cfg.help, "h", false, "help text for this program then exit")
	flag.StringVar(&cfg.outFormatName, "o", aft.Ledger,
		fmt.Sprintf("output format name: %q or %q", aft.Ledger, aft.ModuleCSV))
	flag.StringVar(&cfg.thisAccount, "t", "", fmt.Sprintf(
		"this account name: account the transactions belong to e.g. %q", "Assets:Current"))

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
func parseTransactions(r *csv.Reader, thisAccount, currency string,
	crf aft.CSVRecordFormat) ([]aft.Transaction, error) {
	var ts []aft.Transaction

	for {
		fs, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return ts, fmt.Errorf("r.Read: %w", err)
		}

		var t aft.Transaction

		t.Currency = currency
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

// Usage prints the help text for csv2trn.
func usage() {
	fmt.Fprint(os.Stderr, `
CSV2trn filters financial transactions from comma-separated values (CSV) records,
in an account statement, to a standard format.

Usage:

	csv2trn [flags]

The flags are:

`)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}
