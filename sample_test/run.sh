#!/bin/bash

code_folder='src/'
curr_folder=$(pwd)'/'
servers=(A B C D E)
ports=(10001 10002 10003 10004 10005)

# shutdown servers
shutdown() {
	for port in ${ports[@]}; do
		kill -15 $(lsof -ti:$port)
	done
}
trap shutdown EXIT

# initialize servers
cp config.txt ${code_folder}/config.txt
cd $code_folder
for server in ${servers[@]}; do
	./server $server config.txt > ../server_${server}.log 2>&1 &
done

# run 2 tests
timeout -s SIGTERM 5s ./client a config.txt < ${curr_folder}input1.txt > ${curr_folder}output1.log 2>&1
timeout -s SIGTERM 5s ./client a config.txt < ${curr_folder}input2.txt > ${curr_folder}output2.log 2>&1

cd $curr_folder
echo "Difference between your output and expected output:"
diff output1.log expected1.txt
diff output2.log expected2.txt
