#!/bin/bash

# Check if a filename was given as an argument
if [ -n "$1" ]; then
    # If a filename is provided, add that specific file to git
    echo "--------------------------------------------------------------"
    echo "Adding file: $1"
    git add "$1"
else
    # If no filename is provided, add all files
    echo "--------------------------------------------------------------"
    echo "Adding all modified files..."
    git add .
fi

# Prompt for a commit message
echo "--------------------------------------------------------------"
echo "Enter commit message: "
read commit_message

# Check if commit message is empty
if [ -z "$commit_message" ]; then
    echo "Aborting commit: No commit message provided."
    exit 1
fi

# Commit the changes
echo "Committing changes..."
git commit -m "$commit_message"

# Pull the latest changes from the remote repository
echo "Pulling the latest changes from the remote repository..."
git pull --rebase

# Push the changes to the remote repository
echo "Pushing the changes to the remote repository..."
git push

# Check if the push was successful
if [ $? -eq 0 ]; then
    echo "Push successful!"
else
    echo "Push failed. Please check for errors."
fi









































