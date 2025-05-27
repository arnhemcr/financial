# Arnhemcr/financial

Arnhemcr/financial translates financial transactions
from an arbitrary [comma-separated values (CSV)] format into a standard format.
It allows transactions in different formats to be merged for analysis.
Arnhemcr/financial is a Go module, which contains a package and a program.

## Package transaction

This package represents financial transactions as Transaction structures.
It offers:

  - parsing a transaction from a CSV record
  - stringing a transaction to either this module's CSV record or a [Ledger] journal entry

The parser is configured for the format of the record by the CSVFormat structure.

## Program translate

This program is a [filter] that translates a CSV statement on standard input
to a list of transactions on standard output.

The input format, a CSVFormat structure, is usually read from an XML file (-f flag),
and its default is this module's CSV record.
The output format is either this module's CSV records or the default Ledger journal entries.

A transaction is the transfer of an amount of currency from one account to another 
with a memo (or description) on a date.
This account is the name or number of the account that this statement and
its transactions belong to.

This program is explained by example using files in the translate/ directory.

### Translate CSV account statement into Ledger journal entries

```
translate -f national_bank.xml <national_bank.csv
```

In this example, the National Bank CSV statement contains all the fields required for a transaction.
It also contains a header line, which translate warns is not parseable.

### Translate CSV account statement into this module's CSV records

```
translate -f local_CU.xml -t Assets:Saving -o csv <local_CU.csv
```

This account is required for a transaction.
In this example, the Local Credit Union statement does not contain this account,
so it is set to its Ledger name from the translate command (-t flag).
The output format is also set to this module's CSV record (-o flag).

### Merge CSV account statements into Ledger journal entries

```
translate -f national_bank.xml -t Assets:Current -o csv <national_bank.csv >all.csv
translate -f local_CU.xml -t Assets:Saving -o csv <local_CU.csv >>all.csv
sed -f adjust.sed all.csv | sort -t , -k 1 | translate
```

In this example, the CSV statements are translated into this module's CSV records.
The stream editor (sed) replaces account numbers with names,
guesses the name of the other account from the memo,
and, for transfers between the accounts with statements, deletes mirrored credit transactions.
The CSV records are sorted by date ascending then translated into Ledger journal entries.

[comma-separated values (CSV)]: https://www.ietf.org/rfc/rfc4180.txt
[filter]: https://en.wikipedia.org/wiki/Filter_(software)
[Ledger]: https://ledger-cli.org "Ledger command-line accounting"
