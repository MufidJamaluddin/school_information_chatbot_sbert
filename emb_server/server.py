from sentence_transformers import SentenceTransformer
from flask import jsonify, Flask
from gevent.pywsgi import WSGIServer

ori_model_path = 'model/ori_model'
fine_tuned_model_path = 'model/fine_tuned'

ori_model = SentenceTransformer(ori_model_path)
fine_tuned_model = SentenceTransformer(fine_tuned_model_path)

app = Flask(__name__)

@app.route("/encode-ori/<text>")
async def encode_ori(text):
    return jsonify(ori_model.encode(text).tolist())

@app.route("/encode-ft/<text>")
async def encode_ft(text):
    return jsonify(fine_tuned_model.encode(text).tolist())

http_server = WSGIServer(("0.0.0.0", 5500), app)
http_server.serve_forever()
