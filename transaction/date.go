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
	"fmt"
	"time"
)

/*
ParseModuleDate returns the date in this module's layout from the start of text.
This module's date layout is "2006-01-02" or time.DateOnly.
If it fails to parse a date, parseModuleDate returns the error.
*/
func ParseModuleDate(text string) (string, error) {
	const dLen = len(time.DateOnly)

	var (
		tLen = len(text)
		d    string
	)

	if dLen <= tLen {
		d = text[0:dLen]
	} else {
		d = text[0:tLen]
	}

	return parseDate(d, time.DateOnly)
}

/*
ParseDate returns the date parsed from text according to the layout.
It assumes the layout is valid e.g. "2006-01-02".
If it fails to parse a date, parseDate returns the error.
*/
func parseDate(text, layout string) (string, error) {
	date, err := time.Parse(layout, text)
	if err != nil {
		return "", fmt.Errorf("parseDate: %w", err)
	}

	return date.Format(time.DateOnly), nil
}
