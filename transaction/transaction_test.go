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
	"strings"
	"testing"
)

func TestHappyConfig(t *testing.T) {
	t.Parallel()

	// test the CSV formats
	for _, format := range [...]CSVFormat{kbFull, mini, pcu} {
		err := format.Validate()
		if err != nil {
			t.Fatalf("wrong format.Validate: expected==nil, got==%v\n", err)
		}

		if !format.IsValid() {
			t.Fatalf("wrong format.IsValid: expected==true, got==false\n")
		}
	}
}

func TestHappyLedger(t *testing.T) {
	t.Parallel()

	trn0 := Transaction{
		Date:        "2025-05-05",
		ThisAccount: "ABC", OtherAccount: "XYZ",
		Memo:   "Transfer",
		Amount: -1.23,
	}

	expect := "2025-05-05 Transfer\n  ABC  -1.23\n  XYZ\n"

	got := trn0.StringFormat("ledger")
	if got != expect {
		t.Fatalf("wrong StringFormat: expected==%q, got==%q\n", expect, got)
	}
}

func TestHappyTransactKBAmount(t *testing.T) {
	t.Parallel()

	format := kbFull

	// test credit transaction with amount field
	flds := []string{"ZZ-YYYY-XXXXXXX-WW", "29-12-2023", "Automatic Payment Rates MISS E MACD ;Ref: Rates MISS E MACD",
		"AP", "Rates", "E", "", "", "", "", "MISS E MACD", "AA-BBBB-CCCCCCC-DD", "162.00", "", "162.00", "1434.23"}

	var trn Transaction

	trn.ThisAccount = "Assets:Current:KB06"

	err := trn.ParseCSV(flds, format)
	if err != nil {
		t.Fatalf("wrong trn.ParseCSV: expected==nil, got==%v", err)
	}

	expectAmount := 162.00
	gotAmount := trn.Amount

	if gotAmount != expectAmount {
		t.Fatalf("wrong amount: expected==%v, got==%v\n", expectAmount, gotAmount)
	}

	expect := "2023-12-29"
	got := trn.Date

	if got != expect {
		t.Fatalf("wrong date: expected==%q, got==%q\n", expect, got)
	}

	expect = "Automatic Payment Rates MISS E MACD ;Ref: Rates MISS E MACD"
	got = trn.Memo

	if got != expect {
		t.Fatalf("wrong memo: expected==%q, got==%q\n", expect, got)
	}

	expect = "AA-BBBB-CCCCCCC-DD"
	got = trn.OtherAccount

	if got != expect {
		t.Fatalf("wrong that account: expected==%q, got==%q\n", expect, got)
	}

	expect = "Assets:Current:KB06"
	got = trn.ThisAccount

	if got != expect {
		t.Fatalf("wrong this account: expected==%q, got==%q\n", expect, got)
	}
}

func TestHappyTransactMini(t *testing.T) {
	t.Parallel()

	format := mini

	// test transaction with minimal number of fields
	flds := []string{"2025-04-17", "A penny for your thoughts.", ".01"}

	var trn Transaction

	trn.ThisAccount = "Mini"

	err := trn.ParseCSV(flds, format)
	if err != nil {
		t.Fatalf("wrong trn.ParseCSV: expected==nil, got==%v", err)
	}

	if !trn.IsValid() {
		t.Fatalf("wrong trn.IsValid: expected==true, got==false\n")
	}

	expect := "2025-04-17,Mini,Imbalance,A penny for your thoughts.,0.01,\n"
	got := trn.StringFormat("csv")

	if got != expect {
		t.Fatalf("wrong StringFormat: expected==%q, got==%q\n", expect, got)
	}
}

