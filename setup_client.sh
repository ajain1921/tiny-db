#!/bin/bash

for i in {1,2,4}
do
	echo "Copying code to $i server"
	sshpass -p $UIUC_PASSWORD scp -o StrictHostKeyChecking=no -r ./mp0 "$netid@sp23-cs425-220$i.cs.illinois.edu:/home/$netid"
	sshpass -p $UIUC_PASSWORD ssh -o StrictHostKeyChecking=no "$netid@sp23-cs425-220$i.cs.illinois.edu" "ls mp0"
done



for i in {1,2,4}
do
	echo "SSHing to $i server"
	sshpass -p $UIUC_PASSWORD ssh -o StrictHostKeyChecking=no -n -f "$netid@sp23-cs425-220$i.cs.illinois.edu" "sh -c 'cd mp0 && make && nohup python3 -u generator.py 2 | ./bin/node node$i sp23-cs425-2209.cs.illinois.edu 4321 > /dev/null 2>&1 &'"
done
