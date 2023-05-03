#!/bin/bash

# Path to the file containing client commands and inputs
input_file="input.txt"
client_path='../bin/client'
server='../bin/server'
servers=(A B C D E)
log_folder='logs/'

rm -rf ${log_folder} && mkdir ${log_folder}
rm -rf *.log

cd ..
make client && make server
cd another_test

# Declare an associative array to store named pipes for each client
client_pipes=()

for s in ${servers[@]}; do
	$server $s config.txt > ${log_folder}server_${s}.log 2>&1 &
done

client_num=3
# Create named pipes for each client and store them in the array
for i in $(seq 1 $client_num); do
    pipe="/tmp/${i}_pipe"
    if [[ ! -p "$pipe" ]]; then
        mkfifo "$pipe"
    fi
    client_pipes+="$pipe"
done

echo $client_pipes

pids=""
# Start the clients
for client in "${!client_pipes[@]}"; do
    echo "YO: ${client_pipes[$client]}"
    $client_path "$client" config.txt < "${client_pipes[$client]}" > output_$client.log &
    sleep 0.1
    pids="$pids $!"
done

# read from pipe and output
read_from_pipe() {
  while true; do
    if read line < "${client_pipes[0]}"; then
      echo "Result: $line"
    fi
    sleep 0.5
  done
}

# Call the function in the background
read_from_pipe &

# Read each line from the file and provide input to the correct client
while read -r line; do
    # Extract the client name and input
    client=$(echo "$line" | awk '{print $1}')
    input=$(echo "$line" | awk '{$1=""; print $0}')

    echo $client $input

    # Provide the input to the appropriate named pipe
    if [[ -p "${client_pipes[$client]}" ]]; then
        echo "$input" > "${client_pipes[$client]}"
    else
        echo "Invalid client name: $client"
    fi

    sleep 0.25
done < "$input_file"

shutdown() {
    pkill server
    pkill client
    for pipe in "${client_pipes[@]}"; do
        rm "$pipe"
    done
}
trap shutdown EXIT

wait $pids



# Remove the named pipes
