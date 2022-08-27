# foldersize

get total size of folder on windows


## Usage

```
Usage:
  foldersize [OPTION] [PATTERN]

  Print total size of folders.
  Sort from largest to smallest size.

  OPTION:
    -m      ... display in MB
    -k      ... display in KB
```


## Example

```
$ foldersize
C:\Tools\CentOS7              3328390500
C:\Tools\LibreOfficePortable   712310154
C:\Tools\Go                    448664262
C:\Tools\vim81-kaoriya-win64    49216198

$ foldersize C:\Tools\G*
C:\Tools\Go  448664262

$ foldersize -k C:\Tools\G*
C:\Tools\Go     438149
```


## Installation

windows binary is [here](https://github.com/inazak/foldersize/releases)

