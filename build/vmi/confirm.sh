#!/bin/bash
 
# Function to display the confirmation prompt
function confirm() {
    while true; do
        read -p "Do you want to proceed? (Y/n) " yn
        case $yn in
            [Yy]* ) return 0;;
            [Nn]* ) return 1;;
            [Cc]* ) exit;;
            * ) echo "Please answer YES, NO, or CANCEL.";;
        esac
    done
}
 
# Example usage of the confirm function
if confirm; then
    echo "User chose YES. Executing the operation..."
    # Place your code here to execute when user confirms
else
    echo "User chose NO. Aborting the operation..."
    # Place your code here to execute when user denies
fi