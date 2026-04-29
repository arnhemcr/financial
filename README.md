# Arnhemcr/financial

This [Go] module offers programs to:

 * translate a financial transaction statement
   from an arbitrary [comma-separated values (CSV)] format to a [Ledger] journal
 * merge multiple Ledger journals into one general journal for reporting and analysis

According to the Ledger 3 manual:

> "Importing csv files is a lot of work, ..."
>
> &mdash; [The `convert` command]

This module aims to make it a little easier.
Its sole dependency is the [Go standard library].

## One CSV statement to one Ledger journal

Program csv2trn reads a statement, extracts a transaction from each CSV record
and writes the transactions in a standard format/
The format of the input CSV record is configured in an [XML] file.
The output format is either Ledger entry (the default) or this module's CSV record.

Assuming [Go has been installed], build csv2trn in its directory with `go build`.

As an example of a CSV statement, the Ledger manual gives one from ValuFirst Credit Union
(see ["The convert command" in the Ledger 3 manual]).
As an example of csv2trn, that statement can be translated to a Ledger journal with:
```
cat VFCU.csv | ./csv2trn -t Assets:ValuFirst:Checking -f VFCU.xml >VFCU.journal
```
Command flags name this account, the Ledger account this statement belongs to,
and the XML configuration file.
The statement contains lines which are not transaction records, so csv2trn warns about them.
Print the journal with `cat VFCU.journal` :
```
2011-12-12 Tuscan IT #00037657
 Assets:ValuFirst:Checking  -29.73
 Imbalance
2011-12-13 ID: 1741472662 CO: XXAA.COM PAYMNT
 Assets:ValuFirst:Checking  -236.65
 Imbalance
...
2011-12-13 CASH DEPOSIT
 Assets:ValuFirst:Checking  45
 Imbalance
```
Note the CSV records do not give the other account in each transaction,
so other account defaults to Imbalance.
<!--
Adding `-o mcsv` to the command above translates that statement to this module's CSV records:
```
...
2011-12-12,Assets:ValuFirst:Checking,Imbalance,,Tuscan IT #00037657,-29.73,$
2011-12-13,Assets:ValuFirst:Checking,Imbalance,,ID: 1741472662 CO: XXAA.COM PAYMNT,-236.65,$
...
2011-12-13,Assets:ValuFirst:Checking,Imbalance,,CASH DEPOSIT,45,$
```
-->

Get help on csv2trn with `./csv2trn -h`
and get documentation, including further examples, with `go doc`.

## Multiple statements each to its own journal then merged to one general journal

<!--
## Program mrglent

Mrglent merges multiple Ledger financial journals containing entries
(containing transactions) into one journal.
It is a filter that:
  - removes entries marked as mirrors to avoid double transfers between accounts with journals
  - sorts the remaining entries by date ascending

Mrglent only merges entries: transactions with dates.
Automatic transactions, comments and command directives in the input journals
are not currently copied to the output journal.
It uses the same "YYYY-MM-DD" date layout, for input and output,
which other arnhemcr/financial programs use for output.

Install mrglent with `go install`.
Get documentation, including an example, with `go doc`.
-->

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

Get more information with `go doc -all` in the transaction directory.

[comma-separated values (CSV)]: https://en.wikipedia.org/wiki/Comma-separated_values
[filters]: https://en.wikipedia.org/wiki/Filter_(software)
[Go]: https://go.dev
[Go has been installed]: https://go.dev/doc/install
[Go standard library]: https://pkg.go.dev/std
[Ledger]: https://ledger-cli.org
["The convert command" in the Ledger 3 manual]: https://ledger-cli.org/doc/ledger3.html#The-convert-command
[XML]: https://en.wikipedia.org/wiki/XML
