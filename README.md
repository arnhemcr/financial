# Arnhemcr/financial

This [Go] module offers programs to:

 * translate a [comma-separated values (CSV)] financial transaction statement
   in an arbitrary format to a [Ledger] journal
 * merge multiple Ledger journals into one general journal for reporting and analysis

According to the Ledger 3 manual:

> "Importing csv files is a lot of work, ..."
>
> &mdash; [The convert command]

This module aims to make it a little easier.
Its sole dependency is the [Go standard library].

## Translate CSV transaction records to Ledger journal entries

Program csv2trn reads a statement, extracts a transaction from each CSV record
and writes the transactions in a standard format.
The format of the input CSV record is configured in an [XML] file.
The output format is either Ledger entry (the default "lent")
or this module's CSV record ("mcsv").

Assuming [Go has been installed], build csv2trn in its directory with `go build` .

As an example of a CSV statement,
the Ledger manual gives one from ValuFirst Credit Union (VFCU)
(see ["The convert command" in the Ledger 3 manual]).
As an example of csv2trn, that statement can be translated to a Ledger journal with:
```
cat VFCU.csv | ./csv2trn -t Assets:ValuFirst:Checking -f VFCU.xml
```
CSV2trn warns about lines which are not transaction records.
Entries in the journal include:
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
Ledger entries require names for both this and other account.
However, VFCU transaction records do not contain these details.
So this account, the Ledger account that the statement belongs to, is set with a command flag,
while other account defaults to "Imbalance".

Install csv2trn with `go install` then verify by getting the help text with `csv2trn -h` .

## Set Ledger account names

CSV transaction records contain a memo or description, and some contain account numbers.
The stream editor sed can be configured to set Ledger account names from these details.

In the sed directory,
the following example adds entries from an individual's
current account at National Bank (NB) and emergency fund at Local Credit Union (LCU)
to their Ledger journals for those accounts:
```
# Initialise the journals with opening balances.
cp NB_0.journal NB.journal
cp LCU_0.journal LCU.journal

# Add entries translated from the statements to the journals.
cat NB.csv | csv2trn -o mcsv -f NB.xml | sed -f accounts.sed | csv2trn >>NB.journal
cat LCU.csv | csv2trn -o mcsv -t Assets:Emergency -f LCU.xml | sed -f accounts.sed | \
	csv2trn >>LCU.journal
```
View the journals with `cat NB.journal` and `cat LCU.journal` .
In the example, sed modifies transactions in this module's CSV record format ("mcsv"),
which is also the default input for csv2trn.

## Mark mirror Ledger journal entries

In the example above,
Ledger accounts Assets:Current and Assets:Emergency each have a journal.
A transfer between those accounts has an entry in both journals:
a debit entry from one mirroring a credit to the other
(see the entry with memo "To emergency fund" in the journals).
To merge these journals into a general journal,
one entry for each mirrored transfer must be removed
so the transfer happens once not twice.

Program mcsv2lent marks the credit entry of transfers between Ledger accounts with journals.
The names of Ledger accounts with journals are read from an XML file.
Marked entries are removed during merging.

Build, install and verify mcsv2lent from its directory. Then in the sed directory,
repeat the last example replacing the second csv2trn with mcsv2lent:
```
cp NB_0.journal NB.journal
cp LCU_0.journal LCU.journal

cat NB.csv | csv2trn -o mcsv -f NB.xml | sed -f accounts.sed | \
	mcsv2lent -f ../mcsv2lent/journalAccounts.xml >>NB.journal
cat LCU.csv | csv2trn -o mcsv -t Assets:Emergency -f LCU.xml | sed -f accounts.sed | \
	mcsv2lent -f ../mcsv2lent/journalAccounts.xml >>LCU.journal
```
The credit "To emergency fund" entry in LCU.journal is marked with mirror entry comments,
while the debit entry in NB.journal is as before.

## Merge Ledger journals to a general journal

Build, install and verify mrglent from its directory.

```
cat ../sed/NB.journal ../sed/LCU.journal >general.journal | mrglent >general.journal
```

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
