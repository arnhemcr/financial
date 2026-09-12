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
CSV2trn translates financial transactions from [CSV] records in an account statement
to formats including [Ledger] journal entries (lent).

It:
  - reads the statement from standard input
  - parses a transaction from each CSV record following the input format (default this module's CSV records (mcsv))
  - warns if a record cannot be parsed into a transaction to standard error
  - writes transactions to standard output in the other format ordered by date ascending

Usage:

	csv2trn [flags]

The flags are:

	-c string
		currency symbol or word e.g. "$" or "GBP"; takes precedence over currency field from input
	-f string
	 	name of file containing input CSV record format in XML
	-o string
	  	output format name: "lent" or "mcsv" (default "mcsv")
	-t string
		the Ledger name of this account e.g. "Assets:Current"; takes precedence over this account field from input

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

// The configuration returned by parseFlags.
type config struct {
	currency         string
	inFormatFileName string
	outFormatName    string
	thisAccount      string
}

func main() {
	var err error

	log.SetPrefix("csv2trn: ")
	log.SetFlags(0)

	cfg := parseFlags()

	if cfg.currency != "" && !aft.IsLedgerCurrency(cfg.currency) {
		log.Fatalf("expect currency symbol or word not %q", cfg.currency)
	}

	switch cfg.outFormatName {
	case aft.Ledger, aft.ModuleCSV:
		// This output format name is valid.
	default:
		log.Fatalf("expect output format name not %q", cfg.outFormatName)
	}

	inFormat := aft.NewModuleCSVRecordFormat()
	if cfg.inFormatFileName != "" {
		inFormat, err = aft.NewCSVRecordFormat(cfg.inFormatFileName)
		if err != nil {
			log.Fatal(err)
		}
	}

	if cfg.thisAccount == "" && inFormat.ThisAccountI == 0 {
		log.Fatal("expect this account from command line flag or field in input but neither are set")
	}

	r := csv.NewReader(os.Stdin)

	// The number of fields in a record is checked by this module; disable the reader's check.
	r.FieldsPerRecord, r.ReuseRecord = -1, true

	ts, err := parseCSVStatement(r, cfg, inFormat)
	if err != nil {
		log.Fatal(err)
	}

	stringTransactions(ts, os.Stdout, cfg.outFormatName)
}

/*
ParseFlags returns this program's configuration parsed from command line flags.
If the flags are invalid, this program exits with a non-zero status.
*/
func parseFlags() (c config) {
	flag.StringVar(&c.currency, "c", "",
		fmt.Sprintf("currency symbol or word e.g. %q or %q; takes precedence over currency field from input", "$", "GBP"))
	flag.StringVar(&c.inFormatFileName, "f", "", "name of file containing input CSV record format in XML")
	flag.StringVar(&c.outFormatName, "o", aft.ModuleCSV,
		fmt.Sprintf("output format name: %q or %q", aft.Ledger, aft.ModuleCSV))
	flag.StringVar(&c.thisAccount, "t", "", fmt.Sprintf(
		"the Ledger name of this account e.g. %q%s",
		"Assets:Current", "; takes precedence over this account field from input"))

	flag.Usage = usage
	flag.Parse()

	return c
}

/*
ParseCSVStatement reads a CSV account statement,
parses a transaction from each CSV record following the input format
then returns the transactions.
If it fails to read the statement, parseCSVStatement returns an error.
If it fails to parse a transaction, parseCSVStatement logs a warning then continues.
*/
func parseCSVStatement(r *csv.Reader, c config, f aft.CSVRecordFormat) (ts []aft.Transaction, err error) {
	for {
		fs, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return ts, err // This error will be wrapped by main.
		}

		t := aft.Transaction{Currency: c.currency, ThisAccount: c.thisAccount}

		err = t.ParseCSV(fs, f)
		if err != nil {
			n, _ := r.FieldPos(0)
			log.Printf("line %v: %v", n, err)

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
CSV2trn translates financial transactions from CSV records in an account statement
to other formats including Ledger journal entries (lent).

It:
 - reads the statement from standard input
 - parses a transaction from each CSV record following the input format (default this module's CSV records (mcsv))
 - warns if a record cannot be parsed into a transaction to standard error
 - writes transactions to standard output in the other format ordered by date ascending

Usage:

	csv2trn [flags]

The flags are:

`)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}
