#!/usr/bin/env zsh

if [ $# != 0 ]
then
    logfile="$1"-test.log
else
    logfile=test.log
fi

`rm -f *.csv`

STORES=("badger" "leveldb")

echo "name,set,get,set-mixed,get-mixed,del" >> nofsync_throughputs.csv
echo "name,set,get,set-mixed,get-mixed,del" >> nofsync_time.csv
echo "name,set,get,set-mixed,get-mixed,del" >> fsync_throughputs.csv
echo "name,set,get,set-mixed,get-mixed,del" >> fsync_time.csv

for i in "${STORES[@]}"
do
    data=`grep -e ^${i}/nofsync  ${logfile}|awk '{print $4}'|xargs| tr ' ' ','`
    echo "${i}/nofsync,${data}" >> nofsync_throughputs.csv
    data=`grep -e ^${i}/nofsync  ${logfile}|awk '{print $7}'|xargs| tr ' ' ','`
    echo "${i}/nofsync,${data}" >> nofsync_time.csv

    data=`grep -e ^${i}/fsync  ${logfile}|awk '{print $4}'|xargs| tr ' ' ','`
    echo "${i}/fsync,${data}" >> fsync_throughputs.csv
    data=`grep -e ^${i}/nofsync  ${logfile}|awk '{print $7}'|xargs| tr ' ' ','`
    echo "${i}/fsync,${data}" >> fsync_time.csv
done