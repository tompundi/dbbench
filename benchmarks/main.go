package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/tompundi/dbbench/testdb"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	duration = flag.Duration("d", time.Minute, "test duration for each case")
	c        = flag.Int("c", runtime.NumCPU(), "concurrent goroutines")
	keysize  = flag.Int64("keysize", 8, "key size")
	valsize  = flag.Int64("valsize", 256, "value size")
	txnum    = flag.Int64("txnum", 1000000, "tx number")
	fix      = flag.Bool("fix", true, "fix")
	fsync    = flag.Bool("fsync", false, "fsync")
	s        = flag.String("s", "badger", "store type")
	data     = make([]byte, *valsize)
)

func main() {

	flag.Parse()
	fmt.Printf("duration=%v, c=%d size=%d\n", *duration, *c, *valsize)

	var memory bool
	var path string

	store, path, err := getStore(*s, *fsync, path)
	if err != nil {
		panic(err)
	}
	if !memory {
		defer os.RemoveAll(path)
	}

	defer store.Close()
	name := *s
	if *fsync {
		name = name + "/fsync"
	} else {
		name = name + "/nofsync"
	}

	testBatchWrite(name, store)
	testSet(name, store)
	testGet(name, store)
	testGetSet(name, store)
	testDelete(name, store)
}

func testBatchWrite(name string, store testdb.Store) {
	var wg sync.WaitGroup
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), *duration)
	defer cancel()

	var totalEntries uint64
	var totalSize uint64
	for j := 0; j < *c; j++ {
		wg.Add(1)
		go func(proc int) {
			batchSize := uint64(1000)
			var keyList, valList [][]byte

			for i := uint64(0); i < batchSize; i++ {
				keyList = append(keyList, genKey(i))
				valList = append(valList, make([]byte, *valsize))
			}
			addKeySize := (uint64(unsafe.Sizeof(keyList[0]))+uint64(unsafe.Sizeof(keyList[0][0]))*uint64(*keysize))*batchSize + uint64(unsafe.Sizeof(keyList))
			addValSize := (uint64(unsafe.Sizeof(valList[0]))+uint64(unsafe.Sizeof(valList[0][0]))*uint64(*valsize))*batchSize + uint64(unsafe.Sizeof(valList))
			addSize := addKeySize + addValSize

		LOOP:
			for {
				select {
				case <-ctx.Done():
					break LOOP
				default:
					if *fix == true && totalEntries >= uint64(*txnum) {
						break LOOP
					}
					for i := range keyList {
						rand.Read(keyList[i])
						rand.Read(valList[i])
					}
					err := store.PSet(keyList, valList)
					if err != nil {
						panic(err)
					}
					atomic.AddUint64(&totalEntries, uint64(len(keyList)))
					atomic.AddUint64(&totalSize, addSize)
				}
			}
			wg.Done()
		}(j)
	}
	wg.Wait()
	diffTimes := float64(totalEntries) / float64(*txnum)
	if diffTimes == 0 {
		diffTimes = 1
	}
	seconds := float64(time.Since(start).Seconds())
	sizeMB := float64(totalSize) / 1e6
	fmt.Printf("%s batch-write size: %.2f MB; for total %d entries, took: %.2f s \n", name, sizeMB, totalEntries, seconds)
	fmt.Printf("%s batch-write size: %.2f MB; for total %d entries, took: %.2f s \n", name, sizeMB/diffTimes, *txnum, seconds/diffTimes)
}

