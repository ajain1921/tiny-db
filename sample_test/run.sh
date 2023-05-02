#!/bin/bash

code_folder='src/'
log_folder='logs/'
curr_folder=$(pwd)'/'
servers=(A B C D E)
client='../bin/client'
server='../bin/server'

ports=(10001 10002 10003 10004 10005)

# shutdown servers
shutdown() {
	# for port in ${ports[@]};`` do
		# kill -15 $(lsof -ti:$port)
		pkill server
	# done
}
trap shutdown EXIT

rm -rf ${log_folder} && mkdir ${log_folder}
rm -rf *.log

# initialize servers``
cd ..
make client && make server
cd sample_test


for s in ${servers[@]}; do
	$server $s config.txt > ${log_folder}server_${s}.log 2>&1 &
done

sleep 4

# run 2 tests
# timeout -s SIGTERM 5s $client a config.txt < ${curr_folder}input1.txt > ${curr_folder}output1.log 2>&1
# timeout -s SIGTERM 5s $client a config.txt < ${curr_folder}input2.txt > ${curr_folder}output2.log 2>&1

$client a config.txt < ${curr_folder}input1.txt > ${curr_folder}output1.log 2>&1 &
$client b config.txt < ${curr_folder}input2.txt > ${curr_folder}output2.log 2>&1 &

sleep 100

cd $curr_folder
echo "Difference between your output and expected output:"
diff output1.log expected1.txt
diff output2.log expected2.txt
