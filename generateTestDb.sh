curl -X POST http://localhost:8083/insert \
     -H "Content-Type: application/json" \
     -d '{
        "task-identifier": "text-summarization",
        "gpu-memory": 4,
        "gpu-cores": 50,
        "image": "synergcseiitb/bart-large-cnn-text_summarization",
        "startup-latency": 30,
        "min-latency": 2,
        "mean-latency": 8,
        "max-latency": 16,
        "accuracy": 85,
        "batch-size": 32,
        "end-point": "/summarize",
        "port": 4444,
        "capacity": 30
     }'

curl -X POST http://localhost:8083/insert \
     -H "Content-Type: application/json" \
     -d '{
        "task-identifier": "text-summarization",
        "gpu-memory": 4,
        "gpu-cores": 70,
        "image": "synergcseiitb/bart-large-cnn-text_summarization",
        "startup-latency": 30,
        "min-latency": 2,
        "mean-latency": 7,
        "max-latency": 14,
        "accuracy": 85,
        "batch-size": 32,
        "end-point": "/summarize",
        "port": 4444,
        "capacity": 35
     }'

curl -X POST http://localhost:8083/insert \
     -H "Content-Type: application/json" \
     -d '{
        "task-identifier": "text-summarization",
        "gpu-memory": 4,
        "gpu-cores": 90,
        "image": "synergcseiitb/bart-large-cnn-text_summarization",
        "startup-latency": 28,
        "min-latency": 1.5,
        "mean-latency": 6,
        "max-latency": 11,
        "accuracy": 85,
        "batch-size": 32,
        "end-point": "/summarize",
        "port": 4444,
        "capacity": 45
     }'

curl -X POST http://localhost:8083/insert \
     -H "Content-Type: application/json" \
     -d '{
        "task-identifier": "text-summarization",
        "gpu-memory": 4,
        "gpu-cores": 100,
        "image": "synergcseiitb/bart-large-cnn-text_summarization",
        "startup-latency": 28,
        "min-latency": 1.5,
        "mean-latency": 5,
        "max-latency": 10,
        "accuracy": 85,
        "batch-size": 32,
        "end-point": "/summarize",
        "port": 4444,
        "capacity": 50
     }'
