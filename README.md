# Arnhemcr/financial

Arnhemcr/financial is a native Go module containing:

* package transaction which represents financial transactions, and offers methods to:
  * parse a transaction from a [comma-separated values] (CSV) record in an arbitrary format
  * string a transaction to the package's own CSV record format or [Ledger] format
* program translate: a [filter] which translates financial transactions from
  the package's own CSV format to Ledger format.
  It can be updated to translate CSV account statements in other formats.

[comma-separated values]: https://www.ietf.org/rfc/rfc4180.txt
[filter]: https://en.wikipedia.org/wiki/Filter_(software)
[Ledger]: https://ledger-cli.org "Ledger command-line accounting"
