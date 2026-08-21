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
	"time"
)

var ErrDateLayout = errors.New("expect Go-style date layout")

// Reports whether the string is a Go-style date layout.
func IsDateLayout(layout string) bool {
	d, _ := time.Parse(layout, layout)

	return d.Format(time.DateOnly) == time.DateOnly
}

/*
ParseDate parses a date from the string, according to the layout, and returns the date in this module's layout.
It assumes the layout is valid e.g. "2006-01-02", which can be verified by calling [IsDateLayout].
This module's date layout is YYYY-MM-DD also known as Go [time.DateOnly] and [ISO 8601 extended date].
If it fails to parse a date, parseDate returns the error.

[ISO 8601 extended date]: https://en.wikipedia.org/wiki/ISO_8601#Calendar_dates
*/
func ParseDate(text, layout string) (string, error) {
	t := trimDate(text, layout)

	d, err := time.Parse(layout, t)
	if err != nil {
		return "", err // This error will be wrapped by ParseCSV.
	}

	return d.Format(time.DateOnly), nil
}

/*
ParseDate parses a date from the string, according to this module's layout, and returns the date in that layout.
If it fails to parse a date, parseModuleDate returns the error.
*/
func ParseModuleDate(text string) (string, error) {
	return ParseDate(text, time.DateOnly)
}

/*
TrimDate assumes the text starts with a date and returns it trimmed to the length of the date layout.
If text or layout are empty string, trimDate returns empty string.
*/
func trimDate(text, layout string) string {
	tl, ll := len(text), len(layout)

	switch {
	case tl == 0 || ll == 0:
		return ""
	case ll < tl:
		return text[0:ll]
	default:
		return text
	}
}