func testSet(name string, store testdb.Store) {
	var wg sync.WaitGroup
	wg.Add(*c)

	ctx, cancel := context.WithTimeout(context.Background(), *duration)
	defer cancel()

	var totalCount int64
	start := time.Now()
	for j := 0; j < *c; j++ {
		go func() {
			i := uint64(j)
		LOOP:
			for {
				select {
				case <-ctx.Done():
					break LOOP
				default:
					if *fix == true && totalCount >= *txnum {
						break LOOP
					}
					store.Set(genKey(i), data)
					i += uint64(*c)
					totalCount++
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	dur := time.Since(start)
	ops := float64(totalCount) / (float64(dur) / 1e9)
	ns := float64(dur) / float64(totalCount*int64(*c))
	fmt.Printf("%s set rate: %.2f op/s, mean: %.2f ns, took: %d s\n", name, ops, ns, int(dur.Seconds()))
	fmt.Printf("%s set rate: %.2f op/s, for: %d tx, took: %.2f s\n", name, ops, *txnum, float64(*txnum)/ops)
}

func testGet(name string, store testdb.Store) {
	var wg sync.WaitGroup
	wg.Add(*c)

	ctx, cancel := context.WithTimeout(context.Background(), *duration)
	defer cancel()
	var totalCount int64
	start := time.Now()
	for j := 0; j < *c; j++ {
		index := uint64(j)
		go func() {
			i := index
		LOOP:
			for {
				select {
				case <-ctx.Done():
					break LOOP
				default:
					if *fix == true && totalCount >= *txnum {
						break LOOP
					}
					_, ok, _ := store.Get(genKey(i))
					if !ok {
						i = index
					}
					i += uint64(*c)
					totalCount++
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	dur := time.Since(start)
	ops := float64(totalCount) / (float64(dur) / 1e9)
	ns := float64(dur) / float64(totalCount*int64(*c))
	fmt.Printf("%s get rate: %.2f op/s, mean: %.2f ns, took: %d s\n", name, ops, ns, int(dur.Seconds()))
	fmt.Printf("%s get rate: %.2f op/s, for: %d tx, took: %.2f s\n", name, ops, *txnum, float64(*txnum)/ops)
}

func testGetSet(name string, store testdb.Store) {
	var wg sync.WaitGroup
	wg.Add(*c)

	ch := make(chan struct{})
	var totalCount int64
	setCount := uint64(0)
	go func() {
		for {
			select {
			case <-ch:
				return
			default:
				store.Set(genKey(setCount), data)
				setCount++
			}
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), *duration)
	defer cancel()

	totalCount += int64(setCount)
	start := time.Now()
	for j := 0; j < *c; j++ {
		index := uint64(j)
		go func() {
			i := index
		LOOP:
			for {
				select {
				case <-ctx.Done():
					break LOOP
				default:
					if *fix == true && totalCount >= *txnum {
						break LOOP
					}
					_, ok, _ := store.Get(genKey(i))
					if !ok {
						i = index
					}
					i += uint64(*c)
					totalCount++
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(ch)
	dur := time.Since(start)
	ops := float64(totalCount) / (float64(dur) / 1e9)
	ns := float64(dur) / float64(totalCount*int64(*c))

	if setCount == 0 {
		fmt.Printf("%s setmixed rate: -1 op/s, mean: -1 ns, took: %d s\n", name, int(dur.Seconds()))
	} else {
		fmt.Printf("%s setmixed rate: %.2f op/s, mean: %.2f ns, took: %d s\n", name, float64(setCount)/(float64(dur)/1e9), float64(dur)/float64(setCount), int(dur.Seconds()))
	}
	fmt.Printf("%s getmixed rate: %.2f op/s, mean: %.2f ns, took: %d s\n", name, ops, ns, int(dur.Seconds()))
	fmt.Printf("%s getmixed rate: %.2f op/s, for: %d tx, took: %.2f s\n", name, ops, *txnum, float64(*txnum)/ops)
}

func testDelete(name string, store testdb.Store) {
	var wg sync.WaitGroup
	wg.Add(*c)

	ctx, cancel := context.WithTimeout(context.Background(), *duration)
	defer cancel()

	var totalCount int64
	start := time.Now()
	for j := 0; j < *c; j++ {
		index := uint64(j)
		go func() {
			i := index
		LOOP:
			for {
				select {
				case <-ctx.Done():
					break LOOP
				default:
					if *fix == true && totalCount >= *txnum {
						break LOOP
					}
					store.Del(genKey(i))
					i += uint64(*c)
					totalCount++
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	dur := time.Since(start)
	ops := float64(totalCount) / (float64(dur) / 1e9)
	ns := float64(dur) / float64(totalCount*int64(*c))
	fmt.Printf("%s del rate: %.2f op/s, mean: %.2f ns, took: %d s\n", name, ops, ns, int(dur.Seconds()))
	fmt.Printf("%s del rate: %.2f op/s, for: %d tx, took: %.2f s\n", name, ops, *txnum, float64(*txnum)/ops)
}

func genKey(i uint64) []byte {
	r := make([]byte, *keysize)
	r[0] = 'k'
	//binary.BigEndian.PutUint64(r[0:], i)
	binary.BigEndian.PutUint64(r[1:], i)
	return r
}

func getStore(s string, fsync bool, path string) (testdb.Store, string, error) {
	var store testdb.Store
	var err error
	switch s {
	default:
		err = fmt.Errorf("unknown store type: %v", s)
	case "leveldb":
		if path == "" {
			path = "leveldb.db"
		}
		store, err = testdb.NewLevelDBStore(path, fsync)
	case "badger":
		if path == "" {
			path = "badger.db"
		}
		store, err = testdb.NewBadgerStore(path, fsync)
	}

	return store, path, err
}
