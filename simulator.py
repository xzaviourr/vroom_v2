import asyncio
import aiohttp
import json
from datetime import datetime
import time
from aiohttp import web

async def make_post_request(session, url, payload):
    async with session.post(url, json=payload) as response:
        return await response.json()
    
async def process_requests(allRequests):
    tasks = []

    async with aiohttp.ClientSession() as session:
        for request in allRequests:
            url = 'http://localhost:8083/run'
            
            payload = {
                "task-identifier": "text-summarization",
                "deadline": 60000.0,
                "accuracy": 80,
                "args": json.dumps({
                    "text": request
                }),
                "response-url": "http://localhost:12367/response",
                "request-size": 1
            }
            task = asyncio.create_task(make_post_request(session, url, payload))
            tasks.append(task)
        
        responses = await asyncio.gather(*tasks)
        return responses

datasetPath = "dataset/16/30.json"

dataset = []
format_string = '%a %b %d %H:%M:%S %z %Y'

with open(datasetPath) as file:
    for line in file:
        record = json.loads(line.strip())
        try:
            timestamp = record['created_at']
            text = record['text']
            date_object = datetime.strptime(timestamp, format_string)

            dataset.append([date_object, text])
        except:
            pass

# Convert to time difference
initialTs = dataset[0][0]
for ind in range(len(dataset)):
    td = (dataset[ind][0] - initialTs).total_seconds()
    dataset[ind][0] = td

# Second wise dataset
datasetPerSec = [[] for _ in range(60)]
for ind in range(len(dataset)):
    datasetPerSec[int(dataset[ind][0])].append(dataset[ind][1])



async def main():
    # Start the simulation
    for sec in range(60):
        start_time = time.time()
        allRequests = datasetPerSec[sec]
        print("Total requests this second : ", len(datasetPerSec[sec]))

        await process_requests(allRequests[:25])

        end_time = time.time()
        time_taken = (end_time - start_time) * 1000 # milli seconds
        
        remaining_time = 1000 - time_taken
        if remaining_time > 0:
            await asyncio.sleep(remaining_time / 1000.0)
        break

if __name__ == "__main__":
    asyncio.run(main())