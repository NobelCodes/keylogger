from flask import Flask, request, jsonify
import os
import time
import json
from datetime import datetime

app = Flask(__name__)

# Upload Folder
UPLOAD_FOLDER = './downloads'
os.makedirs(UPLOAD_FOLDER, exist_ok=True)

# Path to the JSON file
DATA_FILE = 'data.json'

# Function to append new data to the JSON file
def save_data(new_data):
    # Check if the file exists; if not, create it
    if not os.path.exists(DATA_FILE):
        with open(DATA_FILE, 'w') as f:
            json.dump([], f)  # Initialize with an empty list

    # Load existing data
    with open(DATA_FILE, 'r') as f:
        data = json.load(f)

    # Append new data
    data.append(new_data)

    # Save the updated data
    with open(DATA_FILE, 'w') as f:
        json.dump(data, f, indent=4)

# File upload endpoint with ID
@app.route('/upload/<int:id>', methods=['POST'])
def upload_file(id):
    if 'file' not in request.files:
        return 'No file part', 400
    file = request.files['file']
    if file.filename == '':
        return 'No selected file', 400

    # Extract file extension
    file_ext = os.path.splitext(file.filename)[1]

    # Save file with timestamp and ID as name
    timestamp = time.strftime("%Y%m%d%H%M%S")
    file_path = os.path.join(UPLOAD_FOLDER, f"{timestamp}_{id}{file_ext}")
    file.save(file_path)

    return f"File saved as {file_path}", 200

# POST endpoint to accept JSON data and store it in data.json
@app.route('/post', methods=['POST'])
def post_data():
    if request.is_json:
        # Get the client IP address and current timestamp
        client_ip = request.remote_addr
        timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

        # Get posted data
        data = request.get_json()
        data_with_meta = {
            'ip': client_ip,
            'timestamp': timestamp,
            'data': data
        }

        # Save data to data.json
        save_data(data_with_meta)

        return jsonify({"message": "Data received and stored", "data": data_with_meta}), 201
    return jsonify({"message": "Invalid JSON"}), 400

# GET endpoint to retrieve all data from data.json
@app.route('/data', methods=['GET'])
def get_data():
    # Check if data.json exists
    if os.path.exists(DATA_FILE):
        with open(DATA_FILE, 'r') as f:
            data = json.load(f)
        return jsonify(data), 200
    else:
        return jsonify({"message": "No data found"}), 404

# Sample homepage
@app.route('/')
def hello_world():
    return 'Hello from investigator, I am just out of curiosity. Don’t blame me; I won’t modify anything, just trying to test how it works!'

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)
