#!/bin/bash

# Function to check the connection status for a given IP and port
check_port() {
    local ip=$1
    local port=$2
    local status

    # Try to connect and wait for up to 3 seconds for a response
    nc -z -w3 $ip $port >/dev/null 2>&1

    # Check the exit code of the previous command
    if [ $? -eq 0 ]; then
        status="Success"
    else
        status=$(nc -z -w3 $ip $port 2>&1)
        if [[ "$status" =~ "succeeded!" ]]; then
            status="Refused"
        else
            status="Timed out"
        fi
    fi

    echo "$status"
}

# Check if the filename is provided as an argument
if [ $# -eq 0 ]; then
    echo "Usage: $0 <filename>"
    exit 1
fi

# Check if the file exists
if [ ! -f "$1" ]; then
    echo "Error: File $1 not found."
    exit 1
fi

# Loop through the file, extracting IPs and checking the ports
while read -r ip; do
    ssh_status=$(check_port "$ip" 22)       # SSH
    winrm_status=$(check_port "$ip" 5985)   # WinRM
    winrm_https_status=$(check_port "$ip" 5986)  # WinRM over HTTPS

    echo "$ip: SSH($ssh_status) WinRM($winrm_status) WinRM over HTTPS($winrm_https_status)"
done < "$1"
