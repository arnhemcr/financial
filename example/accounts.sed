#
# accounts.sed
#

# Substitute account number for Ledger account name.
s/,01-2345-6789012-34,/,Assets:Current,/
s/,12-3456-7890123-45,/,Assets:Emergency,/
s/,98-7654-3210987-65,/,Expenses:Rates,/
s/,01-0101-0101010-10,/,Income:Salary,/

# Match the memo and
# substitute default other account name Imbalance for Ledger account name.
/,Cash,/s/,Imbalance,/,Expenses:Cash,/
/,From current,/s/,Imbalance,/,Assets:Current,/
/,Net interest,/s/,Imbalance,/,Income:NetInterest,/
