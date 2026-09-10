# Arnhemcr/financial

This [Go] module offers three [filter] programs to:
* help translate financial transactions from [comma-separated values (CSV)] records in an account statement
  to journal entries for the [Ledger] command-line accounting system
* merge multiple Ledger journals into one general journal for reporting and analysis

A financial transaction is the transfer of an amount of currency from one account to another on a date.
It is described by a memo and code, also called the description and transaction type.
A statement lists transactions on an account made over a period of time.
Each transaction has two accounts: this account, the one which the transaction and its statement belong to, and the other account.

This module supports the following layouts for transaction details:
* Amount: decimal number with optional sign e.g. "1234.56", "-98.765" "+1234".
  Decimal separators other than '.', thousands separators and amounts including currencies are not supported.
* Date: YYYY-MM-DD also known as Go time.DateOnly and [ISO 8601 extended date].
  Program `csv2trn` supports input CSV records with other date layouts.

This module's only dependency is the standard library in the [Go installation].
The examples depend on a [Ledger installation], [pipelines] and [redirection] of output to a file.
They also use the stream editor `sed` but other programs that can match text strings and substitute one string for another could be used instead.

## Program `csv2trn`

This program translates financial transactions from CSV records in an account statement to formats including Ledger journal entries.
The output is ordered by date ascending.

In the `csv2trn` directory, build and install the program with `go install`. 
Then verify the program by getting its help text with `csv2trn -?`.

In the `example` directory, translate a National Bank statement with:
```
cat NB.csv | csv2trn -f NB.xml -o lent
```
The input CSV record format is configured in the XML file, while the output format is set to Ledger journal entries (`lent`).

Now translate a Local Credit Union statement:
```
cat LCU.csv | csv2trn -f LCU.xml -o lent -t Assets:Emergency 
```
In contrast to the bank's CSV records, those from the credit union do not have this or other account, and they are in reverse order.
The program sets this account to `Assets:Emergency` and other account defaults to `Imbalance`, and outputs the entries ordered by date ascending.

This module's two remaining programs are used when merging multiple Ledger journals into one general journal.
Those who have just one account and one Ledger journal can stop reading here.

## Accounts, journals and mirror entries

The bank and credit union accounts above belong to the same person.
To get a complete picture of their finances,
the journals for those accounts are merged into one general journal.

Transactions between accounts with journals have two entries: a debit in one journal mirrored by a credit in the other.
One of those entries must be discarded during merging so the transaction appears in the general journal once not twice.

## Program `mcsv2lent`

This program translates financial transactions from this module's CSV records (`mcsv`) to Ledger journal entries (`lent`).
It also encloses credit mirror entries with comments.

In the `mcsv2lent` directory, build, install and verify the program.

In the `example` directory, create Ledger journals from the bank and credit union statements with:
```
# Initialise the journals with opening balances.
cp NB_0.journal NB.journal
cp LCU_0.journal LCU.journal

# Add transactions from the bank and credit union statements to the journals.
cat NB.csv | csv2trn -f NB.xml -t Assets:Current -c GBP | sed -f accounts.sed | \
	mcsv2lent -f journalAccounts.xml >>NB.journal
cat LCU.csv | csv2trn -f LCU.xml -t Assets:Emergency -c GBP | sed -f accounts.sed | \
	mcsv2lent -f journalAccounts.xml >>LCU.journal
```
This time, `csv2trn` outputs transactions in its default format: this module's CSV records (`mcsv`).
The stream editor `sed` substitutes Ledger account names for account numbers and for `Imbalance` by matching the memo.
Then `mcsv2lent` translates this module's CSV records to Ledger journal entries (`lent`).
It also encloses the credit entry of each transaction between accounts with journals using mirror entry comments (see the mirror entry in `LCU.journal`).
The names of journalled accounts are listed in the XML file.

Check the balance on each account according to its journal with:
```
ledger -f NB.journal register Current
ledger -f LCU.journal register Emergency
```
The bank and credit union account balances should be 53.86 and 42.42 GBP respectively.

## Program `mrglent`

This program merges financial transactions in Ledger entry (`lent`) format from multiple journals.
It also discards entries enclosed with mirror comments.
The output is ordered by date ascending.

In the `mrglent` directory, build, install and verify the program.

In the `example` directory, merge the bank and credit union journals into a general journal with:
```
cat NB.journal LCU.journal | mrglent >general.journal
```

Check the balances on both accounts according to the general journal with:
```
ledger -f general.journal balance
```
The bank and credit account balances should again be 53.86 and 42.42 GBP respectively.

[comma-separated values (CSV)]: https://en.wikipedia.org/wiki/Comma-separated_values
[filter]: https://en.wikipedia.org/wiki/Filter_(software)
[Go]: https://go.dev
[Go installation]: https://go.dev/doc/install
[ISO 8601 extended date]: https://en.wikipedia.org/wiki/ISO_8601#Calendar_dates
[Ledger]: https://en.wikipedia.org/wiki/Ledger_(software)
[Ledger installation]: https://ledger-cli.org/download.html
[pipelines]: https://en.wikipedia.org/wiki/Pipeline_(software)
[redirection]: https://en.wikipedia.org/wiki/Redirection_(computing)
