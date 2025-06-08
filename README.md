# Arnhemcr/financial

This Go module translates financial transactions.
It can read transactions in an arbitrary [comma-separated values (CSV)] format.
And it can write those transactions in various standard formats.
This module enables merging of transactions from CSV account statements in different formats.
Arnhemcr/financial contains a program and supporting package.

## Program translate

Translate [filters] transaction records from a CSV account statement to a standard format.

The input CSV record format is usually configured by an [XML] file.
It defaults to this module's CSV record.
If a statement line cannot be translated to a transaction, translate prints the error.

The output formats for transactions are [Ledger] journal entry, the default,
or this module's CSV record.
The output order is date ascending.

For more information and examples to try, run `go doc` in the translate directory.

## Package transaction

This package represents a financial transaction as an instance of type Transaction.
A transaction is the transfer of an amount of currency from one account to another.
The transfer takes place on a date.
It is described by a memo and code,
also known as the description and transaction type respectively.
A statement and its transactions belong to an account known as this account.

This package offers:

  - parsing a transaction from a CSV record
  - stringing a transaction to either a Ledger journal entry or this module's CSV record

The parser is configured by an instance of type CSVRecordFormat.

For more information, run `go doc -all` in the transaction directory.

[comma-separated values (CSV)]: https://en.wikipedia.org/wiki/Comma-separated_values
[filters]: https://en.wikipedia.org/wiki/Filter_(software)
[Ledger]: https://en.wikipedia.org/wiki/Ledger_(software)
[XML]: https://en.wikipedia.org/wiki/XML
