# Arnhemcr/financial

This Go module translates financial transactions
from [comma-separated values (CSV)] to a standard format.
It allows transactions from CSV account statements in different formats
to be merged for analysis.
Arnhemcr/financial contains a program and supporting package.

## Program translate

Translate [filters] transaction records from a CSV account statement into a standard format.

The input CSV record format is usually configured by an [XML] file.
It defaults to this module's CSV record.
If a statement line cannot be translated into a transaction, it is reported as an error.

The output transaction format defaults to a [Ledger] journal entry.
It can also be set to this module's CSV record.
Transactions are ordered by date ascending.

For more information and examples to try, run `go doc` in the translate directory.

## Package transaction

This package represents a financial transaction as an instance of type Transaction.
A transaction is the transfer of an amount of currency from one account to another.
The transfer takes place on a date and described by a memo.
A statement and its transactions belong to an account, which is known as this account.

This package offers:

  - parsing a transaction from a CSV record
  - stringing a transaction to either a Ledger journal entry or this module's CSV record

The parser is configured by an instance of type CSVRecordFormat.

For more information, run `go doc -all` in the transaction directory.

[comma-separated values (CSV)]: https://en.wikipedia.org/wiki/Comma-separated_values
[filters]: https://en.wikipedia.org/wiki/Filter_(software)
[Ledger]: https://en.wikipedia.org/wiki/Ledger_(software)
[XML]: https://en.wikipedia.org/wiki/XML
