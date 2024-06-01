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

class Pod:
    def __init__(self, memory, compute):
        self.memory = memory
        self.compute = compute

def get_service_ip(service_name, namespace='default'):
    service_ip_cmd = f"kubectl get service {service_name} -n {namespace} -o jsonpath='{{.spec.clusterIP}}'"
    service_ip_output = subprocess.run(service_ip_cmd, shell=True, capture_output=True, text=True)
    service_ip = service_ip_output.stdout.strip()
    return service_ip

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

def send_post_request(service_ip, port):
    api_url = f"http://{service_ip}:{port}/summarize"
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
    
async def send_async_post_request(session, service_ip, port):
    api_url = f"http://{service_ip}:{port}/summarize"  # Assuming the endpoint is /summarize
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

async def measure_overall_throughput(num_requests, pods, ports):
    ips = [get_service_ip(f"{name}-service") for name in pods]
    n = len(pods)
    sessions = [aiohttp.ClientSession() for _ in range(n)]

    start_time = time.time()
    
    tasks = [send_async_post_request(sessions[ind%n], ips[ind%n], ports[ind%n]) for ind in range(num_requests)]
    results = await asyncio.gather(*tasks)

    end_time = time.time()
    
    for session in sessions:
        session.close()

    avg_latency = sum(latency for success, latency in results if success) / num_requests   
    return num_requests / (end_time - start_time) , avg_latency

def create_pod_yaml(mem, com, name, port1, port2):
    yaml_content_pod = {
        "apiVersion": "v1",
        "kind": "Pod",
        "metadata": {"name": name, "labels": {"name": name}},
        "spec": {
            "hostIPC": True,
            "restartPolicy": "OnFailure",
            "securityContext": {"runAsUser": 1000},
            "containers": [
                {
                    "name": "ts1",
                    "image": "synergcseiitb/bart-large-cnn-text_summarization",
                    "imagePullPolicy": "Never",
                    "ports": [{"containerPort": port1}],
                    "resources": {
                        "requests": {"nvidia.com/vcore": com, "nvidia.com/vmem": mem},
                        "limits": {"nvidia.com/vcore": com, "nvidia.com/vmem": mem},
                    },
                }
            ],
        },
    }

    yaml_string = dump(yaml_content_pod)
    # Write YAML content to a file
    with open("pod_request.yaml", "w") as yaml_file:
        yaml_file.write(yaml_string)

    yaml_content_service = {
        "apiVersion": "v1",
        "kind": "Service",
        "metadata": {"name": f"{name}-service"},
        "spec": {
            "selector": {"name": name},
            "ports": [
                {
                    "protocol": "TCP",
                    "port": port2,
                    "targetPort": port1
                }
            ]
        }
    }
    
    yaml_string = dump(yaml_content_service)
    # Write YAML content to a file
    with open("pod_service.yaml", "w") as yaml_file:
        yaml_file.write(yaml_string)

def measure_start_time(start_time, pod_name, port):
    pod_ip = get_service_ip(f"{pod_name}-service")

    # Poll for API readiness
    print("Waiting for the API server to become ready...")
    max_attempts = 300  # Maximum number of attempts
    attempt = 0
    while attempt < max_attempts:
        if send_post_request(pod_ip, port):
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

