#!/usr/bin/env zsh

STORES=("badger" "leveldb")

if [ $# != 0 ]
then
    logfile="$1"-test.log
else
    logfile=test.log
fi

rm -f *.csv

echo "name,batch-write time,set,set(txnum) time,get,get(txnum),set-mixed,get-mixed,del" >> benchmarks/nofsync_fix_time.csv
echo "name,batch-write size(MB),set,set(txnum) time,get,get(txnum),set-mixed,get-mixed,del" >> benchmarks/nofsync_fix_throughputs.csv

echo "name,batch write time,set,set(txnum),get,get(txnum),set-mixed,get-mixed,del" >> benchmarks/fsync_fix_time.csv
echo "name,batch-write size(MB),set,set(txnum) time,get,get(txnum),set-mixed,get-mixed,del" >> benchmarks/fsync_fix_throughputs.csv

for i in "${STORES[@]}"
do
    data=$(grep -e ^"${i}"/nofsync  benchmarks/${logfile}|awk '{print $10}'|xargs| tr ' ' ',')
    echo "${i}/nofsync,${data}" >> benchmarks/nofsync_fix_time.csv
    data=$(grep -e ^"${i}"/nofsync  benchmarks/${logfile}|awk '{print $4}'|xargs| tr ' ' ',')
    echo "${i}/nofsync,${data}" >> benchmarks/nofsync_fix_throughputs.csv

    data=$(grep -e ^"${i}"/fsync  benchmarks/${logfile}|awk '{print $7}'|xargs| tr ' ' ',')
    echo "${i}/fsync,${data}" >> benchmarks/fsync_fix_time.csv
    data=$(grep -e ^"${i}"/fsync  benchmarks/${logfile}|awk '{print $4}'|xargs| tr ' ' ',')
    echo "${i}/fsync,${data}" >> benchmarks/fsync_fix_throughputs.csv
done