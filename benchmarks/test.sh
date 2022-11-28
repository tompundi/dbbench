#!/usr/bin/env zsh

SIZE=256

STORES=("badger" "leveldb")

export LD_LIBRARY_PATH=/usr/local/lib

`rm  -fr .*db`
`rm  -fr *.db`
`rm  -fr pogreb.*`
`rm -f test.log`

echo "=========== test nofsync ==========="
for i in "${STORES[@]}"
do
	go run main.go -d 5s -size ${SIZE} -s "$i" >> test.log 2>&1
done

`rm  -fr .*db`
`rm  -fr *.db`
`rm  -fr pogreb.*`

echo ""
echo "=========== test fsync ==========="

for i in "${STORES[@]}"
do
	go run main.go -d 5s -size ${SIZE} -s "$i" -fsync >> test.log 2>&1
done

`rm  -fr .*db` 
`rm  -fr *.db`
`rm  -fr pogreb.*`