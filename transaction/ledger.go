/*
Copyright (C) 2025 Andrew Flint.

This file is part of arnhemcr/financial/transaction.

Arnhemcr/financial/transaction is free software:
you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Arnhemcr/financial/transaction is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with arnhemcr/financial/transaction.
If not, see <https://www.gnu.org/licenses/>.
*/
package transaction

import (
	"fmt"
)

// StringLedger returns a string representing this transaction in [Ledger] format.
func (trn Transaction) StringLedger() string {
	const sep = "  " // The hard separator (see section 5.6 of [Ledger 3 manual]).

	return fmt.Sprintf("%v %v\n%v%v%v%.2f\n%v%v\n",
		trn.Date, trn.Memo,
		sep, trn.ThisAccount, sep, trn.Amount,
		sep, trn.OtherAccount)
}
