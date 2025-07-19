# Arnhemcr/financial

This Go module translates financial transactions 
from arbitrary [comma-separated values (CSV)] to a standard format.
It allows transactions from account statements in a variety of CSV formats
to be merged for analysis and reporting.

## Program csv2trn

CSV2trn [filters] transactions from CSV records, in an account statement, to a standard format.
The format of the input CSV records is configured by [XML].
The format of the output transactions can be
either [Ledger] journal entries or this module's CSV records.

Assuming [Go has been installed], build csv2trn in its directory by running `go build`.

Translating arbitrary CSV account statements is a challenge for financial software.
["The convert command" in the Ledger 3 manual] shows the issues
with a statement from ValuFirst Credit Union.
As an example of csv2trn, that statement can be translated to Ledger journal entries
by running:
```
cat VFCU.csv | ./csv2trn -f VFCU.xml -t Assets:ValuFirst:Checking -c $
```
For help on csv2trn run `./csv2trn -h`,
and for documentation, including further examples, run `go doc`.

## Program mrglent

Mrglent merges multiple Ledger financial journals containing entries,
also known as transactions, into one journal. It is a filter that:
  - removes mirror entries, that have been marked with the transaction code "MT",
    to avoid double transfers between accounts with journals
  - sorts the remaining entries by date ascending

Mrglent only merges entries: transactions that have dates.
Automatic transactions, comments and command directives in the input journals
are not currently copied to the output journal.
It uses the same "YYYY-MM-DD" date layout, for input and output,
that other arnhemcr/financial programs use for output.

Install mrglent by running `go install`.
For documentation, including an example, run `go doc`.

## Package transaction

This package represents a financial transaction as an instance of type Transaction.
A transaction is the transfer of an amount of currency from one account to another on a date.
It is described by a memo and code,
also known as the transaction's description and type respectively.
A transaction belongs to an account known as this account.

This package offers:

 - parsing a transaction from a CSV record;
   an instance of type CSVRecordFormat configures the parser for the record format
 - stringing a transaction to either a Ledger journal entry or this module's CSV record

For more information, see `go doc -all` in the transaction directory.

[comma-separated values (CSV)]: https://en.wikipedia.org/wiki/Comma-separated_values
[filters]: https://en.wikipedia.org/wiki/Filter_(software)
[Go has been installed]: https://go.dev/doc/install
[Ledger]: https://en.wikipedia.org/wiki/Ledger_(software)
["The convert command" in the Ledger 3 manual]: https://ledger-cli.org/doc/ledger3.html#The-convert-command
[XML]: https://en.wikipedia.org/wiki/XML
