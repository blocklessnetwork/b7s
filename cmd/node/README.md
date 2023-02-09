# Node

## Usage

## Readme

# TODOs

- Consider using two separate databases for peer store and function storing.
Using the same database for both might mean that one will suffer performance penalty because of the other.
Even if they are kept the same - perhaps the keys could have different prefixes for the two record types - to have clearer separation.

- Add validation tags to config structures.