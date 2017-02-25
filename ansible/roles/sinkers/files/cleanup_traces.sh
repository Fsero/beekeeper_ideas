#!/usr/bin/env bash
set -euo pipefail

FREE_SPACE_THRESHOLD=90
MOUNT_TO_WATCH="/"
TRACES_DIR="/var/log/traces"
LAST_MODIFIED_FILES_TO_KEEP_IN_MINUTES="60"


get_free_disk_left() {
    percent=$(df -h ${MOUNT_TO_WATCH} --output=pcent | tail -1 | xargs | tr -d '%')
    FREE_SPACE=$((100-percent))
    echo $FREE_SPACE
}

check_if_cleanup_is_needed () {
    if [[ $FREE_SPACE -le $FREE_SPACE_THRESHOLD ]]; then
        find $TRACES_DIR -mmin +$LAST_MODIFIED_FILES_TO_KEEP_IN_MINUTES -print0 | xargs -0 -L 5000 rm &> /dev/null
    fi
}

get_free_disk_left
check_if_cleanup_is_needed