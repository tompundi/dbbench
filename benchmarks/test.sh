#!/usr/bin/env zsh

KEYSIZE=16
VALSIZE=256
TXNUM=1000000
FIXTIME=false

STORES=("badger" "leveldb")

rm  -fr .*db
rm  -fr *.db
rm  -fr pogreb.*
rm -f test.log

echo "=========== test nofsync ==========="
for i in "${STORES[@]}"
do
	go run main.go -d=10s -keysize=${KEYSIZE} -valsize=${VALSIZE} -txnum=${TXNUM} -fix=${FIXTIME} -s "$i" >> test.log 2>&1
done

rm  -fr .*db
rm  -fr *.db
rm  -fr pogreb.*

echo ""
echo "=========== test fsync ==========="

for i in "${STORES[@]}"
do
	go run main.go -d=10s -keysize=${KEYSIZE} -valsize=${VALSIZE} -txnum=${TXNUM} -fix=${FIXTIME} -s "$i" -fsync >> test.log 2>&1
done
