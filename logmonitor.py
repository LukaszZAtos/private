import os
import subprocess
import pyinotify

# Configuration
log_file_path = "/var/log/syslog"  # Adjust the path to your syslog file
target_word = "WORD"  # Replace with the word you want to trigger on
command_to_execute = "your_command_here"  # Replace with the Linux command to execute

# Define a callback function to execute the command
def execute_command(event):
    if target_word in event.pathname:
        print(f"Triggered by {target_word}. Executing command: {command_to_execute}")
        subprocess.Popen(command_to_execute, shell=True)

# Initialize the inotify watcher
wm = pyinotify.WatchManager()
mask = pyinotify.IN_MODIFY  # Watch for file modification events

# Create a Notifier with the callback function
notifier = pyinotify.Notifier(wm, execute_command)

# Add a watch for the syslog file
wm.add_watch(log_file_path, mask)

try:
    print(f"Monitoring syslog for '{target_word}'... Press Ctrl+C to exit.")
    notifier.loop()
except KeyboardInterrupt:
    print("Monitoring stopped.")
