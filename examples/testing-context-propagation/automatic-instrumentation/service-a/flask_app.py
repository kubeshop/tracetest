from flask import Flask, request, jsonify, make_response
import requests
import os

app = Flask(__name__)

api_port = os.getenv('API_PORT', '8800')
service_b_address = os.getenv('SERVICE_B_URL', 'http://service-b:8801')

@app.route('/sendData', methods=['POST'])
def send_data():
  data = request.json

  data['messageFromA'] = 'Hello from Service A'

  response = requests.post(f'{service_b_address}/augmentData', json=data)
  if response.status_code == 200:
    augmented_data = response.json()
    data.update(augmented_data)
  else:
    return make_response(jsonify({'error': 'Failed to augment data'}), 500)

  return jsonify(data)

if __name__ == '__main__':
  print('Running on port: ' + api_port)
  app.run(host='0.0.0.0', port=api_port)
