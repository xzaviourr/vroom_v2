import subprocess
from yaml import dump
import time
import requests
from typing import List
import asyncio
import aiohttp
import logging

aiohttp_logger = logging.getLogger("aiohttp")
aiohttp_logger.setLevel(logging.ERROR)
logging.basicConfig(filename='debug.log', level=logging.ERROR)

def get_pod_ip(pod_name):
    pod_ip_cmd = f"kubectl get pod {pod_name} -o jsonpath='{{.status.podIP}}'"
    pod_ip_output = subprocess.run(pod_ip_cmd, shell=True, capture_output=True, text=True)
    pod_ip = pod_ip_output.stdout.strip()
    return pod_ip

def check_pod_readiness(pod_name):
    readiness_cmd = f"kubectl get pod {pod_name} -o jsonpath='{{.status.conditions[?(@.type==\"Ready\")].status}}'"
    readiness_output = subprocess.run(readiness_cmd, shell=True, capture_output=True, text=True)
    readiness_str = readiness_output.stdout.strip()
    return readiness_str == "True"

def send_post_request(pod_ip):
    api_url = f"http://{pod_ip}:4444/summarize"
    headers = {"Content-Type": "application/json"}
    data = {
        "text": "Scientists have discovered a new species of dinosaur in China. The new species belongs to the theropod family, which includes other well-known dinosaurs like the T. rex and velociraptor. The researchers named the new species Haplocheirus sollers, which means 'simple-handed skillful hunter'. The dinosaur lived around 160 million years ago and had long, slender arms and a unique skull."
    }
    try:
        response = requests.post(api_url, headers=headers, json=data)
        return response.status_code == 200
    except Exception as e:
        print(f"Error while sending POST request: {e}")
        return False
    
async def send_async_post_request(session, pod_ip):
    api_url = f"http://{pod_ip}:4444/summarize"  # Assuming the endpoint is /summarize
    headers = {"Content-Type": "application/json"}
    data = {
        "text": "Scientists have discovered a new species of dinosaur in China. The new species belongs to the theropod family, which includes other well-known dinosaurs like the T. rex and velociraptor. The researchers named the new species Haplocheirus sollers, which means 'simple-handed skillful hunter'. The dinosaur lived around 160 million years ago and had long, slender arms and a unique skull."
    }
    try:
        start_time = time.time()
        async with session.post(api_url, headers=headers, json=data) as response:
            end_time = time.time()
            latency = (end_time - start_time) * 1000  # Convert to milliseconds
            result = await response.text()
            return response.status == 200, latency
    except Exception as e:
        print(f"Error while sending POST request: {e}")
        return False, None

def measure_start_time(start_time):
    pod_name = "ts"
    pod_ip = get_pod_ip(pod_name)

    # Poll for API readiness
    print("Waiting for the API server to become ready...")
    max_attempts = 120  # Maximum number of attempts
    attempt = 0
    while attempt < max_attempts:
        if send_post_request(pod_ip):
            print("API server is ready.")
            break
        else:
            attempt += 1
            time.sleep(1)  # Wait for 1 seconds before retrying
    else:
        print("Timeout: API server did not become ready within the specified time.")

    # Calculate startup time (if API server became ready)
    if attempt < max_attempts:
        end_time = time.time()
        startup_time_ms = (end_time - start_time) * 1000
        return startup_time_ms
    return -1

async def measure_throughput(num_requests):
    pod_name = "ts"
    pod_ip = get_pod_ip(pod_name)

    start_time = time.time()
    async with aiohttp.ClientSession() as session:
        tasks = [send_async_post_request(session, pod_ip) for _ in range(num_requests)]
        results = await asyncio.gather(*tasks)
    end_time = time.time()

    avg_latency = sum(latency for success, latency in results if success) / num_requests   
    return num_requests / (end_time - start_time) , avg_latency

def run_simulation(memory_values:List, compute_values:List, load:List):
    for mem in memory_values:
        for com in compute_values:
            # Create YAML content
            yaml_content = {
                "apiVersion": "v1",
                "kind": "Pod",
                "metadata": {"name": "ts", "labels": {"name": "ts"}},
                "spec": {
                    "hostIPC": True,
                    "restartPolicy": "OnFailure",
                    "securityContext": {"runAsUser": 1000},
                    "containers": [
                        {
                            "name": "ts",
                            "image": "synergcseiitb/bart-large-cnn-text_summarization",
                            "imagePullPolicy": "Never",
                            "ports": [{"containerPort": 4444}],
                            "resources": {
                                "requests": {"nvidia.com/vcore": com, "nvidia.com/vmem": mem},
                                "limits": {"nvidia.com/vcore": com, "nvidia.com/vmem": mem},
                            },
                        }
                    ],
                },
            }

            # Convert YAML content to string
            yaml_string = dump(yaml_content)

            # Write YAML content to a file
            with open("pod_request.yaml", "w") as yaml_file:
                yaml_file.write(yaml_string)

            start_time = time.time()
            # Apply YAML file using kubectl
            subprocess.run(["kubectl", "apply", "-f", "pod_request.yaml"])
            pod_name = "ts"
            
            while not check_pod_readiness(pod_name):
                time.sleep(1)

            startup_time = measure_start_time(start_time)

            for num_requests in load:
                print(f"Running for - mem:{mem}|com:{com}|load:{num_requests}\n")
                throughput, latency = asyncio.run(measure_throughput(num_requests))

                output_str = f"{mem},{com},{round(startup_time,3)},{num_requests},{round(throughput,3)},{round(latency,3)}\n"
                with open("results.csv", 'a') as file:
                    file.write(output_str)

            subprocess.run(["kubectl", "delete", "pod", "ts"])

if __name__ == "__main__":
    memory_values = [2, 4, 6, 8, 10, 12, 14]
    compute_values = [20, 40, 60, 80, 100]
    load = [1, 5, 20]

    output_str = f"memory,compute,startup_time,load,throughput,latency\n"
    with open("results.csv", 'a') as file:
        file.write(output_str)

    run_simulation(memory_values, compute_values, load)


