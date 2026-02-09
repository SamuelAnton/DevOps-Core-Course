"""
DevOps Info Service
Main application module
"""
import os
import socket
import platform
from datetime import datetime, timezone
from flask import Flask, jsonify, request
import logging

app = Flask(__name__)

# Configuration
HOST = os.getenv('HOST', '0.0.0.0')
PORT = int(os.getenv('PORT', 5000))
DEBUG = os.getenv('DEBUG', 'False').lower() == 'true'

# Application start time
start_time = datetime.now()

# Logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


# Function that collects system info
def get_system_info():
    """Collect system information."""
    return {
        'hostname': socket.gethostname(),
        'platform': platform.system(),
        'platform_version': platform.version(),
        'architecture': platform.machine(),
        'cpu_count': os.cpu_count(),
        'python_version': platform.python_version()
    }


# Function that gets total uptime of a service
def get_uptime():
    delta = datetime.now() - start_time
    seconds = int(delta.total_seconds())
    hours = seconds // 3600
    minutes = (seconds % 3600) // 60
    human = f"{hours} hour{'s' if hours != 1 else ''}, {minutes} minute{'s' if minutes != 1 else ''}"
    return {
        'seconds': seconds,
        'human': f"{hours} hour{'s' if hours != 1 else ''}, {minutes} minute{'s' if minutes != 1 else ''}"
    }


# API Endpoints
# Main - service and system information
@app.route('/')
def index():
    """Main endpoint - service and system information."""
    # Collect all required information
    system_info = get_system_info()
    uptime_info = get_uptime()
    
    response = {
        "service": {
            "name": "devops-info-service",
            "version": "1.0.0",
            "description": "DevOps course info service",
            "framework": "Flask"
        },
        "system": {
            "hostname": system_info['hostname'],
            "platform": system_info['platform'],
            "platform_version": system_info['platform_version'],
            "architecture": system_info['architecture'],
            "cpu_count": system_info['cpu_count'],
            "python_version": system_info['python_version']
        },
        "runtime": {
            "uptime_seconds": uptime_info['seconds'],
            "uptime_human": uptime_info['human'],
            "current_time": datetime.now(timezone.utc).isoformat().replace('+00:00', 'Z'),
            "timezone": "UTC"
        },
        "request": {
            "client_ip": request.remote_addr,
            "user_agent": request.headers.get('User-Agent', 'Unknown'),
            "method": request.method,
            "path": request.path
        },
        "endpoints": [
            {"path": "/", "method": "GET", "description": "Service information"},
            {"path": "/health", "method": "GET", "description": "Health check"}
        ]
    }
    
    return jsonify(response)


# Health check
@app.route('/health')
def health():
    return jsonify({
        'status': 'healthy',
        'timestamp': datetime.now(timezone.utc).isoformat().replace('+00:00', 'Z'),
        'uptime_seconds': get_uptime()['seconds']
    })


# Error Handling
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


# Start of the service
if __name__ == '__main__':
    logger.info('Application starting...')
    logger.info(f'Running on http://{HOST}:{PORT}')
    logger.info(f'Debug mode: {DEBUG}')
    app.run(host=HOST, port=PORT, debug=DEBUG)