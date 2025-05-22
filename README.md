# Arnhemcr/financial

Arnhemcr/financial translates financial transactions
from an arbitrary [comma-separated values (CSV)] format 
into a standard format.
It allows transactions in different formats,
from different financial institutions or account types,
to be combined.

Arnhemcr/financial is a Go module, which contains a package and a program.

## Package transaction

Package transaction represents financial transactions.
It offers:

* a configurable parser for CSV records into transactions
* stringers for transactions into the package's default CSV records or [Ledger] journal entries

## Program translate

Program translate is a [filter] which:

* reads an account statement line by line
* parses a transaction from the CSV record on each line <!--Ed: mention config-->
* strings each transaction into a default CSV record or Ledger journal entry
* writes the transaction strings <!--ordered by date ascending-->

[comma-separated values (CSV)]: https://www.ietf.org/rfc/rfc4180.txt
[filter]: https://en.wikipedia.org/wiki/Filter_(software)
[Ledger]: https://ledger-cli.org "Ledger command-line accounting"
