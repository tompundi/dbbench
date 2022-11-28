#!/usr/bin/env zsh

SIZE=256

# STORES=("badger" "bbolt" "bolt" "leveldb" "kv" "buntdb" "pebble" "pogreb" "nutsdb"  "btree" "btree/memory" "map" "map/memory")

STORES=("badger" "leveldb")

export LD_LIBRARY_PATH=/usr/local/lib

`rm  -fr .*db`
`rm  -fr *.db`
`rm  -fr pogreb.*`
`rm -f benchmarks/test.log`

echo "=========== test nofsync ==========="
for i in "${STORES[@]}"
do
	go run main.go -d 10s -size ${SIZE} -s "$i" >> benchmarks/test.log 2>&1
done

`rm  -fr .*db`
`rm  -fr *.db`
`rm  -fr pogreb.*`

echo ""
echo "=========== test fsync ==========="

for i in "${STORES[@]}"
do
	go run main.go -d 10s -size ${SIZE} -s "$i" -fsync >> benchmarks/test.log 2>&1
done

`rm  -fr .*db` 
`rm  -fr *.db`
`rm  -fr pogreb.*`