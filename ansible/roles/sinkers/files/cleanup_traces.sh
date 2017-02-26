#!/usr/bin/env bash
set -euo pipefail

FREE_SPACE_THRESHOLD=10
MOUNT_TO_WATCH="/"
TRACES_DIR="/var/log/traces"
LAST_MODIFIED_FILES_TO_KEEP_IN_MINUTES="60"


get_free_disk_left() {
    percent=$(df -h ${MOUNT_TO_WATCH} --output=pcent | tail -1 | xargs | tr -d '%')
    FREE_SPACE=$((100-percent))
}

check_if_cleanup_is_needed () {
    if [[ $FREE_SPACE -le $FREE_SPACE_THRESHOLD ]]; then
        local files_to_remove=$(find $TRACES_DIR -mmin +$LAST_MODIFIED_FILES_TO_KEEP_IN_MINUTES -print)
        logger -t info "[cleanup] removing old traces"
        logger -t info "[cleanup] removing $files_to_remove"
        find $TRACES_DIR -mmin +$LAST_MODIFIED_FILES_TO_KEEP_IN_MINUTES -print0 | xargs -0 -L 5000 rm &> /dev/null
    fi
}

get_free_disk_left
check_if_cleanup_is_needed