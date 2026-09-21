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
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

var (
	errAmountMismatch  = errors.New("expect amount text and number to match")
	errCreditDebit     = errors.New("expect string in credit and empty string in debit or visa versa")
	errDecimal         = errors.New("expect decimal number string")
	errDecimalZero     = errors.New("expect non-zero decimal number")
	errPositiveDecimal = errors.New("expect positive decimal number as credit or debit")
)

/*
IsDecimal reports whether the string represents a decimal floating-point number
without an exponent but optionally with a sign.
Adapted from [Go's floating-point literals], the syntax is:

	signed_decimal_float_lit = [ sign ] ( decimal_digits "." [ decimal_digits ] | "." decimal_digits ) .
	sign = ( "+" | "-" ) .

[Go's floating-point literals]: https://go.dev/ref/spec#Floating-point_literals
*/
func isDecimal(d string) bool {
	var hasDigit, hasPoint bool

	for i, r := range d {
		switch {
		case unicode.IsDigit(r):
			hasDigit = true
		case i == 0 && (r == '-' || r == '+'):
			// This string starts with a sign.
		case !hasPoint && r == '.':
			hasPoint = true
		default:
			return false
		}
	}

	return hasDigit
}

/*
ParseAmount returns the amount of a transaction as both string and floating-point values.
It is parsed from the amount, credit or debit field of a CSV record according to the format.
ParseAmount assumes the format is valid (see [CSVRecordFormat.Validate]).
The value string must represent a decimal floating-point number with an optional sign but without an exponent.
The amount cannot be zero.
If parseAmount fails to parse a non-zero amount, it returns the error.
*/
func parseAmount(fields []string, f CSVRecordFormat) (vText string, v float64, err error) {
	var negative bool

	if f.AmountI == 0 {
		// If field index for amount is zero then indexes for credit and debit will be non-zero.
		c, d := fields[f.CreditI], fields[f.DebitI]

		switch {
		case c != "" && d == "":
			vText, v, err = parsePositiveDecimal(c)
		case d != "" && c == "":
			vText, v, err = parsePositiveDecimal(d)
			negative = true
		default:
			err = fmt.Errorf("%w not %q and %q", errCreditDebit, c, d)
		}
	} else {
		vText, v, err = parseNonZeroDecimal(fields[f.AmountI])
	}

	if err != nil {
		return "", 0, err // This error will be wrapped by ParseCSV.
	}

	if negative {
		vText = "-" + vText
		v *= -1
	}

	return vText, v, nil
}

/*
ParseNonZeroDecimal returns the non-zero decimal number parsed from the string as both string and floating-point values.
If it fails to parse a non-zero number, parseNonZeroDecimal returns the first error.
*/
func parseNonZeroDecimal(d string) (nText string, n float64, err error) {
	if !isDecimal(d) {
		return "", 0, fmt.Errorf("%w not %q", errDecimal, d)
	}

	n, err = strconv.ParseFloat(d, 64)
	if err != nil {
		return "", 0, err // This error will be wrapped by ParseCSV.
	}

	if n == 0 {
		return "", 0, fmt.Errorf("%w not %v", errDecimalZero, d)
	}

	return d, n, nil
}

/*
ParsePositiveDecimal returns the positive decimal number parsed from the string
as both string and floating-point values.
If it fails to parse a positive number, parsePositiveDecimal returns the first error.
*/
func parsePositiveDecimal(d string) (nText string, n float64, err error) {
	nText, n, err = parseNonZeroDecimal(d)
	if err != nil {
		return "", 0, err
	}

	if n < 0 {
		return "", 0, fmt.Errorf("%w not %v", errPositiveDecimal, n)
	}

	return nText, n, nil
}
