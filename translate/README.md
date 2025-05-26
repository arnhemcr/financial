# Arnhemcr/transaction program translate example

Program translate is used in the following example to combine transactions
from two CSV account statements, in different formats,
into a file of Ledger journal entries.

```
translate -f national_bank.xml -t Assets:Current -o csv <national_bank.csv >all.csv
translate -f local_CU.xml -t Assets:Saving -o csv <local_CU.csv >>all.csv
sort -t , -k 1 all.csv | translate >all.ledger
```