func TestHappyTransactPCUCredit(t *testing.T) {
	t.Parallel()

	format := pcu

	// test credit transaction with credit field
	flds := []string{"28/11/2019", "HealthAndLif eInsuranceAn dSubs ARNHEMCR BP", "", "123.00", "316.69"}

	var trn Transaction

	trn.ThisAccount = "Assets:Current:PCUS1"

	err := trn.ParseCSV(flds, format)
	if err != nil {
		t.Fatalf("wrong trn.ParseCSV: expected==nil, got==%v", err)
	}

	expect := "2019-11-28,Assets:Current:PCUS1,Imbalance,HealthAndLif eInsuranceAn dSubs ARNHEMCR BP,123.00,\n"
	got := trn.StringFormat("csv")

	if got != expect {
		t.Fatalf("wrong StringFormat: expected==%q, got==%q\n", expect, got)
	}
}

func TestHappyTransactPCUDebit(t *testing.T) {
	t.Parallel()

	format := pcu

	// test debit transaction with debit field
	flds := []string{"07/01/2020", "554PHP 18832946 Best of Health", "16.92", "", "265.01"}

	var trn Transaction

	trn.ThisAccount = "Assets:Current:PCUS1"

	err := trn.ParseCSV(flds, format)
	if err != nil {
		t.Fatalf("wrong trn.ParseCSV: expected==nil, got==%v", err)
	}

	expect := "2020-01-07,Assets:Current:PCUS1,Imbalance,554PHP 18832946 Best of Health,-16.92,\n"
	got := trn.StringFormat("csv")

	if got != expect {
		t.Fatalf("wrong StringFormat: expected==%q, got==%q\n", expect, got)
	}
}

func TestHappyTransactOutIn(t *testing.T) {
	t.Parallel()

	trn0 := Transaction{
		Date:        "2025-05-02",
		ThisAccount: "Assets:Current:KB05", OtherAccount: "Assets:Current:PCUS1",
		Memo:   "To term deposit",
		Amount: 42.00,
	}
	if !trn0.IsValid() {
		t.Fatalf("wrong trn0.IsValid: expected==true, got==false\n")
	}

	// format then parse a transaction using the standard CSV format
	expect := "2025-05-02,Assets:Current:KB05,Assets:Current:PCUS1,To term deposit,42.00,\n"
	got := trn0.StringFormat("csv")

	if got != expect {
		t.Fatalf("wrong trn0.StringFormat: expected==%q, got==%q\n", expect, got)
	}

	var trn1 Transaction

	got = strings.TrimRight(got, "\n")

	err := trn1.ParseCSV(strings.Split(got, ","), GetPkgFormat())
	if err != nil {
		t.Fatalf("wrong trn1.ParseCSV: expected==nil, got==%v", err)
	}
}

func TestUnhappyConfigIndexes(t *testing.T) {
	t.Parallel()

	format := kbFull

	// field index cannot be out of range
	format.AmountI = format.NFields + 1

	err := format.Validate()
	if err == nil {
		t.Fatalf("wrong format.Validate: expected!=nil, got==nil\n")
	}

	format = kbFull

	// field indexes cannot share a non-zero value
	format.CreditI, format.DebitI = 1, 1

	if format.IsValid() {
		t.Fatalf("wrong format.IsValid: expected==false, got==true\n")
	}
}

func TestUnhappyConfigMandatory(t *testing.T) {
	t.Parallel()

	format := kbFull

	// date field index cannot be zero
	format.DateI = 0

	err := format.Validate()
	if err == nil {
		t.Fatalf("wrong format.Validate: expected!=nil, got==nil\n")
	}

	format = kbFull

	// date format cannot be empty string
	format.DateLayout = ""

	err = format.Validate()
	if err == nil {
		t.Fatalf("wrong format.Validate: expected!=nil, got==nil\n")
	}

	// date format must be a Go date format
	format.DateLayout = "gibberish"

	err = format.Validate()
	if err == nil {
		t.Fatalf("wrong format.Validate: expected!=nil, got==nil\n")
	}

	format = kbFull

	// memo field index cannot be zero
	format.MemoI = 0

	err = format.Validate()
	if err == nil {
		t.Fatalf("wrong format.Validate: expected!=nil, got==nil\n")
	}
}

