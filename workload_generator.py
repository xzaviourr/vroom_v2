import asyncio
import aiohttp
import json
from datetime import datetime
import time
from aiohttp import web
import requests

def send_post_request(session):
    url = 'http://localhost:8083/run'        
    payload = {
        "task-identifier": "text-summarization",
        "deadline": 60000.0,
        "accuracy": 80,
        "args": json.dumps({
            "text": "Scientists have discovered a new species of dinosaur in China. The new species belongs to the theropod family, which includes other well-known dinosaurs like the T. rex and velociraptor. The researchers named the new species Haplocheirus sollers, which means simple-handed skillful hunter. The dinosaur lived around 160 million years ago and had long, slender arms and a unique skull.",
        }),
        "response-url": "http://localhost:12367/response",
        "request-size": 1
    }
    try:
        requests.post(url, json=payload)
    except Exception as e:
        print(f"Error while sending POST request: {e}")
        return False, None

def measure_overall_throughput(arrival_rate):
    session = aiohttp.ClientSession()
    
    for _ in range(1):
        send_post_request(session)
        time.sleep(60)
    
    for _ in range(180):
        for _ in range(3):
            send_post_request(session)
        time.sleep(1)

    session.close()

measure_overall_throughput(1)