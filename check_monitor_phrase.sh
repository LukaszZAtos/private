#!/bin/bash

# Check if a phrase was provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 'phrase_to_search'"
    exit 1
fi

# The phrase to search for
PHRASE="$1"

# Directory to search in
SEARCH_DIR="/home/nadabackup/archive"

# Flag to track if the phrase is found
FOUND=0

# Find .gz files modified in the last 12 hours and search for the phrase
while IFS= read -r -d '' file; do
    if zgrep -iq "$PHRASE" "$file"; then
        echo "Phrase '$PHRASE' found in $file"
        FOUND=1
        break
    fi
done < <(find "$SEARCH_DIR" -type f -name '*.gz' -mmin -720 -print0)

# Check if the phrase was found
if [ "$FOUND" -eq 1 ]; then
    exit 0
else
    echo "Phrase '$PHRASE' not found"
    exit 1
fi
#!/bin/bash

# Check if a phrase was provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 'phrase_to_search'"
    exit 1
fi

# The phrase to search for
PHRASE="$1"

# Directory to search in
SEARCH_DIR="/home/nadabackup/archive"

# Flag to track if the phrase is found
FOUND=0

# Find .gz files modified in the last 12 hours and search for the phrase
while IFS= read -r -d '' file; do
    if zgrep -iq "$PHRASE" "$file"; then
        echo "Phrase '$PHRASE' found in $file"
        FOUND=1
        break
    fi
done < <(find "$SEARCH_DIR" -type f -name '*.gz' -mmin -720 -print0)

# Check if the phrase was found
if [ "$FOUND" -eq 1 ]; then
    exit 0
else
    echo "Phrase '$PHRASE' not found"
    exit 1
fi
