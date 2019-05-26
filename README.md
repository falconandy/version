# Simple Versioning Library for Go

Parse and compare versions.

Version format:
```
VERSION = PART* POSTFIX
PART = NUMBER (POSTFIX NUMBER-POSTFIX?)?
``` 

For example:
```
1.4.8rc2-wheezy

PARTS: PART(NUMBER=1), PART(NUMBER=4), PART(NUMBER=8, POSTFIX=rc, NUMBER-POSTFIX=2)
POSTFIX: -wheezy
``