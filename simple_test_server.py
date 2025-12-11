from flask import Flask, request, jsonify

app = Flask(__name__)

@app.route("/", methods=["GET"])
def index():
    return '<button onclick="fetch(`/echo`, {method: `post`, body: {}})">heck</button>'

@app.route("/echo", methods=["POST"])
def echo():
    data = request.get_json(silent=True)

    if data is None:
        return jsonify({"received": request.data.decode()}), 200
    
    return jsonify({"json": data}), 200


if __name__ == "__main__":
    print("running at http://localhost:8999")
    app.run(host="localhost", port=8999, threaded=True)