def run_simulation(pods:List, load:List):
    for pod in pods:
        start_time = time.time()
        
        create_pod_yaml(pod[0].memory, pod[0].compute, "ts1", 4444, 12345)
        subprocess.run(["kubectl", "apply", "-f", "pod_request.yaml"])
        subprocess.run(["kubectl", "apply", "-f", "pod_service.yaml"])
        
        create_pod_yaml(pod[1].memory, pod[1].compute, "ts2", 4444, 12346)
        subprocess.run(["kubectl", "apply", "-f", "pod_request.yaml"])
        subprocess.run(["kubectl", "apply", "-f", "pod_service.yaml"])

        create_pod_yaml(pod[2].memory, pod[2].compute, "ts3", 4444, 12347)
        subprocess.run(["kubectl", "apply", "-f", "pod_request.yaml"])
        subprocess.run(["kubectl", "apply", "-f", "pod_service.yaml"])

        create_pod_yaml(pod[3].memory, pod[3].compute, "ts4", 4444, 12348)
        subprocess.run(["kubectl", "apply", "-f", "pod_request.yaml"])
        subprocess.run(["kubectl", "apply", "-f", "pod_service.yaml"])

        while check_pod_readiness("ts1") == False:
            time.sleep(1)
        
        while check_pod_readiness("ts2") == False:
            time.sleep(1)
        
        while check_pod_readiness("ts3") == False:
            time.sleep(1)

        while check_pod_readiness("ts4") == False:
            time.sleep(1)

        startup_time = measure_start_time(start_time, "ts1", 12345)
        startup_time = measure_start_time(start_time, "ts2", 12346)
        startup_time = measure_start_time(start_time, "ts3", 12347)
        startup_time = measure_start_time(start_time, "ts4", 12348)

        for num_requests in load:
            print(f"Running for - mem1:{pod[0].memory}|mem2:{pod[1].memory}|mem3:{pod[2].memory}|mem4:{pod[3].memory}|com1:{pod[0].compute}|com2:{pod[1].compute}|com3:{pod[2].compute}|com4:{pod[3].compute}|load:{num_requests}\n")
            throughput, latency = asyncio.run(measure_overall_throughput(num_requests, ["ts1", "ts2", "ts3", "ts4"], [12345, 12346, 12347, 12348]))

            output_str = f"{pod[0].memory},{pod[0].compute},{pod[1].memory},{pod[1].compute},{pod[2].memory},{pod[2].compute},{pod[3].memory},{pod[3].compute},{round(startup_time,3)},{num_requests},{round(throughput,3)},{round(latency,3)}\n"
            with open("results.csv", 'a') as file:
                file.write(output_str)

        subprocess.run(["kubectl", "delete", "pod", "ts1"])
        subprocess.run(["kubectl", "delete", "pod", "ts2"])
        subprocess.run(["kubectl", "delete", "pod", "ts3"])
        subprocess.run(["kubectl", "delete", "pod", "ts4"])

if __name__ == "__main__":
    pods = [
        [Pod(2, 10), Pod(2, 10), Pod(2, 10), Pod(2, 10)],   # 8, 40
        [Pod(2, 20), Pod(2, 20), Pod(2, 10), Pod(2, 10)],   # 8, 60
        [Pod(2, 20), Pod(2, 20), Pod(2, 20), Pod(2, 20)],   # 8, 80
        [Pod(2, 40), Pod(2, 20), Pod(2, 20), Pod(2, 20)],   # 8, 100
        [Pod(4, 10), Pod(2, 10), Pod(2, 10), Pod(2, 10)],   # 10, 40
        [Pod(4, 20), Pod(2, 20), Pod(2, 10), Pod(2, 10)],   # 10, 60
        [Pod(4, 20), Pod(2, 20), Pod(2, 20), Pod(2, 20)],   # 10, 80
        [Pod(4, 40), Pod(2, 20), Pod(2, 20), Pod(2, 20)],   # 10, 100
        [Pod(4, 10), Pod(4, 10), Pod(2, 10), Pod(2, 10)],   # 12, 40
        [Pod(4, 20), Pod(4, 20), Pod(2, 10), Pod(2, 10)],   # 12, 60
        [Pod(4, 20), Pod(4, 20), Pod(2, 20), Pod(2, 20)],   # 12, 80
        [Pod(4, 40), Pod(4, 20), Pod(2, 20), Pod(2, 20)],   # 12, 100
        [Pod(4, 10), Pod(4, 10), Pod(4, 10), Pod(2, 10)],   # 14, 40
        [Pod(4, 20), Pod(4, 20), Pod(4, 10), Pod(2, 10)],   # 14, 60
        [Pod(4, 20), Pod(4, 20), Pod(4, 20), Pod(2, 20)],   # 14, 80
        [Pod(4, 40), Pod(4, 20), Pod(4, 20), Pod(2, 20)],   # 14, 100
    ]
    load = [2, 4, 8, 16, 32, 64, 128]

    output_str = f"memory1,compute1,memory2,compute2,memory3,compute3,memory4,compute4,startup_time,load,throughput,latency\n"
    with open("results.csv", 'a') as file:
        file.write(output_str)

    run_simulation(pods, load)