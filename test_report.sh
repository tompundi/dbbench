#!/usr/bin/env zsh

SIZE=256

# STORES=("badger" "bbolt" "bolt" "leveldb" "kv" "buntdb" "pebble" "pogreb" "nutsdb" "btree" "btree/memory" "map" "map/memory")

STORES=("badger" "leveldb")

`rm  -fr .*db`
`rm  -fr *.db`
`rm  -fr pogreb.*`
`rm -f benchmarks/test.log`

echo "=========== test nofsync ==========="
for i in "${STORES[@]}"
do
	go run main.go -d 1m -size ${SIZE} -s "$i" >> benchmarks/test.log 2>&1
done

`rm  -fr .*db`
`rm  -fr *.db`
`rm  -fr pogreb.*`

echo ""
echo "=========== test fsync ==========="

for i in "${STORES[@]}"
do
	go run main.go -d 1m -size ${SIZE} -s "$i" -fsync >> benchmarks/test.log 2>&1
done

`rm  -fr .*db` 
`rm  -fr *.db`
`rm  -fr pogreb.*`


logfile=test.log

`rm -f benchmarks/*.csv`

echo "name,set,get,set-mixed,get-mixed,del" >> benchmarks/nofsync_throughputs.csv
echo "name,set,get,set-mixed,get-mixed,del" >> benchmarks/nofsync_time.csv
echo "name,set,get,set-mixed,get-mixed,del" >> benchmarks/fsync_throughputs.csv
echo "name,set,get,set-mixed,get-mixed,del" >> benchmarks/fsync_time.csv

for i in "${STORES[@]}"
do
    data=`grep -e ^${i}/nofsync  benchmarks/${logfile}|awk '{print $4}'|xargs| tr ' ' ','`
    echo "${i}/nofsync,${data}" >> benchmarks/nofsync_throughputs.csv
    data=`grep -e ^${i}/nofsync  benchmarks/${logfile}|awk '{print $7}'|xargs| tr ' ' ','`
    echo "${i}/nofsync,${data}" >> benchmarks/nofsync_time.csv

    data=`grep -e ^${i}/fsync  benchmarks/${logfile}|awk '{print $4}'|xargs| tr ' ' ','`
    echo "${i}/fsync,${data}" >> benchmarks/fsync_throughputs.csv
    data=`grep -e ^${i}/nofsync  benchmarks/${logfile}|awk '{print $7}'|xargs| tr ' ' ','`
    echo "${i}/fsync,${data}" >> benchmarks/fsync_time.csv
done