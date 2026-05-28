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
CSV2trn [filters] financial transactions from [comma-separated values (CSV)] records
in an account statement to a selected format.

A transaction is the transfer of an amount of currency between accounts on a particular day.
It is described by a memo and code, also called the description and transaction type.
A statement and its records belong to an account, which is called this account in transactions from those records.

CSV2trn reads a statement from standard input.
It parses each line as a transaction CSV record following an input format
and warns about lines that cannot be parsed to standard error.
The input format defaults to this module's CSV record (mcsv), but it is usually loaded from an XML file.
In XML, the mcsv format is:

	<CSVRecordFormat>
	    <NFields>7</NFields><!-- The number of fields in the record. -->

	    <!-- The index of each field: 1..NFields or zero if it is not provided. -->
	    <DateI>1</DateI>
	        <DateLayout>2006-01-02</DateLayout><!-- The default Go date layout time.DateOnly. -->
	    <ThisAccountI>2</ThisAccountI>
	    <OtherAccountI>3</OtherAccountI>
	    <CodeI>4</CodeI>
	    <MemoI>5</MemoI>
	    <AmountI>6</AmountI>
	        <CreditI>0</CreditI>
	        <DebitI>0</DebitI>
	    <CurrencyI>7</CurrencyI>
	</CSVRecordFormat>

If the other account field is not provided then its default value is "Imbalance".

CSV2trn orders transactions by date ascending and writes them to standard output in the selected format:
[Ledger] journal entries (lent) or mcsv.

Usage:

	csv2trn [flags]

The flags are:

	-c string
	  	currency: symbol "$" or code e.g. "GBP"; overrides currency field from input
	-f string
	 	name of file containing input CSV record format in XML
	-h	write this help text then exit
	-o string
	  	output format name: Ledger journal entry "lent" or "mcsv" (default "mcsv")
	-t string
	  	the Ledger name of this account e.g. "Assets:Current"; overrides this account field from input

See also [this package's README].

[comma-separated values (CSV)]: https://en.wikipedia.org/wiki/Comma-separated_values
[filters]: https://en.wikipedia.org/wiki/Filter_(software)
[Ledger]: https://ledger-cli.org
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
		"CSV records do not contain that field and its flag is not set")
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

	inFormat := aft.NewModuleCSVRecordFormat()
	if cfg.formatFileName != "" {
		inFormat, err = aft.NewCSVRecordFormat(cfg.formatFileName)
		if err != nil {
			log.Fatal(err)
		}
	}

	if cfg.thisAccount == "" && inFormat.ThisAccountI == 0 {
		log.Fatal(errThisAccount)
	}

	r := csv.NewReader(os.Stdin)
	/*
		The number of fields in a record is checked by aft.ParseCSV,
		so disable the reader's check.
	*/
	r.FieldsPerRecord, r.ReuseRecord = -1, true

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

	flag.StringVar(&cfg.currency, "c", "", fmt.Sprintf("currency: symbol %q or code e.g. %q%s",
		"$", "GBP", "; overrides currency field from input"))
	flag.StringVar(&cfg.formatFileName, "f", "", "name of file containing input CSV record format in XML")
	flag.StringVar(&cfg.outFormatName, "o", aft.ModuleCSV,
		fmt.Sprintf("output format name: Ledger journal entry %q or %q",
			aft.Ledger, aft.ModuleCSV))
	flag.StringVar(&cfg.thisAccount, "t", "", fmt.Sprintf(
		"the Ledger name of this account e.g. %q%s",
		"Assets:Current", "; overrides this account field from input"))

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
			return ts, fmt.Errorf("parseCSVStatement: %w", err)
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
CSV2trn filters financial transactions from comma-separated values (CSV) records in an account statement
to a selected format.

A transaction is the transfer of an amount of currency between accounts on a particular day.
It is described by a memo and code, also called the description and transaction type.
A statement and its records belong to an account, which is called this account in transactions from those records.

CSV2trn reads a statement from standard input.
It parses each line as a transaction CSV record following an input format
and warns about lines that cannot be parsed to standard error.
The input format defaults to this module's CSV record (mcsv), but it is usually loaded from an XML file.
In XML, the mcsv format is:

    <CSVRecordFormat>
        <NFields>7</NFields><!-- The number of fields in the record. -->

        <!-- The index of each field: 1..NFields or zero if it is not provided. -->
        <DateI>1</DateI>
	    <DateLayout>2006-01-02</DateLayout><!-- The default Go date layout time.DateOnly. -->
        <ThisAccountI>2</ThisAccountI>
        <OtherAccountI>3</OtherAccountI>
        <CodeI>4</CodeI>
        <MemoI>5</MemoI>
        <AmountI>6</AmountI>
            <CreditI>0</CreditI>
            <DebitI>0</DebitI>
        <CurrencyI>7</CurrencyI>
    </CSVRecordFormat>

If the other account field is not provided then its default value is
"Imbalance".

CSV2trn orders transactions by date ascending and writes them to standard output
in the selected format: Ledger journal entries (lent) or mcsv.

Usage:

	csv2trn [flags]

The flags are:

`)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}
