#
# adjust.sed
#

# Translate account numbers into names.
s/,01-0101-0101010-10,/,Income:Salary,/
s/,12-3456-7890123-45,/,Assets:Saving,/
s/,98-7654-3210987-65,/,Expenses:Rates,/

# Use memo as a clue to the name of the other account.
/,Holiday savings,/s/,Imbalance,/,Assets:Current,/
/,Net interest,/s/,Imbalance,/,Income:Interest,/

# Transfers between asset accounts e.g. the current and saving accounts
# appear as a transaction on both account statements: 
# a debit from one account and a mirroring credit to the other.
#
# To prevent a double transfer,
# assume the debit transaction exists and delete the mirroring credit transaction. 
/,Assets:.*,Assets:.*,.*,[0-9]/d
