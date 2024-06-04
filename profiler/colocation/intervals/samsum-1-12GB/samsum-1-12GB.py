import subprocess
from yaml import dump
import time
import requests
from typing import List
import asyncio
import aiohttp
import logging
import threading

aiohttp_logger = logging.getLogger("aiohttp")
aiohttp_logger.setLevel(logging.ERROR)

def monitor_gpu_utilization(interval, stop_event, results):
    while not stop_event.is_set():
        result = subprocess.run(['nvidia-smi', '--query-gpu=utilization.gpu', '--format=csv,noheader,nounits'], stdout=subprocess.PIPE)
        utilization = int(result.stdout.decode('utf-8').strip())
        results.append((time.time(), utilization))
        time.sleep(interval)

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
    
    tasks = []
    counter = 1

    start_time = time.time()
    for second in range(15):
        for req in range(num_requests):
            tasks.append(send_async_post_request(sessions[counter%n], ips[counter%n], ports[counter%n]))
            counter += 1
        time.sleep(1)
    results = await asyncio.gather(*tasks)
    end_time = time.time()

    for session in sessions:
        session.close()

    latencies = []
    for success, latency in results:    
        if success:
            latencies.append(latency)

    total_time = sum(latency for success, latency in results if success)
    avg_latency = total_time / num_requests   
    return (num_requests*15) / (end_time - start_time) , avg_latency, latencies

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
                    "image": "synergcseiitb/bart-large-cnn-samsum-text_summarization",
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

def run_simulation(pods:List, load:List, num_colocation:int, filename:str):
    for pod in pods:
        # Create YAML content
        start_time = time.time()
        for ind in range(1, num_colocation+1):
            create_pod_yaml(pod[ind-1].memory, pod[ind-1].compute, f"ts{ind}", 5555, 12344 + ind)
            subprocess.run(["kubectl", "apply", "-f", "pod_request.yaml"])
            subprocess.run(["kubectl", "apply", "-f", "pod_service.yaml"])

        for ind in range(1, num_colocation+1):
            while check_pod_readiness(f"ts{ind}") == False:
                time.sleep(1)

        startup_time = 0
        for ind in range(1, num_colocation+1):
            startup_time = measure_start_time(start_time, f"ts{ind}", 12344 + ind)

        for num_requests in load:
            stop_event = threading.Event()
            results = []
            monitoring_thread = threading.Thread(target=monitor_gpu_utilization, args=(0.2, stop_event, results))
            monitoring_thread.start()

            s = "Running for - "
            for ind in range(1, num_colocation+1):
                s = s + f"mem{ind}:{pod[ind-1].memory}|com{ind}:{pod[ind-1].compute}|"
            s += f"Load:{num_requests}/sec\n"
            print(s)

            throughput, latency, latencies = asyncio.run(measure_overall_throughput(num_requests, [f"ts{ind}" for ind in range(1, num_colocation+1)], [12344+ind for ind in range(1, num_colocation+1)]))

            stop_event.set()
            monitoring_thread.join()

            s = ""
            for ind in range(1, num_colocation+1):
                s = s + f"{pod[ind-1].memory},{pod[ind-1].compute},"
            s += f'{round(startup_time, 3)},{num_requests},{round(throughput, 3)},{round(latency, 3)},"{latencies}","{results}"\n'
            with open(filename, 'a') as file:
                file.write(s)

        for ind in range(1, num_colocation+1):
            subprocess.run(["kubectl", "delete", "pod", f"ts{ind}"])

if __name__ == "__main__":
    pods = [
        [Pod(12, 10)],   # 12, 10
        [Pod(12, 20)],   # 12, 20
        [Pod(12, 30)],   # 12, 30
        [Pod(12, 40)],   # 12, 40
        [Pod(12, 50)],   # 12, 50
        [Pod(12, 60)],   # 12, 60
        [Pod(12, 70)],   # 12, 70
        [Pod(12, 80)],   # 12, 80
        [Pod(12, 90)],   # 12, 90
        [Pod(12, 100)],   # 12, 100
    ]
    load = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]
    num_colocation = 1
    file_name = "samsum-1-12GB.csv"

    s = ""
    for ind in range(1, num_colocation+1):
        s += f"memory{ind},compute{ind},"
    s += "startup_time,arrival_rate,throughput,average_latency,latencies,utilization\n"
    with open(file_name, 'a') as file:
        file.write(s)

    run_simulation(pods, load, num_colocation, file_name)