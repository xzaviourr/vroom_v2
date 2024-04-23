from flask import Flask, request, jsonify
from transformers import pipeline
import torch

app = Flask(__name__)

# Initialize the summarization pipeline
summarizer = pipeline("summarization", model="facebook/bart-large-cnn", device=0)

if torch.cuda.is_available():
    device = torch.device('cuda:0')
else:
    device = torch.device('cpu')

print("Using device:", device)

@app.route('/summarize', methods=['POST'])
def summarize():
    # Get the text from the JSON body of the request
    text = request.json['text']

    # Generate the summary
    summary = summarizer(text, max_length=130, min_length=30, do_sample=False)

    # Return the summary as a JSON response
    response = {
        'summary': summary[0]['summary_text']
    }
    return jsonify(response)

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=4444)

