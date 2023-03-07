from flask import Flask, request
import json

app = Flask(__name__)

@app.route('/')
def home():
    return "App Works!!!"

@app.route("/server_request")
def server_request():
    print(request.args.get("param"))
    return "served"
