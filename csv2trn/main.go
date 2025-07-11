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

## National Bank

Try:

	cat NB_01_emergency.csv | csv2trn -f NB.xml 

This module's transactions have the following fields:
  - Date
  - This account, which is the account this transaction belongs to
  - Other account, which is optional
  - Code, which is also known as the type and is optional
  - Memo, which is also known as the description
  - Amount
  - Currency, which is optional

National Bank CSV account statements contain all these fields except currency.
In this example, csv2trn reads CSV records from a National Bank account statement.
A transaction is parsed from the record on each line according to the input format in [XML].
Each transaction is written in the default output format: [Ledger] journal entry.
A warning is written about the unparseable header line, but csv2trn continues.

Try:

	cat NB_02_current.csv | csv2trn -f NB.xml -c GBP -t Assets:Current

In this example, command line flags set the currency 
and override this account number, from the CSV records, with a Ledger hierachical name.

## Local Credit Union

Try:

	cat LCU.csv | csv2trn -f LCU.xml -t Assets:Saving

Local Credit Union statements contain debit and credit instead of amount,
and they do not contain this or other account.
In this example, the XML input format configures csv2trn to get the amount from credit and debit. 
And this account must be set with a flag.

## Towards Ledger

Try:

	cat NB_02_current.csv | csv2trn -f NB.xml -o modcsv | sed -f adjust.sed | csv2trn

Ledger expects both this and other accounts to be hierarchical names.
If a CSV record contains an other account number, it can be mapped to the name.
If other account is blank, which csv2trn translates into "Imbalance", 
its name can often be inferred from the memo.
In this example, the output format is set to this module's CSV.
The Unix stream editor maps and infers.
Finally, csv2trn filters from its default input to its default output format:
this module's CSV to Ledger journal.

To merge multiple Ledger journals into one, see this module's program mrglent.

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
	outFormatName  string
	thisAccount    string
}

var (
	errOutFormatName = errors.New("invalid output format name")
	errThisAccount   = errors.New("cannot get this account: " +
		"CSV records do not contain that field and its flag is empty string")
)

func main() {
	log.SetPrefix("csv2trn: ")
	log.SetFlags(0)

	cfg := parseFlags()
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
		The number of fields in a record is checked by aft.ParseCSV,
		so disable the reader's check.
	*/
	r.FieldsPerRecord = -1
	r.ReuseRecord = true

	ts, err := parseCSVStatement(r, cfg, inFormat)
	if err != nil {
		log.Fatal(err)
	}

	stringTransactions(ts, os.Stdout, cfg.outFormatName)
}

/*
ParseFlags returns this program's configuration parsed from command line flags.
If help was requested, parseFlags writes this program's help text then exits.
If the flags are invalid, this program exits with a non-zero status.
*/
func parseFlags() config {
	var cfg config

	flag.StringVar(&cfg.currency, "c", "", fmt.Sprintf(
		"currency for transaction amounts: symbol %q or a code e.g. %q", "$", "GBP"))
	flag.StringVar(&cfg.formatFileName, "f", "",
		"file name containing input CSV record format in XML")
	flag.StringVar(&cfg.outFormatName, "o", aft.Ledger,
		fmt.Sprintf("output format name: %q or %q", aft.Ledger, aft.ModuleCSV))
	flag.StringVar(&cfg.thisAccount, "t", "", fmt.Sprintf(
		"this account name, which transactions belong to e.g. %q",
		"Assets:Current"))

	var help bool

	flag.BoolVar(&help, "h", false, "write this help text then exit")

	flag.Usage = usage
	flag.Parse()

	if help {
		usage()
		os.Exit(0)
	}

	return cfg
}

/*
ParseCSVStatement reads a CSV account statement,
parses a transaction's fields from those in the CSV record on each line
then returns the transactions.
If it fails to read the statement, parseCSVStatement returns an error.
If it fails to parse a transaction, parseCSVStatement logs a warning then continues.
*/
func parseCSVStatement(r *csv.Reader, cfg config, crf aft.CSVRecordFormat) ([]aft.Transaction, error) {
	var ts []aft.Transaction

	for {
		fs, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return ts, err
		}

		var t aft.Transaction

		t.Currency, t.ThisAccount = cfg.currency, cfg.thisAccount

		err = t.ParseCSV(fs, crf)
		if err != nil {
			n, _ := r.FieldPos(0)
			log.Printf("%v on line %v", err, n)

			continue
		}

		ts = append(ts, t)
	}

	return ts, nil
}

/*
StringTransactions writes the transactions in the named format.
It assumes the transactions are in date order ascending or descending.
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

// Usage writes the help text for this program.
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
