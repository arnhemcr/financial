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
	"strconv"
)

/*
ParseAmount returns the floating-point number parsed from the string.
If it fails to parse a number, parseAmount returns the error.
*/
func parseAmount(s string) (float64, error) {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.00, fmt.Errorf("parseamount: %w", err)
	}

	return n, nil
}

var errNotPositive = errors.New("amount must be positive (0 < value)")

/*
ParsePositiveAmount returns the positive floating-point number parsed from the string.
If it fails to parse a positive number, parsePositiveAmount returns the first error.
*/
func parsePositiveAmount(s string) (float64, error) {
	n, err := parseAmount(s)
	if err != nil {
		return 0.00, err
	}

	if n <= 0.00 {
		return 0.00, errNotPositive
	}

	return n, nil
}

// StringAmount returns the floating-point number as a string.
func stringAmount(n float64) string {
	return strconv.FormatFloat(n, 'f', 2, 64)
}
