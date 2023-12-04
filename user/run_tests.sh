#!/bin/bash

echo -e "\e[1mWelcome to User Manager Test Suite.\e[0m"
echo "The test suite only focuses on critical functionalities for key scenarios."
echo "Some tests may take longer; your patience is appreciated."
echo "Running tests..."

echo -e "\n"
go test -v ./...
echo -e "\n"

echo -e "\e[1mUser Manager Test Suite has completed."