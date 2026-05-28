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
	"fmt"
	"time"
)

// Reports whether dl is a Go-style date layout.
func IsDateLayout(dl string) bool {
	d, _ := time.Parse(dl, dl)

	return d.Format(time.DateOnly) == time.DateOnly
}

/*
ParseDate returns the date parsed according to the layout from the start of text.
It assumes the layout is valid e.g. "2006-01-02".
The layout can be verified by calling function IsDateLayout.
If it fails to parse a date, parseDate returns the error.
*/
func ParseDate(text, layout string) (string, error) {
	d := trimDate(text, layout)

	date, err := time.Parse(layout, d)
	if err != nil {
		return "", fmt.Errorf("ParseDate: %w", err)
	}

	return date.Format(time.DateOnly), nil
}

/*
ParseModuleDate returns the date in this module's default layout from the start of text.
The layout is YYYY-MM-DD also known as [time.DateOnly] and [ISO 8601 extended date].
If it fails to parse a date, parseModuleDate returns the error.

[ISO 8601 extended date]: https://en.wikipedia.org/wiki/ISO_8601#Calendar_dates
*/
func ParseModuleDate(text string) (string, error) {
	return ParseDate(text, time.DateOnly)
}

// TrimDate returns the start of the text trimmed to the length of layout or the text whichever is shorter.
func trimDate(text, layout string) string {
	var (
		dlLen = len(layout)
		tLen  = len(text)
	)

	switch {
	case dlLen == 0 || tLen == 0:
		return ""
	case dlLen < tLen:
		return text[0:dlLen]
	default:
		return text
	}
}
