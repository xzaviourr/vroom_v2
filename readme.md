User gives task. How ?
- Task
- Latency
- Accuracy

Function variant database ?
- Variant id
- task_identifier
- GPU memory
- GPU cores
- image location / name
- latency
- accuracy

Scheduler
- Filter the possible variants
- Look at the available resources on different nodes
- Generate a pod specification file for that particular variant and node
- Launch the pod


CREATE USER 'vroom'@'%' IDENTIFIED WITH mysql_native_password BY 'vroom';
GRANT ALL PRIVILEGES ON *.* TO 'vroom'@'%';