# Arnhemcr/financial

This [Go] module offers [filter] programs to:
 * help translate financial transactions from a [comma-separated values (CSV)] statement
   into a journal for the [Ledger] command-line accounting system
 * merge Ledger journals into a general journal for reporting and analysis

Its only dependency is the [Go standard library].

# Example

An individual has two accounts: one with National Bank (NB) and the other with Local Credit Union (LCU).
These institutions provide CSV statements with different sets of transactions details and record formats.
This example translates those CSV statements into Ledger journals.
It then merges those journals into a general journal.

## Dependencies

* a [Go installation]
* a program to match text strings and substitute one for another e.g. the stream editor [sed]
* connecting programs together in a pipeline and redirecting output to a file
* a [Ledger installation]

## Install programs

Install this module's program csv2trn from its directory with `go install`.
Verify by viewing its help text with `csv2trn -h`.
Then install and verify programs mcsv2lent and mrglent.

## Translate CSV statements into Ledger journals

In the example directory, translate the statements into journals with:
```
# Initialise both journals.
cp NB_0.journal NB.journal
cp LCU_0.journal LCU.journal

cat NB.csv | csv2trn -f NB.xml -c GBP | sed -f accounts.sed | \
	mcsv2lent -f journalAccounts.xml >>NB.journal
cat LCU.csv | csv2trn -f LCU.xml -t Assets:Emergency -c GBP | sed -f accounts.sed | \
	mcsv2lent -f journalAccounts.xml >>LCU.journal
```

The pipelines load a CSV statement with cat, process its transactions then append entries to the journals.
A transaction is the transfer of an amount of currency between accounts on a particular day.
It is described by a memo and code, also called the description and transaction type.
A statement and its records belong to an account, which in a transaction is called this account.

Program csv2trn reads the statement line by line, parses transactions from CSV records following the input format in XML
and warns about lines that cannot be parsed.
Records from LCU do not provide this account, so it is set to its Ledger name Assets:Emergency.
If other account is not provided it defaults to Imbalance.
The program writes transactions in this module's CSV record format (mcsv) ordered by date ascending.

The stream editor sed substitutes Ledger account names for account numbers and for Imbalance by matching the transaction's memo.

This module has specific layouts for some details of a transaction:

* Amount: decimal integer or decimal expressions both with optional signs e.g. 1234, +1234 and -1234.56.
  This module does not support decimal separators other than '.', thousands separators or currencies in amounts.
* Date: YYYY-MM-DD or [ISO 8601] extended date. 
  Program csv2trn can be configured to read other layouts through the CSV input record format.

## Mark mirror entries in Ledger journals

Transfers between accounts with journals have two entries: a debit in one mirrored by a credit in the other.
When merging journals into a general journal, one of these entries must be discarded so the transfer happens once not twice.

Returning to the example above, mcsv2lent reads transactions in mcsv format then writes them in Ledger entry format (lent).
For transfers between accounts with journals, whose Ledger account names are on the list in XML, the credit entry is marked with "mirror entry" comments.

Use Ledger to verify the LCU journal with `ledger -f LCU.journal register Assets:Emergency`.
There are three entries with a current balance of 42.42 GBP.
Verify the "To emergency fund" entry in that journal is marked as a mirror with `cat LCU.journal`.

## Merge Ledger journals into a general journal

Merge the journals into a general journal with:
```
cat NB.journal LCU.journal | mrglent >general.journal
```
Program mrglent reads the journals and writes entries ordered by date ascending.
All other journal content is discarded including mirror entries, automatic transactions and command directives as well as block and global comments.

Verify the general journal with `ledger -f general.journal register Assets:Emergency` which has the same entries and balance as above.
Then verify the accounts and their balances with `ledger -f general.journal balance` are:
```
 96.28 GBP  Assets
 53.86 GBP    Current
 42.42 GBP    Emergency
-37.79 GBP  Equity:OpeningBalances
 32.63 GBP  Expenses
 20.00 GBP    Cash
 12.63 GBP    Rates
  4.38 GBP  Imbalance
-95.50 GBP  Income
 -0.13 GBP    NetInterest
-95.37 GBP    Salary
----------
         0
```

[comma-separated values (CSV)]: https://en.wikipedia.org/wiki/Comma-separated_values
[filter]: https://en.wikipedia.org/wiki/Filter_(software)
[Go]: https://go.dev
[Go installation]: https://go.dev/doc/install
[Go standard library]: https://pkg.go.dev/std
[ISO 8601]: https://en.wikipedia.org/wiki/ISO_8601#Calendar_dates
[Ledger]: https://ledger-cli.org
[Ledger installation]: https://ledger-cli.org/download.html
[sed]: https://www.gnu.org/software/sed/manual/sed.html
