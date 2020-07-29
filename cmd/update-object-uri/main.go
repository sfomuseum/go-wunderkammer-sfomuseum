package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"github.com/aaronland/go-wunderkammer/oembed"
	"github.com/jtacoma/uritemplates"
	"github.com/tidwall/pretty"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const SFOMUSEUM_URI_TEMPLATE string = "sfom://id/{id}"

var sfomuseum_uri_template *uritemplates.UriTemplate

func init() {

	t, err := uritemplates.Parse(SFOMUSEUM_URI_TEMPLATE)

	if err != nil {
		panic(err)
	}

	sfomuseum_uri_template = t
}

func main() {

	format_json := flag.Bool("format", false, "Emit results as formatted JSON.")
	as_json := flag.Bool("json", false, "Emit results as a JSON array.")

	to_stdout := flag.Bool("stdout", true, "Emit to STDOUT")
	to_devnull := flag.Bool("null", false, "Emit to /dev/null")

	workers := flag.Int("workers", runtime.NumCPU(), "The number of concurrent workers to append data URLs with")
	timings := flag.Bool("timings", false, "Log timings (time to wait to process, time to complete processing")

	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	writers := make([]io.Writer, 0)

	if *to_stdout {
		writers = append(writers, os.Stdout)
	}

	if *to_devnull {
		writers = append(writers, ioutil.Discard)
	}

	if len(writers) == 0 {
		log.Fatal("Nothing to write to.")
	}

	wr := io.MultiWriter(writers...)

	reader := bufio.NewReader(os.Stdin)

	count := int32(0)

	throttle := make(chan bool, *workers)

	for i := 0; i < *workers; i++ {
		throttle <- true
	}

	mu := new(sync.RWMutex)
	wg := new(sync.WaitGroup)

	t0 := time.Now()

	for {

		select {
		case <-ctx.Done():
			break
		default:
			// pass
		}

		body, err := reader.ReadBytes('\n')

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Failed to read bytes, %v", err)
		}

		body = bytes.TrimSpace(body)

		var rec *oembed.Photo

		err = json.Unmarshal(body, &rec)

		if err != nil {
			log.Fatalf("Failed to unmarshal OEmbed record, %v", err)
		}

		t1 := time.Now()

		<-throttle

		if *timings {
			log.Printf("Time to wait to process %s, %v\n", rec.URL, time.Since(t1))
		}

		wg.Add(1)

		go func(rec *oembed.Photo) {

			t2 := time.Now()

			defer func() {

				throttle <- true
				wg.Done()

				if *timings {
					log.Printf("Time to complete processing for %s, %v\n", rec.URL, time.Since(t2))
				}
			}()

			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}

			u, err := url.Parse(rec.ObjectURI)

			if err != nil {
				log.Fatal(err)
			}

			path := u.Path
			id := filepath.Base(path)

			uri_values := make(map[string]interface{})
			uri_values["id"] = id

			object_uri, err := sfomuseum_uri_template.Expand(uri_values)

			if err != nil {
				log.Fatal(err)
			}

			rec.ObjectURI = object_uri

			body, err := json.Marshal(rec)

			if err != nil {
				log.Fatalf("Failed to marshal record, %v", err)
			}

			if *format_json {
				body = pretty.Pretty(body)
			}

			new_count := atomic.AddInt32(&count, 1)

			mu.Lock()
			defer mu.Unlock()

			if *as_json && new_count > 1 {
				wr.Write([]byte(","))
			}

			wr.Write(body)
			wr.Write([]byte("\n"))

		}(rec)

	}

	if *as_json {
		wr.Write([]byte("]"))
	}

	wg.Wait()

	if *timings {
		log.Printf("Time to process %d records, %v\n", count, time.Since(t0))
	}
}
