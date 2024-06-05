# import time
# import threading
# from prometheus_api_client import PrometheusConnect

# def fetch_gpu_utilization():
#     prometheus_url = "http://localhost:9090"  # URL of your Prometheus server
#     prometheus = PrometheusConnect(url=prometheus_url)

#     query = 'sum(rate(DCGM_FI_DEV_GPU_UTIL{job="dcgm-exporter"}[5m]))'  # Example query for GPU utilization

#     while True:
#         result = prometheus.custom_query(query)
#         print("GPU Utilization:", result)
#         time.sleep(1)  # Fetch result every second

# def main():
#     # Create a separate thread for fetching GPU utilization
#     gpu_thread = threading.Thread(target=fetch_gpu_utilization)
#     gpu_thread.daemon = True  # Daemonize the thread so it terminates when the main thread terminates
#     gpu_thread.start()

#     # Keep the main thread running
#     while True:
#         time.sleep(1)

# if __name__ == "__main__":
#     main()

import time
import threading
from datetime import datetime, timedelta
from prometheus_api_client import PrometheusConnect

def fetch_gpu_utilization():
    prometheus_url = "http://localhost:9090"  # URL of your Prometheus server
    prometheus = PrometheusConnect(url=prometheus_url)

    query = 'DCGM_FI_DEV_GPU_UTIL'  # Metric query for GPU utilization

    end_time = datetime.now()  # End time (current time)
    start_time = end_time - timedelta(seconds=60)  # Start time (1 hour ago)
    step = 1  # Step interval in seconds

    while True:
        try:
            # Execute range query to fetch GPU utilization data
            result = prometheus.custom_query_range(query=query, start_time=start_time, end_time=end_time, step=step)
            print("Query Result:", result[0]['values'])  # Print query result for debugging
            if result and 'result' in result and 'values' in result['result']:
                # Process the query result
                for value in result['result']['values']:
                    timestamp = value[0]
                    gpu_utilization = float(value[1])
                    print("Timestamp:", timestamp, "GPU Utilization:", gpu_utilization)

        except Exception as e:
            print("Error fetching GPU utilization:", e)

        time.sleep(step)  # Wait for the next query interval

def main():
    # Create a separate thread for fetching GPU utilization
    gpu_thread = threading.Thread(target=fetch_gpu_utilization)
    gpu_thread.daemon = True  # Daemonize the thread so it terminates when the main thread terminates
    gpu_thread.start()

    # Keep the main thread running
    while True:
        time.sleep(1)

if __name__ == "__main__":
    main()
