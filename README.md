# APM
Arrietty Package Manager  
PLEASE Set ENV ARRIETTY_PM_PATH.  
The recommended is the root of this repository.

## Commands
- get  
If @TAG was not there, the latest is fetched.
It's a copy of `go get`.
```shell
apm get $REPO@$TAG
# apm get github.com/xxx/yyy@v0.0.1 for github release tag "v0.0.1"
# apm get github.com/xxx/yyy for latest github release
```