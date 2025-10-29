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
	"strings"
)

var (
	errAmount         = errors.New("amount cannot be zero")
	errCreditDebit    = errors.New("credit and debit cannot both be empty string or both non-empty string")
	errPositiveNumber = errors.New("number must be positive")
)

/*
DepunctAmount removes any punctuation bytes from an amount so it can be parsed as a value.
It returns the depunctuated amount.

With punctuation bytes "$," and amount "$5,432.10",
depunctAmount returns "5342.10 from which value 5432.10 can be parsed.
*/
func depunctAmount(amount, puncts string) string {
	if puncts == "" {
		return amount
	}

	bs, cs := []byte(amount), []byte{}
	for _, b := range bs {
		if 0 <= strings.IndexByte(puncts, b) {
			continue
		}

		cs = append(cs, b)
	}

	return string(cs)
}

/*
ParseAmount parses the value of a transaction
from either the amount, credit or debit fields.
The value cannot be zero.
If it fails to parse a non-zero value, parseAmount returns the first error.
*/
func parseAmount(fields []string, crf CSVRecordFormat) (float64, error) {
	a, c, d := fields[crf.AmountI], fields[crf.CreditI], fields[crf.DebitI]

	var (
		v   float64
		err error
	)

	ps := crf.AmountPuncts

	switch {
	case a != "":
		v, err = parseFloat(depunctAmount(a, ps))
	case c != "" && d == "":
		v, err = parsePositiveFloat(depunctAmount(c, ps))
	case d != "" && c == "":
		v, err = parsePositiveFloat(depunctAmount(d, ps))

		v *= -1
	default:
		return 0, errCreditDebit
	}

	switch {
	case err != nil:
		return 0, err
	case v == 0:
		return 0, errAmount
	default:
		return v, nil
	}
}

/*
ParseFloat returns the floating-point number parsed from the string.
If it fails to parse a number, parseFloat returns the error.
*/
func parseFloat(s string) (float64, error) {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse number: %w", err)
	}

	return n, nil
}

/*
ParsePositiveFloat returns the positive floating-point number parsed from the string.
If it fails to parse a positive number, parsePositiveFloat returns the first error.
*/
func parsePositiveFloat(s string) (float64, error) {
	n, err := parseFloat(s)

	switch {
	case err != nil:
		return 0, err
	case n <= 0:
		return 0, errPositiveNumber
	default:
		return n, nil
	}
}

// StringAmount returns the floating-point number as a string.
func stringAmount(n float64) string {
	return strconv.FormatFloat(n, 'f', -1, 64)
}
