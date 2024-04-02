Variant-
    * Object Detection 
    * Resnet 101
    * Pytorch

List of commands-
    * Build - sudo docker build -t gpu-object-detection .
    * Run - sudo docker run -d --gpus all -p 5123:5123 --runtime=nvidia gpu-object-detection
    * Test - 
        Input - curl -X POST -H "Content-Type: application/json" -d '{"image_url": "http://images.cocodataset.org/val2017/000000039769.jpg"}' http://localhost:5123/predict
        Ouput - toothbrush

List of items in dockerfile that can be changed-
    * Base image should be nvidia/cuda:12.1.0-runtime-ubuntu20.04, on 18.04 there were many version errors due to python3.6 and pip9
    * timm package needs to be installed(pipreqs cannot capture this package from app.py)
    * pipreqs --mode no-pin does not list the versions along with the package
                    (but on ubuntu20.04 should work on default mode - NOT TESTED)


List of commands (New)-
    * Build - sudo docker build -t synergcseiitb/object-detection-resnet:v1 .
    * Run - sudo docker run -d --gpus all -p 5123:5123 --runtime=nvidia synergcseiitb/object-detection-resnet:v1
    * Test - 
        Input - curl -X POST -H "Content-Type: application/json" -d '{"image_url": "http://images.cocodataset.org/val2017/000000039769.jpg"}' http://localhost:5123/predict
        Ouput - toothbrush