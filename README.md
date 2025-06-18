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

More information and examples can be found in the csv2trn directory.
Once there and assuming [Go has been installed], build csv2trn by running `go build`.

Translating arbitrary CSV account statements is a challenge for financial software.
["The convert command" Ledger 3 Manual] shows the issues
with a statement from ValuFirst Credit Union.
As an example of csv2trn, that statement can be translated to Ledger journal entries
by running:
```
cat VFCU.csv | ./csv2trn -f VFCU.xml -t Assets:ValuFirst:Checking -c $
```
For help on csv2trn run `./csv2trn -h`,
and for the documentation, including further examples, run `go doc`.

## Program mergetrn

Mergetrn is another filter that merges transactions from Ledger journals:

  - discarding mirrored transactions that have the code "(MT)"
  - ordering the remaining transactions by date ascending

It assumes the date layout of the journal entries is "YYYY-MM-DD"
like the ones output by csv2trn.

Mirrored transactions are an issue when merging journals of accounts that have
transfers between those accounts.
Each transfer has a debit transaction in one journal mirrored by a credit transaction in another.
If one of the mirrored transactions is not discarded, the result is a double transfer.

An example can be found in the mergetrn directory.
Once there, build mergetrn by running `go build`.
For the documentation, including the example, run `go doc`.

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
["The convert command" Ledger 3 Manual]: https://ledger-cli.org/doc/ledger3.html#The-convert-command
[XML]: https://en.wikipedia.org/wiki/XML
