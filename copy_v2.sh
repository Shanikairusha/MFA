#!/bin/bash

# Define paths
FILE_LIST="copy_list.txt"
SOURCE_DIR="/opt/webapp/DOXPRO/MFA"
DEST_DIR="/mnt/edas1/data/2022_all"
LOG_DIR="rlogs"
mkdir -p "$LOG_DIR"
LOG_FILE="$LOG_DIR/copy_$(date +'%Y%m%d_%H%M%S').log"

# Check if file list exists
if [ ! -f "$FILE_LIST" ]; then
    echo "Error: File list $FILE_LIST not found!" | tee -a "$LOG_FILE"
    exit 1
fi

# Copy missing files
count_total=0
count_copied=0
count_skipped=0
count_deleted=0
count_failed=0
count_renamed=0

while IFS= read -r relative_path; do
    src="$SOURCE_DIR/$relative_path"
    dest="$DEST_DIR/$relative_path"
    ((count_total++))

    # Create destination directory if needed
    mkdir -p "$(dirname "$dest")"

    # If destination file exists
    if [ -f "$dest" ]; then
        if [ ! -s "$dest" ]; then
            echo "DELETE: $relative_path is 0-byte, removing before copy" >> "$LOG_FILE"
            if rm -f "$dest"; then
                ((count_deleted++))
            else
                # Deletion failed, prepare renamed destination
                dir_path="$(dirname "$dest")"
                filename="$(basename "$dest")"
                dest="$dir_path/v2_$filename"
                echo "ERROR: Failed to delete $relative_path, will copy as v2_$filename" >> "$LOG_FILE"
                ((count_renamed++))
            fi
        else
            echo "SKIP: $relative_path exists and is not empty" >> "$LOG_FILE"
            ((count_skipped++))
            continue
        fi
    fi

    # Attempt to copy the file
    if cp "$src" "$dest"; then
        # Confirm file copied correctly
        if [ -s "$dest" ]; then
            echo "COPIED: $src -> $dest" >> "$LOG_FILE"
            ((count_copied++))
        else
            echo "ERROR: $relative_path copied as 0-byte, deleting" >> "$LOG_FILE"
            rm -f "$dest"
            ((count_failed++))
        fi
    else
        echo "ERROR: Failed to copy $relative_path" >> "$LOG_FILE"
        ((count_failed++))
    fi
done < "$FILE_LIST"

echo "Copy complete. Total: $count_total | Copied: $count_copied | Skipped: $count_skipped | Deleted 0-byte: $count_deleted | Renamed: $count_renamed | Failed: $count_failed"
echo "Logs saved to: $LOG_FILE"