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
	"time"
)

var ErrDateLayout = errors.New("expect Go-style date layout")

/*
ParseDate2 parses a date from the string, according to the layout, and returns the date in this module's layout.
It assumes the layout is valid, which can be checked with [ValidateDateLayout].
This module's date layout is YYYY-MM-DD or [ISO 8601 extended date],
which in Go is represented by string constant "2006-01-02" or [time.DateOnly].
If it fails to parse a date, ParseDate returns the first error.

[ISO 8601 extended date]: https://en.wikipedia.org/wiki/ISO_8601#Calendar_dates
*/
func ParseDate2(d, layout string) (string, error) {
	if d == "" {
		return "", fmt.Errorf("ParseDate2: %w not %q", errDate, d)
	}

	v, err := time.Parse(layout, d)
	if err != nil {
		return "", fmt.Errorf("ParseDate2: %w", err)
	}

	return v.Format(time.DateOnly), nil
}

// Deprecated: Call ParseDate2 instead of this function.
func ParseDate(d, layout string) (string, error) {
	t := trimDate(d, layout)

	v, err := time.Parse(layout, t)
	if err != nil {
		return "", fmt.Errorf("ParseDate: %w not %q", err, d)
	}

	return v.Format(time.DateOnly), nil
}

/*
ParseModuleDate2 parses a date from the string, according to this module's layout, and returns the date in that layout.
If it fails to parse a date, ParseModuleDate2 returns the error.
*/
func ParseModuleDate2(d string) (string, error) {
	return ParseDate2(d, time.DateOnly)
}

// Deprecated: Call ParseModuleDate2 instead of this function.
func ParseModuleDate(d string) (string, error) {
	return ParseDate(d, time.DateOnly)
}

/*
ValidateDateLayout returns nil if the date layout is valid.
If not, ValidateDateLayout returns the first error.
*/
func ValidateDateLayout(dl string) error {
	v, err := time.Parse(dl, dl)
	if err != nil {
		return fmt.Errorf("ValidateDateLayout: %w", err)
	}

	if v.Format(time.DateOnly) != time.DateOnly {
		return fmt.Errorf("ValidateDateLayout: %w not %q", ErrDateLayout, dl)
	}

	return nil
}

// Deprecated: Call ValidateDateLayout instead.
func IsDateLayout(layout string) bool {
	d, _ := time.Parse(layout, layout)

	return d.Format(time.DateOnly) == time.DateOnly
}

var errDate = errors.New("expect date")

/*
TrimDate assumes the text starts with a date and returns it trimmed to the length of the date layout.
If text or layout are empty string, trimDate returns empty string.
*/
func trimDate(text, layout string) string {
	tl, ll := len(text), len(layout)

	switch {
	case tl == 0, ll == 0:
		return ""
	case ll < tl:
		return text[0:ll]
	default:
		return text
	}
}
