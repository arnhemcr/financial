# Arnhemcr/financial

Arnhemcr/financial is a native Go module containing:

* package transaction which represents financial transactions, and offers methods to:
  * parse a transaction from a [comma-separated values] (CSV) record in an arbitrary format
  * string a transaction to the package's own CSV record format or [Ledger] format

[comma-separated values]: https://www.ietf.org/rfc/rfc4180.txt
	"Common Format and MIME Type for Comma-Separated Values (CSV) Files"
[Ledger]: https://ledger-cli.org "Ledger command-line accounting"
