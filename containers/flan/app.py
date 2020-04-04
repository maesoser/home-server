from flask import Flask
from flask import request, jsonify
#from .flan import nmap

doing_scans = []
done_scans = []

app = Flask(__name__)
app.config["DEBUG"] = False

@app.route("/api/scans", methods=['GET'])
def get_scan():
    if 'state' in request.args:
        if request.args['state'] == 'doing':
            return jsonify(doing_scans)
        elif request.args['state'] == 'done':
            return jsonify(done_scans)
    return jsonify(doing_scans)

#@app.route("/api/scans", methods=['POST'])
#def add_scan():
#    data = json.loads(request.data)
#    if 'targets' in data:
#        nmap.do_scan(data['targets'])
#    return "Scan Added"

#@app.errorhandler(404)
#def page_not_found(e):
#    return "<h1>404</h1><p>The resource could not be found.</p>", 404

@app.route("/", methods=['GET'])
def get_web():
    return "MainWeb"
