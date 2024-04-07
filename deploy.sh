#!/bin/bash

git pull --rebase origin main
docker compose down
docker compose up -d
echo "Script execution completed."
