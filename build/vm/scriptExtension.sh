#!/bin/bash

# parse options as --option optionValue
POSITIONAL_ARGS=()

while [[ $# -gt 0 ]]; do
case $1 in
    -a|--artifacts-uri)
    artifacts_uri="$2"
    shift # past argument
    shift # past value
    ;;
    -*|--*)
    echo "Unknown option $1"
    exit 1
    ;;
    *)
    POSITIONAL_ARGS+=("$1") # save positional arg
    shift # past argument
    ;;
esac
done
set -- "${POSITIONAL_ARGS[@]}" # restore positional parameters

function guard_against_empty() {
    # parse options as --option optionValue
    POSITIONAL_ARGS=()

    while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--value)
        value="$2"
        shift # past argument
        shift # past value
        ;;
        -m|--error-message)
        message="$2"
        shift # past argument
        shift # past value
        ;;
        -*|--*)
        echo "Unknown option $1"
        exit 1
        ;;
        *)
        POSITIONAL_ARGS+=("$1") # save positional arg
        shift # past argument
        ;;
    esac
    done
    set -- "${POSITIONAL_ARGS[@]}" # restore positional parameters

    if [ -z "$value" ]; then
        echo "$message"
        exit 1
    fi
}

guard_against_empty --value "$artifacts_uri" --error-message "Artifacts URI is required. Use --artifacts-uri|-a to set."

# write file
echo "artifacts_uri" | sudo tee $MODM_HOME/artifacts.uri