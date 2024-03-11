from transformers import YolosFeatureExtractor, YolosForObjectDetection
from PIL import Image
import requests
import torch
from flask import Flask, request, jsonify
import time
import numpy as np
import random

app = Flask(__name__)

feature_extractor = YolosFeatureExtractor.from_pretrained('hustvl/yolos-tiny')
model = YolosForObjectDetection.from_pretrained('hustvl/yolos-tiny')
device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
print(f"Using device: {device}")
model.to(device)

def generate_image():
    # Generate a random image
    width = 224
    height = 224
    channels = 3

    image = np.zeros((height, width, channels), dtype=np.uint8)

    for h in range(height):
        for w in range(width):
            for c in range(channels):
                # Generate random pixel values between 0 and 255
                pixel_value = random.randint(0, 255)
                image[h, w, c] = pixel_value

    return image


@app.route('/predict', methods=['POST'])
def predict():
    batch_size = int(request.json['batch_size'])
    start_time = time.time()
   # Generate 100 images
    images = []
    for _ in range(100):
        image = generate_image()
        images.append(image)
    
    # Convert the images to PIL Image objects
    image_pil_batch = [Image.fromarray(img) for img in images]

    # Initialize lists to store the outputs of each batch
    all_outputs = []

    # Process the images in batches
    for i in range(0, 100, batch_size):
        # Extract a sub-batch of images
        image_pil_subbatch = image_pil_batch[i:i+batch_size]

        # Extract features from the images
        inputs = feature_extractor(images=image_pil_subbatch, return_tensors="pt")
        inputs = inputs.to(device)

        # Run the model on the input tensors
        with torch.no_grad():
            outputs = model(**inputs)

        # Append the outputs to the list
        all_outputs.append(outputs)
    end_time = time.time()
    elapsed_time = end_time - start_time
    return jsonify({'message': str(elapsed_time)})


if __name__ == '__main__':
    app.run(host='0.0.0.0', port='5126', debug=True)

