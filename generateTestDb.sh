curl -X POST http://localhost:8083/insert \
     -H "Content-Type: application/json" \
     -d '{
        "task-identifier": "text-summarization",
        "gpu-memory": 4,
        "gpu-cores": 50,
        "image": "synergcseiitb/bart-large-cnn-text_summarization",
        "startup-latency": 30000,
        "min-latency": 2000,
        "mean-latency": 8000,
        "max-latency": 16000,
        "accuracy": 85,
        "batch-size": 32,
        "end-point": "/summarize",
        "port": 4444,
        "capacity": 2
     }'

curl -X POST http://localhost:8083/insert \
     -H "Content-Type: application/json" \
     -d '{
        "task-identifier": "text-summarization",
        "gpu-memory": 4,
        "gpu-cores": 70,
        "image": "synergcseiitb/bart-large-cnn-text_summarization",
        "startup-latency": 30000,
        "min-latency": 2000,
        "mean-latency": 7000,
        "max-latency": 14000,
        "accuracy": 85,
        "batch-size": 32,
        "end-point": "/summarize",
        "port": 4444,
        "capacity": 3
     }'

curl -X POST http://localhost:8083/insert \
     -H "Content-Type: application/json" \
     -d '{
        "task-identifier": "text-summarization",
        "gpu-memory": 4,
        "gpu-cores": 90,
        "image": "synergcseiitb/bart-large-cnn-text_summarization",
        "startup-latency": 28,
        "min-latency": 1500,
        "mean-latency": 6000,
        "max-latency": 11000,
        "accuracy": 85,
        "batch-size": 32,
        "end-point": "/summarize",
        "port": 4444,
        "capacity": 4
     }'

curl -X POST http://localhost:8083/insert \
     -H "Content-Type: application/json" \
     -d '{
        "task-identifier": "text-summarization",
        "gpu-memory": 4,
        "gpu-cores": 100,
        "image": "synergcseiitb/bart-large-cnn-text_summarization",
        "startup-latency": 28000,
        "min-latency": 15000,
        "mean-latency": 5000,
        "max-latency": 10000,
        "accuracy": 85,
        "batch-size": 32,
        "end-point": "/summarize",
        "port": 4444,
        "capacity": 5
     }'
