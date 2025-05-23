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

func main() {
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)

	cfg := parseFlags()
	if cfg.help {
		usage()
		os.Exit(0)
	}

	var err error

	inFormat := aft.GetPkgFormat()
	if cfg.formatFileName != "" {
		inFormat, err = aft.GetFormat(cfg.formatFileName)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = inFormat.Validate()
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(os.Stdin)

	/*
		Configure how a CSV record with an unexpected number of fields is handled.
		Set to strict, translate will print an error and exit.
		Set to relaxed, it will print an warning but continue.
	*/
	r.FieldsPerRecord = int(inFormat.NFields) // strict
	// r.FieldsPerRecord = -1 // relaxed

	ts, err := parseTransactions(r, cfg.thisAccount, inFormat)
	if err != nil {
		log.Fatal(err)
	}

	stringTransactions(ts, os.Stdout, cfg.outFormatName)
}

// ParseFlags returns the values of flags from the command that ran translate.
func parseFlags() config {
	var cfg config

	flag.StringVar(&cfg.formatFileName, "f", "", "format file name")
	flag.BoolVar(&cfg.help, "h", false, "write this help text then exit")
	flag.StringVar(&cfg.outFormatName, "o", aft.Ledger,
		fmt.Sprintf("output format name: this module's %q or %q", aft.CSV, aft.Ledger))
	flag.StringVar(&cfg.thisAccount, "t", "",
		"this account name, the name of the account that this statement belongs to")

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
func parseTransactions(r *csv.Reader, thisAccount string, cf aft.CSVFormat) ([]aft.Transaction, error) {
	var ts []aft.Transaction

	for {
		fs, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return ts, fmt.Errorf("read: %w", err)
		}

		var t aft.Transaction

		t.ThisAccount = thisAccount

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
	flag.PrintDefaults()
}
