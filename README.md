# Arnhemcr/financial

This Go module translates financial transactions 
from an arbitrary [comma-separated values (CSV)] format to a standard format.
It allows account statements, 
with transactions in different CSV record formats,
to be merged for analysis and reporting.
This module contains a program and its supporting package.

## Program csv2trn

CSV2trn [filters] transactions from CSV records in an account statement to a standard format.
The format of the input CSV record is configured by an [XML] file,
which contains the record's number and position of fields and its date layout.
The output format for transactions is either
[Ledger] journal entries or this module's CSV records.

For more information and examples, see `go doc` in the csv2trn directory.

## Package transaction

This package represents a financial transaction as an instance of type Transaction.
A transaction is the transfer of an amount of currency from one account to another on a date.
It is described by a memo and code,
also known as the transaction's description and type respectively.
A transaction belong to an account known as this account.

This package offers:

  - parsing a transaction from a CSV record; 
    an instance of type CSVRecordFormat configures the parser for the record format
  - stringing a transaction to either a Ledger journal entry or this module's CSV record

For more information, see `go doc -all` in the transaction directory.

[comma-separated values (CSV)]: https://en.wikipedia.org/wiki/Comma-separated_values
[filters]: https://en.wikipedia.org/wiki/Filter_(software)
[Ledger]: https://en.wikipedia.org/wiki/Ledger_(software)
[XML]: https://en.wikipedia.org/wiki/XML
