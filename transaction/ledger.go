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

package transaction

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"
)

const (
	Ledger = "ledger" // The name of the Ledger journal entry format.

	// The transaction code delimiters.
	StartCode = "("
	EndCode   = ")"
)

/*
ParseLedger parses this transaction's date, code and memo fields
from the first line of a Ledger journal entry.
If it fails to parse those fields, ParseLedger returns the first error.

The format of a Ledger journal entry is described in the
"Transactions and Comments" section of the [Ledger 3 manual].

[Ledger 3 manual]: https://ledger-cli.org/doc/ledger3.html
*/
func (t *Transaction) ParseLedger(lines []string) error {
	const fName = "t.ParseLedger"

	if len(lines) == 0 || !startsDigit(lines[0]) {
		return fmt.Errorf("%v: %v:", fName, errEntryStart)
	}

	fs := strings.Fields(lines[0])

	var err error

	t.Date, err = getDate(fs[0])
	if err != nil {
		return fmt.Errorf("%v: %v:", fName, err)
	}

	n := len(fs)

	var i int

	switch {
	case 3 <= n && isStatusMark(fs[1]):
		i = 2
	case 2 <= n:
		i = 1
	default:
		return nil
	}

	t.Code = getCode(fs[i])
	if t.Code != "" {
		i++
	}

	if i < n {
		t.Memo = strings.Join(fs[i:n], " ")
	}

	return nil
}

// StringLedger returns this transaction as a Ledger journal entry.
func (t Transaction) StringLedger() string {
	a := stringAmount(t.Amount)

	cu := t.Currency
	switch cu {
	case "":
		// There is no currency for the amount.
	case "$":
		a = "$" + a
	default:
		a = a + " " + cu
	}

	var co string

	if t.Code != "" {
		co = " "
		if !strings.HasPrefix(t.Code, StartCode) {
			co += StartCode
		}

		co += t.Code

		if !strings.HasSuffix(t.Code, EndCode) {
			co += EndCode
		}
	}

	return fmt.Sprintf("%v%v %v\n %v  %v\n %v\n",
		t.Date, co, t.Memo,
		t.ThisAccount, a,
		t.OtherAccount)
}

var (
	errEntryStart = errors.New("first line of Ledger journal entry " +
		"must start with decimal digit")
)

/*
GetCode returns a transaction code from the Ledger journal entry string.
In an entry, a transaction code is in brackets e.g. code "MT" appears in the entry as "(MT)".
If the string does not contain the code, getCode returns an empty string.
*/
func getCode(code string) string {
	n := len(code)
	if n < 3 ||
		!strings.HasPrefix(code, StartCode) || !strings.HasSuffix(code, EndCode) {
		return ""
	}

	return code[1 : n-1]
}

/*
GetDate returns the actual date from a Ledger journal entry dates string.
The dates string syntax is "actual[=effective]".
If the actual string does not contain a date, getDate returns an error.
*/
func getDate(dates string) (string, error) {
	const separator = "="

	d, _, _ := strings.Cut(dates, separator)

	return ParseDate(d, time.DateOnly)
}

// IsStatusMark reports whether the string is a Ledger journal entry status mark.
func isStatusMark(mark string) bool {
	const cleared, pending = "*", "!"

	switch mark {
	case cleared, pending:
		return true
	default:
		return false
	}
}

// StartsDigit reports whether the string starts with a decimal digit.
func startsDigit(line string) bool {
	if 0 < len(line) && unicode.IsDigit(rune(line[0])) {
		return true
	}

	return false
}
