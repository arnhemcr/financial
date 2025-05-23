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
	"fmt"
)

const Ledger = "ledger" // The name of the Ledger journal entry format.

// StringLedger returns this transaction as a Ledger journal entry.
func (t Transaction) StringLedger() string {
	const s = "  " // The hard separator (see section 5.6 of the Ledger 3 manual).

	a := stringAmount(t.Amount)

	var c string

	if t.Currency != "" {
		c = " " + t.Currency
	}

	return fmt.Sprintf("%v %v\n%v%v%v%v%v\n%v%v\n",
		t.Date, t.Memo,
		s, t.ThisAccount, s, a, c,
		s, t.OtherAccount)
}
