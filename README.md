# PMSS - Poor man's security scaner

## Example dev commands

```
alias mage="go run mage.go"

# generate code
mage generate

# start a server in devcontainer
go run ./pmssd server

# start a client (scanner)
go run ./pmss scan ./pkg/scanner/testdata/
```

