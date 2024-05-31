# Use the NVIDIA CUDA runtime as the base image
FROM nvidia/cuda:12.1.0-runtime-ubuntu20.04
# Set the working directory to /app
# WORKDIR /app
RUN mkdir app
ENV XDG_CACHE_HOME /app/.cache

# Copy the Flask app files to the container

ADD app.py app
# COPY requirements.txt .

# Create cache directory and set permissions
RUN mkdir -p $XDG_CACHE_HOME/hub && \
    chmod -R 777 $XDG_CACHE_HOME

# Install dependencies
RUN apt-get update
RUN apt-get install -y python3 python3-pip
RUN pip3 install timm
RUN pip3 install pipreqs
RUN pipreqs --force --mode no-pin app
RUN pip3 install -r app/requirements.txt


# Expose port 5555 for the Flask app
EXPOSE 5555

# Start the Flask app when the container starts
CMD ["python3", "app/app.py"]
