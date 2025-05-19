#!/usr/bin/env bash

RED="\e[31m"
YELLOW="\e[33m"
GREEN="\e[32m"
BLUE="\e[34m"
GRAY="\e[90m"
RESET="\e[0m"

ERROR_LOG="${RED}[ERROR]${RESET}"
WARN_LOG="${YELLOW}[WARN]${RESET}"
INFO_LOG="${GREEN}[INFO]${RESET}"
PRINT_INDENT_LOG="${GRAY}||===>${RESET}"

log_err() {
    echo -e "$ERROR_LOG $1" >&2
}

log_warn() {
    echo -e "$WARN_LOG $1" >&2
}

log_info() {
    echo -e "$INFO_LOG $1" >&2
}

log_print() {
    echo -e "$PRINT_INDENT_LOG $1" >&2
}

check_go_licenses() {
    command -v go-licenses >/dev/null 2>&1 || {
        log_warn "go-licenses not found. Installing..."
        go install github.com/google/go-licenses@latest || {
            log_err "Failed to install go-licenses"
            exit 1
        }
    }
}

run_license_check() {
    log_info "${YELLOW}Running go-licenses check...${RESET}"

    tmp_output="$(mktemp)"
    
    # Run go-licenses check, excluding the current repository and storing the output
    if ! go-licenses check ./... --ignore=github.com/Phillezi/kthcloud-cli > "$tmp_output" 2>&1; then
        log_err "Failed to generate license report:"
        cat "$tmp_output" | while IFS= read -r line; do
            log_print "\t$line"
        done
        rm -f "$tmp_output"
        return 1
    fi

    log_print "${GREEN}License check passed! No forbidden licenses found.${RESET}"
    rm -f "$tmp_output"
    return 0
}

run_license_report() {
    log_info "${YELLOW}Generating license report...${RESET}"
    
    # Using go-licenses with custom markdown template
    if ! go-licenses report ./... --ignore=github.com/Phillezi/kthcloud-cli --template=./scripts/util/.res/report.md.template > ./licenses/license_report.md; then
        log_err "Failed to generate license report"
        return 1
    fi
    
    log_print "${GREEN}License report generated successfully as 'licenses/license_report.md'.${RESET}"
    return 0
}

run_license_save() {
    log_info "${YELLOW}Generating license copyrights and notices...${RESET}"
    
    # Using go-licenses with custom markdown template
    if ! go-licenses save ./... --ignore=github.com/Phillezi/kthcloud-cli --save_path="./licenses/third_party"; then
        log_err "Failed to generate license copyrights and notices"
        return 1
    fi
    
    log_print "${GREEN}License copyrights and notices generated successfully at 'licenses/third_party/'.${RESET}"
    return 0
}


pushd "$(dirname "${BASH_SOURCE[0]}")/../.." > /dev/null || {
    log_err "Failed to change to script directory."
    exit 1
}

mkdir -p licenses

check_go_licenses

if ! run_license_check; then
    echo -e "${RED}License check failed.${RESET}"
    popd > /dev/null
    exit 1
fi

if ! run_license_report; then
    echo -e "${RED}License report generation failed.${RESET}"
    popd > /dev/null
    exit 1
fi

if ! run_license_save; then
    echo -e "${RED}License save failed.${RESET}"
    popd > /dev/null
    exit 1
fi

echo -e "${GREEN}All license checks and reports completed successfully.${RESET}"
popd > /dev/null
exit 0
