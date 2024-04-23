curl -X POST http://localhost:8083/run \
     -H "Content-Type: application/json" \
     -d '{
        "task-identifier": "text-summarization",
        "deadline": 60000.0,
        "accuracy": 80,
        "args": "{\"text\": \"Scientists have discovered a new species of dinosaur in China. The new species belongs to the theropod family, which includes other well-known dinosaurs like the T. rex and velociraptor. The researchers named the new species Haplocheirus sollers, which means simple-handed skillful hunter. The dinosaur lived around 160 million years ago and had long, slender arms and a unique skull.\"}",
        "response-url": "http://example.com/response",
        "request-size": 1
     }'