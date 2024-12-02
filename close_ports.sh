#!/bin/bash

# Ports to check
PORTS=("8081" "8082" "8083" "8084")

echo "Looking for processes on ports: ${PORTS[*]}"

# Iterate over each port and find processes
for PORT in "${PORTS[@]}"; do
    # Find the process ID (PID) using lsof
    PIDS=$(lsof -ti :"$PORT")
    
    if [ -n "$PIDS" ]; then
        echo "Found process(es) on port $PORT: $PIDS"
        # Kill the processes
        kill -9 $PIDS
        echo "Killed process(es) on port $PORT."
    else
        echo "No processes found on port $PORT."
    fi
done

echo "All specified ports are now free."
