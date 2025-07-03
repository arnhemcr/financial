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
ParseLedger parses the date and optional code for this transaction
from the lines of a Ledger journal entry.
If it fails to parse those fields, ParseLedger returns the first error.

The format of a Ledger journal entry is described in the
"Transactions and Comments" section of the [Ledger 3 manual].

[Ledger 3 manual]: https://ledger-cli.org/doc/ledger3.html
*/
func (t *Transaction) ParseLedger(lines []string) error {
	if len(lines) == 0 || len(lines[0]) == 0 || !unicode.IsDigit(rune(lines[0][0])) {
		return errStartNumber
	}

	fs := strings.Fields(lines[0])

	var err error

	t.Date, err = getDate(fs[0])
	if err != nil {
		return err
	}

	n := len(fs)

	switch {
	case 3 <= n && isStatusMark(fs[1]):
		t.Code = getCode(fs[2])
	case 2 <= n:
		t.Code = getCode(fs[1])
	}

	return nil
}

// StringLedger returns this transaction as a Ledger journal entry.
func (t Transaction) StringLedger() string {
	const sp, sp2 = " ", "  " // single and double space

	a := stringAmount(t.Amount)

	cu := t.Currency
	switch cu {
	case "":
		// There is no currency to associate with the amount.
	case "$":
		a = "$" + a
	default:
		a = a + sp + cu
	}

	var co string

	if t.Code != "" {
		co = sp
		if !strings.HasPrefix(t.Code, StartCode) {
			co += StartCode
		}

		co += t.Code

		if !strings.HasSuffix(t.Code, EndCode) {
			co += EndCode
		}
	}

	return fmt.Sprintf("%v%v %v\n%v%v%v%v\n%v%v\n",
		t.Date, co, t.Memo,
		sp, t.ThisAccount, sp2, a,
		sp, t.OtherAccount)
}

// The first line of a Ledger journal entry must start with a number.
var (
	errStartNumber = errors.New("trn.parseledger: line must start with number")
)

/*
GetCode returns a transaction code from a Ledger journal entry string.
A transaction code appears in brackets in an entry e.g. "MT" appears as "(MT)".
If the string does not contain a code, getCode returns empty string.
*/
func getCode(s string) string {
	n := len(s)
	if n < 3 || !strings.HasPrefix(s, StartCode) || !strings.HasSuffix(s, EndCode) {
		return ""
	}

	return s[1 : n-1]
}

/*
GetDate returns the actual date from a Ledger journal entry date string.
The date string syntax is "actual[=effective]".
If the actual string does not contain a date, getDate returns the error.
*/
func getDate(s string) (string, error) {
	const sep = "="

	d, _, _ := strings.Cut(s, sep)

	return ParseDate(d, time.DateOnly)
}

// IsStatusMark reports whether the string is a Ledger journal entry status mark.
func isStatusMark(s string) bool {
	const cleared, pending = "*", "!" // transaction status marks

	switch s {
	case cleared, pending:
		return true
	default:
		return false
	}
}
