/*
Copyright (C) 2025-2026 Andrew Flint.

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
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

const (
	Ledger = "lent" // The name of the Ledger journal entry format.

	/*
		The start and end lines for Ledger block comments
		(see the "Commenting Your Journal" section of the [Ledger 3 manual].

		[Ledger 3 manual]: https://ledger-cli.org/doc/ledger3.html
	*/
	StartBlockComment = "comment\n"
	EndBlockComment   = "end comment\n"

	// The start and end Ledger global comment lines around a mirror entry.
	StartMirrorEntry = "# mirror entry\n"
	EndMirrorEntry   = "# end mirror entry\n"
)

/*
IsLedgerIndented reports whether the line starts with a white space character
used by Ledger to indent postings and comments belonging to an entry.

See "Transactions and Comments" in the [Ledger 3 manual].
*/
func IsLedgerIndented(line string) bool {
	if len(line) == 0 {
		return false
	}

	switch rune(line[0]) {
	case ' ', '\t':
		return true
	default:
		return false
	}
}

/*
LoadLedgerAccountNames returns a list of Ledger account names loaded from the named XML file.
If it fails to load the list, LedgerAccounts returns the first error.

For example, file LedgerAccountsWithJournals.xml might contain three asset account names:

	<Accounts>
	  <Account>Assets:Current</Account>
	  <Account>Assets:Emergency</Account>
	  <Account>Assets:Savings</Account>
	</Accounts>
*/
func LoadLedgerAccountNames(fileName string) ([]string, error) {
	var as struct {
		Accounts []string `xml:"Account"`
	}

	bs, err := os.ReadFile(fileName)
	if err != nil {
		return as.Accounts, fmt.Errorf("LoadLedgerAccountNames: %w", err)
	}

	err = xml.Unmarshal(bs, &as)
	if err != nil {
		return as.Accounts, fmt.Errorf("LoadLedgerAccountNames: %w", err)
	}

	return as.Accounts, nil
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
		if !strings.HasPrefix(t.Code, startCode) {
			co += startCode
		}

		co += t.Code

		if !strings.HasSuffix(t.Code, endCode) {
			co += endCode
		}
	}

	return fmt.Sprintf("%v%v %v\n %v  %v\n %v\n",
		t.Date, co, t.Memo,
		t.ThisAccount, a,
		t.OtherAccount)
}

const (
	// The transaction code delimiters.
	startCode = "("
	endCode   = ")"
)
