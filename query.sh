curl --location --request POST 'http://10.129.2.22:4444/summarize' \
                --header 'Content-Type: application/json' \
                --data-raw '{
                    "text": "Scientists have discovered a new species of dinosaur in China. The new species belongs to the theropod family, which includes other well-known dinosaurs like the T. rex and velociraptor. The researchers named the new species Haplocheirus sollers, which means \"simple-handed skillful hunter\". The dinosaur lived around 160 million years ago and had long, slender arms and a unique skull."
                }'