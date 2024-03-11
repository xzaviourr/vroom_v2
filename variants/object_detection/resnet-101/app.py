import torch
from transformers import DetrFeatureExtractor, DetrForObjectDetection
from PIL import Image
import requests
from flask import Flask, jsonify, request
import random
import numpy as np
import time
import os
# Initialize the PyTorch device
device = torch.device('cuda') if torch.cuda.is_available() else torch.device('cpu')

# Load the feature extractor and model
feature_extractor = DetrFeatureExtractor.from_pretrained('facebook/detr-resnet-101-dc5')
model = DetrForObjectDetection.from_pretrained('facebook/detr-resnet-101-dc5')
app = Flask(__name__)
# Move the model to the PyTorch device
model.to(device)

# Initialize the Flask app



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


# Define a route that accepts a JSON payload with an image URL
@app.route('/predict', methods=['POST'])
def predict():
    data = request.get_json()
    batch_size = int(1)  # Adjust the batch size according to your GPU memory

    start_time = time.time()
    # Generate 10 thousand images and process them in batches
    for i in range(0, 100, batch_size):
        # Create a batch of image tensors
        images_batch = []
        for j in range(batch_size):
            # Generate or load your images here (replace with your image generation logic)
            image = generate_image()
            images_batch.append(image)

        # Convert the images to PIL Image objects
        image_pil_batch = [Image.fromarray(img) for img in images_batch]

        # Initialize lists to store the outputs of each batch
        all_outputs = []

        # Process the images in batches
        for k in range(0, batch_size, model.config.max_position_embeddings):
            # Extract a sub-batch of images
            image_pil_subbatch = image_pil_batch[k:k+model.config.max_position_embeddings]

            # Extract features from the images
            inputs = feature_extractor(images=image_pil_subbatch, return_tensors="pt")

            # Move the input tensors to the PyTorch device
            inputs = {k: v.to(device) for k, v in inputs.items()}

            # Run the model on the input tensors
            with torch.no_grad():
                outputs = model(**inputs)

            # Append the outputs to the list
            all_outputs.append(outputs)

        # Process the output tensors and send the response
    end_time = time.time()
    elapsed_time = end_time - start_time
    # os.kill(os.getpid(), 9)
    return jsonify({'message': str(elapsed_time)})

# Start the Flask app
if __name__ == '__main__':
    app.run(host='0.0.0.0', port='5123', debug=True)
    
