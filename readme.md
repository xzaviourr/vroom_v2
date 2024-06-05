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

CREATE USER 'vroom'@'%' IDENTIFIED VIA mysql_native_password USING '*705F87B9DB50FED7F3353E653181CECF59401A3B';

GRANT ALL PRIVILEGES ON *.* TO 'vroom'@'%';

password hash : *705F87B9DB50FED7F3353E653181CECF59401A3B

curl -X GET 'http://localhost:8083/run?task_id=image-rec&accuracy=50&deadline=2000'

kubectl port-forward service/kube-prometheus-stack-1717-prometheus 9090:9090 -n prometheus
ssh -L 9090:localhost:9090 ub-10@10.129.2.22