func TestUnhappyConfigNFields(t *testing.T) {
	t.Parallel()

	format := kbFull

	// number of fields cannot be out of range, too low
	format.NFields = minNFields - 1

	err := format.Validate()
	if err == nil {
		t.Fatalf("wrong format.Validate: expected!=nil, got==nil\n")
	}

	// number of fields cannot be out of range, too high
	format.NFields = maxNFields + 1

	err = format.Validate()
	if err == nil {
		t.Fatalf("wrong format.Validate: expected!=nil, got==nil\n")
	}
}

func TestUnhappyConfigOptional(t *testing.T) {
	t.Parallel()

	format := kbFull

	// if amount field index is zero then both credit and debit indexes must be non-zero
	format.AmountI, format.CreditI, format.DebitI = 0, 1, 0

	err := format.Validate()
	if err == nil {
		t.Fatalf("wrong format.Validate: expected!=nil, got==nil\n")
	}
}

func TestUnhappyOutputFormat(t *testing.T) {
	t.Parallel()

	trn := Transaction{
		Date:        "2025-05-02",
		ThisAccount: "Assets:Current:KB05", OtherAccount: "Assets:Current:PCUS1",
		Memo:   "To term deposit",
		Amount: 42.00,
	}

	expect := ""

	got := trn.StringFormat("")
	if got != expect {
		t.Fatalf("wrong StringFormat: expected==%q, got==%q\n", expect, got)
	}

	got = trn.StringFormat("unknown")
	if got != expect {
		t.Fatalf("wrong StringFormat: expected==%q, got==%q\n", expect, got)
	}
}

func TestUnhappyTransact(t *testing.T) {
	t.Parallel()

	var trn Transaction

	// a zero values transaction is not valid
	err := trn.Validate()
	if err == nil {
		t.Fatalf("wrong trn.Validate: expected!=nil, got==nil")
	}

	if trn.IsValid() {
		t.Fatalf("wrong trn.IsValid: expected==false, got==true")
	}
}

func TestUnhappyTransactAccount(t *testing.T) {
	t.Parallel()

	// this and that cannot be the same account
	trn := Transaction{
		Date:        "2025-05-07",
		ThisAccount: "PCUS1", OtherAccount: "PCUS1",
		Memo:   "Oh dear!",
		Amount: 0.01,
	}

	err := trn.Validate()
	if err == nil {
		t.Fatalf("wrong trn.Validate: expected!=nil, got==nil")
	}
}

func TestUnhappyTransactAmount(t *testing.T) {
	t.Parallel()

	format := kbFull

	// amount cannot be zero
	flds := []string{"ZZ-YYYY-XXXXXXX-WW", "29-12-2023", "Automatic Payment Rates MISS E MACD ;Ref: Rates MISS E MACD",
		"AP", "Rates", "E", "", "", "", "", "MISS E MACD", "AA-BBBB-CCCCCCC-DD", "0.00", "", "0.00", "1434.23"}

	var trn Transaction

	err := trn.ParseCSV(flds, format)
	if err == nil {
		t.Fatalf("wrong trn.ParseCSV: expected!=nil, got==nil")
	}

	format = pcu

	// either credit or debit fields must have a value not both
	flds = []string{"07/01/2020", "554PHP 18832946 Best of Health", "16.92", "16.92", "265.01"}

	err = trn.ParseCSV(flds, format)
	if err == nil {
		t.Fatalf("wrong trn.ParseCSV: expected!=nil, got==nil")
	}

	// either credit or debit fields must have a value not neither
	flds = []string{"07/01/2020", "554PHP 18832946 Best of Health", "", "", "265.01"}

	err = trn.ParseCSV(flds, format)
	if err == nil {
		t.Fatalf("wrong trn.ParseCSV: expected!=nil, got==nil")
	}

	// debit field cannot be negative
	flds = []string{"07/01/2020", "554PHP 18832946 Best of Health", "-12.35", "", "265.01"}

	err = trn.ParseCSV(flds, format)
	if err == nil {
		t.Fatalf("wrong trn.ParseCSV: expected!=nil, got==nil")
	}

	// credit field cannot be negative
	flds = []string{"07/01/2020", "554PHP 18832946 Best of Health", "", "-12.35", "265.01"}

	err = trn.ParseCSV(flds, format)
	if err == nil {
		t.Fatalf("wrong trn.ParseCSV: expected!=nil, got==nil")
	}
}

