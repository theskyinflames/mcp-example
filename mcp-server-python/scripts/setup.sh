#!/bin/bash

set -e

# Create a virtual environment
python3 -m venv venv

# Activate the virtual environment
source venv/bin/activate  # On macOS/Linux

# Install the required packages
pip install -r requirements.txt