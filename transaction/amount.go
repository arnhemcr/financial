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
	"errors"
	"fmt"
	"strconv"
)

/*
ParseAmount returns the float64 value parsed from str and nil.
If it fails to parse a value, parseAmount returns an error.
*/
func parseAmount(str string) (float64, error) {
	amt, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0.00, fmt.Errorf("parseAmount: %w", err)
	}

	return amt, nil
}

var errNotPositive = errors.New("amount must be positive (0 < value)")

/*
ParsePositiveAmount returns the positive float64 value parsed from str and nil.
If it fails to parse a positive value, parsePositiveAmount returns the first error.
*/
func parsePositiveAmount(str string) (float64, error) {
	amt, err := parseAmount(str)
	if err != nil {
		return 0.00, err
	}

	if amt <= 0.00 {
		return 0.00, errNotPositive
	}

	return amt, nil
}

// FormatAmount returns a string representing a float64 value.
func formatAmount(amount float64) string {
	return strconv.FormatFloat(amount, 'f', 2, 64)
}