func TestUnhappyTransactDate(t *testing.T) {
	t.Parallel()

	format := kbFull

	// date field cannot have different format from date format
	flds := []string{"ZZ-YYYY-XXXXXXX-WW", "29/12/2023", "Automatic Payment Rates MISS E MACD ;Ref: Rates MISS E MACD",
		"AP", "Rates", "E", "", "", "", "", "MISS E MACD", "AA-BBBB-CCCCCCC-DD", "162.00", "", "162.00", "1434.23"}

	var trn Transaction

	err := trn.ParseCSV(flds, format)
	if err == nil {
		t.Fatalf("wrong trn.ParseCSV: expected!=nil, got==nil")
	}

	// date format cannot be gibberish!
	format.DateLayout = "gibberish"

	err = trn.ParseCSV(flds, format)
	if err == nil {
		t.Fatalf("wrong trn.ParseCSV: expected!=nil, got==nil")
	}

	if trn.IsValid() {
		t.Fatalf("wrong trn.IsValid: expected==true, got==false")
	}
}

func TestUnhappyTransactMemo(t *testing.T) {
	t.Parallel()

	format := kbFull

	// memo cannot be empty string
	flds := []string{"ZZ-YYYY-XXXXXXX-WW", "29-12-2023", "",
		"AP", "Rates", "E", "", "", "", "", "MISS E MACD", "AA-BBBB-CCCCCCC-DD", "162.00", "", "162.00", "1434.23"}

	var trn Transaction

	err := trn.ParseCSV(flds, format)
	if err == nil {
		t.Fatalf("wrong trn.ParseCSV: expected!=nil, got==nil")
	}
}

func TestUnhappyTransactNFields(t *testing.T) {
	t.Parallel()

	format := pcu

	// number of fields cannot be different from NFields
	flds := []string{"28/11/2019", "HealthAndLif eInsuranceAn dSubs ARNHEMCR BP", "123.00", "316.69"}

	var trn Transaction

	err := trn.ParseCSV(flds, format)
	if err == nil {
		t.Fatalf("wrong trn.ParseCSV: expected!=nil, got==nil")
	}
}

func TestUnhappyTransactThisAcct(t *testing.T) {
	t.Parallel()

	format := kbFull

	// as this account field index is non-zero, this account field cannot be empty string
	flds := []string{"", "29-12-2023", "Automatic Payment Rates MISS E MACD ;Ref: Rates MISS E MACD",
		"AP", "Rates", "E", "", "", "", "", "MISS E MACD", "AA-BBBB-CCCCCCC-DD", "162.00", "", "", "1434.23"}

	var trn Transaction

	err := trn.ParseCSV(flds, format)
	if err == nil {
		t.Fatalf("wrong trn.ParseCSV: expected!=nil, got==nil")
	}
}

var kbFull = CSVFormat{ // for Kiwibank full CSV statement
	NFields: 16,
	AmountI: 15, CreditI: 13, DateI: 2, DebitI: 14,
	MemoI: 3, OtherAccountI: 12, ThisAccountI: 1,
	DateLayout: "02-01-2006",
}

var mini = CSVFormat{ // for minimal CSV statement
	NFields: 3,
	AmountI: 3, CreditI: 0, DateI: 1, DebitI: 0,
	MemoI: 2, OtherAccountI: 0, ThisAccountI: 0,
	DateLayout: "2006-01-02",
}

var pcu = CSVFormat{ // for PCU account CSV statement
	NFields: 5,
	AmountI: 0, CreditI: 4, DateI: 1, DebitI: 3,
	MemoI: 2, OtherAccountI: 0, ThisAccountI: 0,
	DateLayout: "02/01/2006",
}
