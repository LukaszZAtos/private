#!/bin/bash

TIMEOUT=2

# Function to check SSH connection
check_ssh() {
  ip=$1
  if timeout $TIMEOUT bash -c "</dev/tcp/$ip/22" >/dev/null 2>&1; then
    echo -n "$ip - SSH Successful | "
  elif [ $? -eq 124 ]; then
    echo -n "$ip - SSH Timed Out | "
  else
    echo -n "$ip - SSH Refused | "
  fi
}

# Function to check WinRM connection
check_winrm() {
  ip=$1
  if timeout $TIMEOUT curl --insecure --silent --max-time $TIMEOUT "http://$ip:5985/wsman" >/dev/null 2>&1; then
    echo -n "WinRM Successful | "
  elif [ $? -eq 124 ]; then
    echo -n "WinRM Timed Out | "
  else
    echo -n "WinRM Refused | "
  fi
}

# Function to check WinRM over HTTPS connection
check_winrm_https() {
  ip=$1
  if timeout $TIMEOUT curl --insecure --silent --max-time $TIMEOUT "https://$ip:5986/wsman" >/dev/null 2>&1; then
    echo "WinRM over HTTPS Successful"
  elif [ $? -eq 124 ]; then
    echo "WinRM over HTTPS Timed Out"
  else
    echo "WinRM over HTTPS Refused"
  fi
}

# Read IP addresses from a file (one per line)
while IFS= read -r ip; do
  check_ssh "$ip"
  check_winrm "$ip"
  check_winrm_https "$ip"
done < ip_list.txt

