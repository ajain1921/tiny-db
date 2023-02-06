#!/bin/bash
echo "Copy code"
sshpass -p $UIUC_PASSWORD scp -o StrictHostKeyChecking=no -r ./mp0 "$netid@sp23-cs425-2209.cs.illinois.edu:/home/$netid"
echo "Sshing to server..."
sshpass -p $UIUC_PASSWORD ssh -o StrictHostKeyChecking=no "$netid@sp23-cs425-2209.cs.illinois.edu" "cd mp0 && rm -f data.csv && make && ./bin/logger 4321"