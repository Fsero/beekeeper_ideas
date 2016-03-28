// This is an example of using elastic's BulkProcessor with Elasticsearch.
//
// See https://github.com/olivere/elastic and
// and https://github.com/olivere/elastic/wiki/BulkProcessor
// for more details.

/*
 * This example illustrates a simple process that performs bulk processing
 * with Elasticsearch using the BulkProcessor in elastic.
 *
 * It sets up a Bulker that runs a loop that will send data into
 * Elasticsearch. A second goroutine is started to periodically print
 * statistics about the process, e.g. the number of successful/failed
 * bulk commits.
 *
 * If you stop your cluster, the Bulker will find out as the After callback
 * gets an error passed. As soon as that happens, Bulker stops sending
 * new data into the BulkProcessor.
 *
 * When you start your cluster again, Bulker will also find out because it
 * has set up an automatic flush interval. This flush will eventually succeed
 * and the After callback will be called again, this time without an error.
 * The Bulker picks that up and starts sending data again.
 *
 * You can pass several parameters to the process:
 * The "url" parameter is the Elasticsearch HTTP endpoint (http://127.0.0.1:9200 by default).
 * The "index" parameter allows you to specify the Elasticsearch index name ("bulker-test" by default).
 * The "worker" parameter allows you to specify the number of parallel workers
 * the BulkProcessor gets started with.
 *
 */

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	elastic "gopkg.in/olivere/elastic.v3"
)

var b *Bulker
var once sync.Once

type BulkerDocument struct {
	document interface{}
	docType  string
}

// Bulker is an example for a process that needs to push data into
// Elasticsearch via BulkProcessor.
type Bulker struct {
	c            *elastic.Client
	p            *elastic.BulkProcessor
	workers      int
	index        string
	indexMapping string
	beforeCalls  int64         // # of calls into before callback
	afterCalls   int64         // # of calls into after callback
	failureCalls int64         // # of successful calls into after callback
	successCalls int64         // # of successful calls into after callback
	seq          int64         // sequential id
	stopC        chan struct{} // stop channel for the indexer
	dataC        chan BulkerDocument

	throttleMu sync.Mutex // guards the following block
	throttle   bool       // throttle (or stop) sending data into bulk processor?
}

type BulkerOptions struct {
	url          string
	workers      int
	index        string
	indexMapping string
}

func GetBulker(options BulkerOptions) *Bulker {
	once.Do(func() {
		client, err := elastic.NewClient(elastic.SetURL(*options.url))
		if err != nil {
			log.Fatal(err)
		}
		defer client.Stop()

		errc := make(chan error)

		// Run the bulker
		b = &Bulker{c: client, workers: *options.workers, index: *options.index, indexMapping: options.indexMapping}
		err = b.Run()
		if err != nil {
			log.Fatal(err)
		}
		defer b.Close()

		// Run the statistics printer
		go func(b *Bulker) {
			for range time.Tick(1 * time.Second) {
				printStats(b)
			}
		}(b)

		// Watch for SIGINT and SIGTERM from the console.
		go func() {
			c := make(chan os.Signal)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			fmt.Printf("caught signal %v\n", <-c)
			errc <- nil
		}()

		// Wait for problems.
		go func() {
			if err := <-errc; err != nil {
				log.Print(err)
				os.Exit(1)
			}
		}()
	})
	return b
}

func addDocument(b *Bulker, d *BulkerDocument) {
	b <- d
}

// printStats retrieves statistics from the Bulker and prints them.
func printStats(b *Bulker) {
	stats := b.Stats()
	var buf bytes.Buffer
	for i, w := range stats.Workers {
		if i > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%d=[%04d]", i, w.Queued))
	}

	fmt.Printf("%s | calls B=%04d,A=%04d,S=%04d,F=%04d | stats I=%05d,S=%05d,F=%05d | %v\n",
		time.Now().Format("15:04:05"),
		b.beforeCalls,
		b.afterCalls,
		b.successCalls,
		b.failureCalls,
		stats.Indexed,
		stats.Succeeded,
		stats.Failed,
		buf.String())
}

// Run starts the Bulker.
func (b *Bulker) Run() error {
	// Recreate Elasticsearch index
	if err := b.ensureIndex(); err != nil {
		return err
	}

	// Start bulk processor
	p, err := b.c.BulkProcessor().
		Workers(b.workers).              // # of workers
		BulkActions(1000).               // # of queued requests before committed
		BulkSize(4096).                  // # of bytes in requests before committed
		FlushInterval(10 * time.Second). // autocommit every 10 seconds
		Stats(true).                     // gather statistics
		After(b.after).                  // call "after" after every commit
		Do()
	if err != nil {
		return err
	}

	b.p = p

	// Start indexer that pushes data into bulk processor
	b.stopC = make(chan struct{})
	b.dataC = make(chan BulkerDocument)
	go b.indexer()

	return nil
}

// Close the bulker.
func (b *Bulker) Close() error {
	b.stopC <- struct{}{}
	<-b.stopC
	close(b.stopC)
	return nil
}

// indexer is a goroutine that periodically pushes data into
// bulk processor unless being "throttled" or "stopped".
func (b *Bulker) indexer() {
	var stop bool

	for !stop {
		select {
		case <-b.stopC:
			stop = true

		default:
			b.throttleMu.Lock()
			throttled := b.throttle
			b.throttleMu.Unlock()

			if !throttled {
				// Sample data structure
				doc := struct {
					Seq int64 `json:"seq"`
				}{
					Seq: atomic.AddInt64(&b.seq, 1),
				}
				// Add bulk request.
				// Notice that we need to set Index and Type here!
				r := elastic.NewBulkIndexRequest().Index(b.index).Type("doc").Doc(doc)
				b.p.Add(r)
			}
			// Sleep for a short time.
			time.Sleep(time.Duration(rand.Intn(7)) * time.Millisecond)
		}
	}

	b.stopC <- struct{}{} // ack stopping
}

// after is invoked by bulk processor after every commit.
// The err variable indicates success or failure.
func (b *Bulker) after(id int64, requests []elastic.BulkableRequest, response *elastic.BulkResponse, err error) {
	printStats(b)
}

// Stats returns statistics from bulk processor.
func (b *Bulker) Stats() elastic.BulkProcessorStats {
	return b.p.Stats()
}

// ensureIndex creates the index in Elasticsearch.
// It will be dropped if it already exists.
func (b *Bulker) ensureIndex() error {
	if b.index == "" {
		return errors.New("no index name")
	}
	exists, err := b.c.IndexExists(b.index).Do()
	if err != nil {
		return err
	}
	if !exists {
		_, err = b.c.CreateIndex(b.index).BodyString(b.indexMapping).Do()
		if err != nil {
			return err
		}
	}
	return nil
}
