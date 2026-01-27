# Overview
This service is a study service in python that returns comprehensive service and system information.

# Prerequisites
python version
dependencies (from requirements.txt):
- 

# Installation
```
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

# Running the Application
```
python app.py
# Or with custom config
PORT=8080 python app.py
HOST=127.0.0.1 PORT=3000 python app.py
```

# API Endpoint
- `GET /` - Service and system information
- `GET /health` - Health check

# Configuration 
| Variable name | Basic value | Description |
|----------|--------|-------------|
| **Host** | 0.0.0.0 | IP of the service |
| **Port** | 5000 | Port of the service |
| **Debug** | False | Debug mode enabeled |
