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

package transaction

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"unicode"
)

const (
	Ledger = "lent" // The name of the Ledger journal entry format.

	/*
		The lines used to mark Ledger block comments (see [Commenting on your Journal] in the Ledger 3 manual).

		[Commenting on your Journal]: https://ledger-cli.org/doc/ledger3.html#Commenting-on-your-Journal
	*/
	StartBlockComment = "comment\n"
	EndBlockComment   = "end comment\n"

	// The Ledger global comments used to mark mirror entries by this module.
	StartMirrorEntry = "# mirror entry\n"
	EndMirrorEntry   = "# end mirror entry\n"
)

/*
IsLedgerCurrency reports whether the string is a currency symbol or word in Ledger
(see [Commodities and Currencies] in the Ledger 3 manual).

[Commodities and Currencies]: https://ledger-cli.org/doc/ledger3.html#Commodities-and-Currencies
*/
func IsLedgerCurrency(currency string) bool {
	for _, r := range currency {
		switch r {
		case '.', ',', '/', '@':
			return false
		}

		if unicode.IsDigit(r) || isLedgerSpace(r) {
			return false
		}
	}

	return true
}

/*
IsLedgerIndented reports whether the line is indented with a Ledger space rune
(see [Transactions and Comments] in the Ledger 3 manual).

[Transactions and Comments]: https://ledger-cli.org/doc/ledger3.html#Transactions-and-Comments
*/
func IsLedgerIndented(line string) bool {
	if line == "" {
		return false
	}

	return isLedgerSpace(rune(line[0]))
}

/*
LoadLedgerAccountNames returns a list of Ledger account names loaded from the named XML file.
If it fails to load the list, LedgerAccounts returns the first error.

For example, file LedgerAccountsWithJournals.xml contains three asset account names:

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

// StringLedger returns this transaction formated as a Ledger journal entry.
func (t Transaction) StringLedger() string {
	co := t.Code
	if co != "" {
		co = " " + startCode + co + endCode
	}

	a, cu := t.AmountText, t.Currency
	switch len(cu) {
	case 0:
		// There is no currency for this amount.
	case 1:
		a = cu + a
	default:
		a = a + " " + cu
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

var errCurrency = errors.New("expect currency symbol or word valid in Ledger")

// IsLedgerSpace reports whether the rune is a space in Ledger.
func isLedgerSpace(r rune) bool {
	switch r {
	case ' ', '\t': // The subset of Go space runes (see [unicode.IsSpace] that are space runes in Ledger 3.
		return true
	default:
		return false
	}
}
