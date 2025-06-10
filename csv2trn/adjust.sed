#
# adjust.sed
#

# Replace account numbers with names.
s/,01-2345-6789012-34,/,Assets:Current,/
s/,12-3456-7890123-45,/,Assets:Saving,/
s/,98-7654-3210987-65,/,Expenses:Rates,/
s/,01-0101-0101010-10,/,Income:Salary,/

# Infer other account name from memo.
/,Holiday savings,/s/,Imbalance,/,Assets:Current,/
/,Net interest,/s/,Imbalance,/,Income:Interest,/

# Transfers between asset and liability accounts
# e.g. the current and saving asset accounts
# appear as a transaction on both statements:
# a debit from one account and a mirroring credit to the other.
# To prevent a double transaction,
# assume the debit exists and comment out the mirroring credit.
/,(Assets|Liabilities):.*,(Assets|Liabilities):.*,.*,[0-9]/d
