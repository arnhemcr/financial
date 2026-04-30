#
# accounts.sed
#

# Substitute account number for Ledger account name.
s/,01-2345-6789012-34,/,Assets:Current,/
s/,12-3456-7890123-45,/,Assets:Emergency,/
s/,98-7654-3210987-65,/,Expenses:Rates,/
s/,01-0101-0101010-10,/,Income:Salary,/

# Substitute default other account name Imbalance for Ledger account name
# depending on the transaction's memo or description.
/,Cash,/s/,Imbalance,/,Expenses:Cash,/
/,To emergency fund,/s/,Imbalance,/,Assets:Current,/
/,Net interest,/s/,Imbalance,/,Income:NetInterest,/
