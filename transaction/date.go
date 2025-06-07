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

import "time"

/*
ParseDate returns the date parsed from s, according to the date layout, as a string.
It assumes the layout is valid for Go e.g. "2006-01-02".
If parseDate fails to parse a date, it returns the error.
*/
func parseDate(s, layout string) (string, error) {
	d, err := time.Parse(layout, s)
	if err != nil {
		return "", err
	}

	return d.Format(time.DateOnly), nil
}
