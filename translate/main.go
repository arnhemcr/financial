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
from a [comma-separated values (CSV)] format to [Ledger] format.

# Overview

Translate is a [filter] which:

  - reads a CSV account statement from standard input
  - parses a transaction's fields from those in the CSV record on each line
  - writes the transactions, in Ledger format ordered by date ascending, to standard output

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

func main() {
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)

	/*
		Format configures the parsing of a transaction's fields
		from those in a CSV record.

		It is currently set from arnhemcr/financial format:

			Date,This account,Other account,Memo,Amount
			2025-05-15,01-2345-6789012-34,98-7654321-09,Collective Coffee,-6.50
	*/
	var format = aft.CSVformat{
		NFields: 5, // The number of fields in the record.

		/*
			The indexes of fields in the record.
			Index values are in the range zero to NFields inclusive.
			A value of zero means this record does not contain that field.
		*/
		DateI:         1,
		ThisAccountI:  2,
		OtherAccountI: 3,
		MemoI:         4,
		AmountI:       5,
		/*
			If the record contains credit and debit fields instead of amount,
			set the credit and debit indexes to non-zero values and
			the amount index to zero.
		*/
		CreditI: 0, DebitI: 0,

		ThisAccount: "", // The name of this account.

		DateFormat: "2006-01-02", // The format of a dates in the record.
	}

	var help bool

	help, format.ThisAccount = parseFlags()
	if help {
		usage()
		os.Exit(0)
	}

	err := format.Validate()
	if err != nil {
		log.Fatalf("validate: %v", err)
	}

	rdr := csv.NewReader(os.Stdin)

	/*
		Configure how a CSV record with an unexpected number of fields is handled.
		Set to strict, translate will print an error and exit.
		Set to relaxed, it will print an warning but continue.
	*/
	rdr.FieldsPerRecord = int(format.NFields) // strict
	// rdr.FieldsPerRecord = -1 // relaxed

	trns, err := parseTransactions(rdr, format)
	if err != nil {
		log.Fatalf("parsetransactions: %v", err)
	}

	stringTransactions(trns, os.Stdout)
}

// ParseFlags returns the values of flags from the command that ran translate.
func parseFlags() (bool, string) {
	var help bool

	flag.BoolVar(&help, "h", false, "write this help text then exit")

	var thisAccount string

	flag.StringVar(&thisAccount, "t", "", "this account name")

	flag.Usage = usage
	flag.Parse()

	return help, thisAccount
}

/*
ParseTransactions reads a CSV account statement,
parses a transaction's fields from those in the CSV record on each line
then returns the transactions.
If it fails to read the statement, parseTransactions returns an error.
*/
func parseTransactions(reader *csv.Reader, format aft.CSVformat) ([]aft.Transaction, error) {
	var trns []aft.Transaction

	for {
		flds, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return trns, fmt.Errorf("read: %w", err)
		}

		var trn aft.Transaction

		err = trn.ParseCSV(flds, format)
		if err != nil {
			lineI, _ := reader.FieldPos(0)
			log.Printf("parsecsv: %v on line %v", err, lineI)

			continue
		}

		trns = append(trns, trn)
	}

	return trns, nil
}

// StringTransactions writes the transactions ordered by date ascending.
func stringTransactions(trns []aft.Transaction, writer *os.File) {
	nTrns := len(trns)

	trnSeq := slices.All(trns)
	if 2 <= nTrns && trns[0].Date > trns[nTrns-1].Date {
		trnSeq = slices.Backward(trns)
	}

	for _, trn := range trnSeq {
		fmt.Fprint(writer, trn.StringLedger())
	}
}

// Usage prints the help text for translate.
func usage() {
	fmt.Fprintf(os.Stderr, "usage: %v [flags]\n", os.Args[0])
	fmt.Fprint(os.Stderr, `
This program translates financial transactions 
from a comma-separated values (CSV) format to Ledger format.
It is a filter which:

  - reads a CSV account statement, belonging to this account, from standard input
  - parses a transaction's fields from those in the CSV record on each line
  - writes the transactions, in Ledger format ordered by date ascending, to standard output

The parsing configuration is built into the format structure in function main.
It can be updated to translate other CSV formats.
The flags are:

`)
	flag.PrintDefaults()
}
