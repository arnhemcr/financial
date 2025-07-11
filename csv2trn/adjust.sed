#
# adjust.sed
#

# Replace account numbers with Ledger account names.
s/,01-2345-6789012-34,/,Assets:Current,/
s/,01-2345-6789012-35,/,Assets:Emergency,/
s/,12-3456-7890123-45,/,Assets:Saving,/
s/,98-7654-3210987-65,/,Expenses:Rates,/
s/,01-0101-0101010-10,/,Income:Salary,/

# Infer other account name from memo.
/,Holiday savings,/s/,Imbalance,/,Assets:Current,/
/,Net interest,/s/,Imbalance,/,Income:NetInterest,/
