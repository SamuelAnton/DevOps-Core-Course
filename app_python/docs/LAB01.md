# Lab 01

## Framework Selection

I choosed `Flask ` because, how it was stated, it is lightweight and I have already worked with it in my project, so it will be easier to use it than study new tools (even if I think that study new tools is great, I want to practice with Flask more).

| Framework | Pros | Cons |
|----------|--------|-------------|
| **Flask** | Lightweight, already know from previous project | Not so good choice for complex projects |
| **FastAPI** | Modern, async, auto-documentation | Complex and requires learning it |
| **Django** | Full-featured, includes ORM | Too complex, requires a lot to learn |


## Best Practices Applied
### Clean Code Organization
Clear code is easier to maintain and read. It is espessialy important if you work in group, so you can spend more time on work but not on understanding in what is going on in this part of code. Also after some time, when you return to your project, you will forgot almost all, so clean code will help you spend less time understanding it. Here are parts of the organization:

- Clear function names
``` python
def get_system_info():
    ...

def get_uptime():
    ...
```
- Proper imports grouping
``` python
import os
import socket
import platform
from datetime import datetime, timezone
from flask import Flask, jsonfy, request
import logging
```
- Comments only where needed
```
"""
DevOps Info Service
Main application module
"""
import os
...

# Configuration
HOST = os.getenv('HOST', '0.0.0.0')
...

# Application start time
start_time = datetime.now()
...

# Logging
logging.basicConfig(
...
```
- Follow PEP 8
```
# Function that collects system info
def get_system_info():
    """Collect system information."""
    return {
        'hostname': socket.gethostname(),
        'platform': platform.system(),
        'architecture': platform.machine(),
        'python_version': platform.python_version()
    }

# Function that gets total uptime of a service
def get_uptime():
    delta = datetime.now() - start_time
    seconds = int(delta.total_seconds())
    hours = seconds // 3600
    minutes = (seconds % 3600) // 60
    return {
        'seconds': seconds,
        'human': f"{hours} hours, {minutes} minutes"
    }
```

### Error Handling
Error handling is crucial, because it will help you find that your service work correctly or not. Good error handling will also allow you to find mistakes in code when testing or using.
``` python
@app.errorhandler(404)
def not_found(error):
    return jsonify({
        'error': 'Not Found',
        'message': 'Endpoint does not exist'
    }), 404

@app.errorhandler(500)
def internal_error(error):
    return jsonify({
        'error': 'Internal Server Error',
        'message': 'An unexpected error occurred'
    }), 500
```

### Logging
Logging is also crucial part of work. It will show you what is happening in the code: events occured, errors met and etc. It will help you to find out hidden mistakes and be sure that everything is working as it should.
``` python
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

logger.info('Application starting...')
logger.debug(f'Request: {request.method} {request.path}')
```

### Dependencies 
It is good practice to collect all project dependencies in a special file, so it will be easier to maintain and download to different machine.
```
Flask==3.1.0
```

### Git Ignore
It is good practice not to push some junk files to your repo, so it will not mess everything up.
```
# Python
__pycache__/
*.py[cod]
venv/
*.log

# IDE
.vscode/
.idea/

# OS
.DS_Store
```

## API Documentation
### GET /

### GET /health

### Testing


## Testing evidence
### Screenshots
### Terminal output

## Challenges & Solutions

## GitHub Community
### Why Stars Matter:

#### Discovery & Bookmarking:

- Stars help you bookmark interesting projects for later reference
- Star count indicates project popularity and community trust
- Starred repos appear in your GitHub profile, showing your interests

#### Open Source Signal:

- Stars encourage maintainers (shows appreciation)
- High star count attracts more contributors
- Helps projects gain visibility in GitHub search and recommendations

#### Professional Context:

- Shows you follow best practices and quality projects
- Indicates awareness of industry tools and trends

### Why Following Matters:

#### Networking:

- See what other developers are working on
- Discover new projects through their activity
- Build professional connections beyond the classroom

#### Learning:

- Learn from others' code and commits
- See how experienced developers solve problems
- Get inspiration for your own projects


#### Collaboration:

- Stay updated on classmates' work
- Easier to find team members for future projects
- Build a supportive learning community


#### Career Growth:

- Follow thought leaders in your technology stack
- See trending projects in real-time
- Build visibility in the developer community


#### GitHub Best Practices:

- Star repos you find useful (not spam)
- Follow developers whose work interests you
- Engage meaningfully with the community
- Your GitHub activity shows employers your interests and involvement