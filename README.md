# go-wunderkammer-sfomuseum

Tools for working with go-wunderkammer files in an SFO Museum context.

## Tools

To build binary versions of these tools run the `cli` Makefile target. For example:

```
$> go build -mod vendor -o bin/update-object-uri cmd/update-object-uri/main.go
```

### update-object-uri

Update to `object_uri` property to in the output of the `go-whosonfirst-data/bin/emit` tool to reflect SFO Museum.

```
$> /usr/local/go-whosonfirst-data/bin/emit \
	-query 'properties.wof:depicts=1159396315' \
	-oembed /usr/local/data/sfomuseum-data-media/data/ \
   
   | bin/update-object-uri \
   	-format \

   | /usr/local/go-wunderkammer-image/bin/append-dataurl \

   | /usr/local/go-wunderkammer/bin/wunderkammer-db \
	-database-dsn 'sql://sqlite3/usr/local/go-wunderkammer/gallery.db
```

## See also

* https://github.com/aaronland/go-wunderkammer
* https://github.com/aaronland/go-wunderkammer-image
* https://github.com/aaronland/ios-wunderkammer