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
Translate translates financial transactions
from a [comma-separated values (CSV)] format
to [Ledger] journal entries or arnhemcr/financial CSV records.

# Overview

Translate is a [filter] which:

  - reads a CSV account statement from standard input
  - parses a transaction's fields from those in the CSV record on each line
  - writes the transactions, by date ascending, to standard output

# Transactions and accounts

A financial transaction on an account statement is the transfer of money (an amount)
between the account the statement belongs to (this account)
and another account (other account) on a particular date.

The CSV records in a statement may contain account numbers e.g. 01-2345-6789012-34.
Financial software, including Ledger, identifies account by name e.g. Assets:Current.

This account is mandatory.
Its value can be set from a field in the CSV records,
and can be overridden from the translate command line (-t flag).
Other account is optional, and it has a default value.

# Parsing and configuration

The parsing configuration is built into the format structure in function main.
It can be updated to translate other CSV formats.

# References

[comma-separated values (CSV)]: https://www.ietf.org/rfc/rfc4180.txt
[filter]: https://en.wikipedia.org/wiki/Filter_(software)
[Ledger]: https://ledger-cli.org "Ledger command-line accounting"
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
	help          bool
	outFormatName string
	thisAccount   string
}

func main() {
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)

	inFormat := aft.GetAFFormat()

	cfg := parseFlags()
	if cfg.help {
		usage()
		os.Exit(0)
	}

	err := inFormat.Validate()
	if err != nil {
		log.Fatalf("validate: %v", err)
	}

	r := csv.NewReader(os.Stdin)

	/*
		Configure how a CSV record with an unexpected number of fields is handled.
		Set to strict, translate will print an error and exit.
		Set to relaxed, it will print an warning but continue.
	*/
	r.FieldsPerRecord = int(inFormat.NFields) // strict
	// r.FieldsPerRecord = -1 // relaxed

	ts, err := parseTransactions(r, inFormat)
	if err != nil {
		log.Fatalf("parsetransactions: %v", err)
	}

	stringTransactions(ts, os.Stdout, cfg.outFormatName)
}

// ParseFlags returns the values of flags from the command that ran translate.
func parseFlags() config {
	var cfg config

	flag.BoolVar(&cfg.help, "h", false, "write this help text then exit")
	flag.StringVar(&cfg.outFormatName, "o", aft.Ledger,
		fmt.Sprintf("output format name: %q or %q", aft.CSV, aft.Ledger))
	flag.StringVar(&cfg.thisAccount, "t", "", "this account name")

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
func parseTransactions(r *csv.Reader, cf aft.CSVFormat) ([]aft.Transaction, error) {
	var ts []aft.Transaction

	for {
		fs, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return ts, fmt.Errorf("read: %w", err)
		}

		var t aft.Transaction

		err = t.ParseCSV(fs, cf)
		if err != nil {
			i, _ := r.FieldPos(0)
			log.Printf("parsecsv: %v on line %v", err, i)

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
	fmt.Fprintf(os.Stderr, "usage: %v [flags]\n", os.Args[0])
	fmt.Fprint(os.Stderr, `
This program translates financial transactions 
from a comma-separated values (CSV) format
to Ledger journal entries or arnhemcr/financial CSV records.
It is a filter which:

  - reads a CSV account statement, belonging to this account, from standard input
  - parses a transaction's fields from those in the CSV record on each line
  - writes the transactions, ordered by date ascending, to standard output

The parsing configuration is built into the format structure in function main.
It can be updated to translate other CSV formats.
The flags are:

`)
	flag.PrintDefaults()
}
