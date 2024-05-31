from flask import Flask, request, jsonify
from transformers import pipeline
import torch
import threading
import time
from queue import Queue

app = Flask(__name__)

# Initialize the summarization pipeline
device = 0 if torch.cuda.is_available() else -1
model = "facebook/bart-large-cnn"
summarizer = pipeline("summarization", model=model, device=device)

# Queue to store incoming texts
text_queue = Queue()
# Dictionary to store results
results = {}
# Lock for thread-safe operations
lock = threading.Lock()

def process_queue():
    while True:
        batch = []
        batch_ids = []
        
        while not text_queue.empty() and len(batch) <= 32:  # Adjust batch size as needed
            request_id, text = text_queue.get()
            batch.append(text)
            batch_ids.append(request_id)
        
        if batch:
            summaries = summarizer(batch, max_length=130, min_length=30, do_sample=False)
            with lock:
                for i, summary in enumerate(summaries):
                    results[batch_ids[i]] = summary['summary_text']
        
        time.sleep(1)  # Adjust the sleep interval as needed

@app.route('/summarize', methods=['POST'])
def summarize():
    text = request.json['text']
    request_id = threading.get_ident()
    
    with lock:
        results[request_id] = None  # Initialize the result for this request
    
    text_queue.put((request_id, text))
    
    while True:
        with lock:
            if results[request_id] is not None:
                summary = results.pop(request_id)
                return jsonify({"summary": summary})
        
        time.sleep(0.1)  # Adjust the sleep interval as needed

if __name__ == '__main__':
    # Start the batch processing thread
    batch_thread = threading.Thread(target=process_queue, daemon=True)
    batch_thread.start()
    
    app.run(debug=True, host='0.0.0.0', port=4444, threaded=True)