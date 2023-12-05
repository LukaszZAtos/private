#!/bin/bash

# Script to set umask 0022 for user ansiblegtsadm

USER_NAME="ansiblegtsadm"
UMASK_VALUE="0022"
USER_HOME=$(getent passwd "$USER_NAME" | cut -d: -f6)

# Function to set umask in a file
set_umask() {
    local file=$1
    if [ -f "$file" ]; then
        echo "Setting umask in $file"
        grep -q "umask" "$file" && sed -i "/umask/c\umask $UMASK_VALUE" "$file" || echo "umask $UMASK_VALUE" >> "$file"
    else
        echo "$file not found."
    fi
}

# Set umask for user-specific files
set_umask "$USER_HOME/.bashrc"
set_umask "$USER_HOME/.bash_profile"
set_umask "$USER_HOME/.profile"

# Set umask for global files (requires root permission)
set_umask "/etc/profile"
set_umask "/etc/bash.bashrc"

echo "umask set to $UMASK_VALUE for user $USER_NAME and globally."
