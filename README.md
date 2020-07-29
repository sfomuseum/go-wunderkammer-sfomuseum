# go-wunderkammer-sfomuseum

## Tools

### update-object-uri

Update to `object_uri` property to in the output of the `go-whosonfirst-data/bin/emit` tool to reflect SFO Museum.

```
$> /usr/local/sfomuseum/go-whosonfirst-data/bin/emit \
	-query 'properties.wof:depicts=1159396315' \
	-oembed /usr/local/data/sfomuseum-data-media/data/ \

   | bin/update-object-uri/main.go \
   	-format
```

## See also

* https://github.com/aaronland/ios-wunderkammer