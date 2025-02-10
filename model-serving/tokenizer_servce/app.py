from flask import Flask, request, jsonify
from transformers import AutoTokenizer

app = Flask(__name__)

# Load the tokenizer
tokenizer = AutoTokenizer.from_pretrained("deepseek-ai/DeepSeek-V3")

@app.route("/tokenize", methods=["POST"])
def tokenize_text():
    data = request.json
    text = data.get("text", "")
    if not text:
        return jsonify({"error": "No text provided"}), 400
    # Tokenize text and convert to IDs
    tokenized = tokenizer.encode(text, add_special_tokens=True)
    return jsonify({"tokens": tokenized})
if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5001)
