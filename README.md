# Farmer Friend Service

###Basic Dependencies

* install go
* add installed go path to global $PATH variable
* set $GOPATH (your go working directory)
* clone this repo at $GOPATH/src/github.com/

### Service Dependencies
```sh 
 $ setup/setup.sh
```

### Migrations
run the migrations from cloned folder 
```sh
$ goose -env=production -pgschema=farmerApp.db status
$ goose -env=production -pgschema=farmerApp.db up
$ goose -env=production -pgschema=farmerApp.db down
$ goose -env=production -pgschema=farmerApp.db create <migration name> -sql
```