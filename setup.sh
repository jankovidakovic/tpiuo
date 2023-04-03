alias kafka-topics="docker exec lab1-kafka1-1 kafka-topics --bootstrap-server=kafka1:9092"
alias run_producer="docker exec -it lab1-kafka1-1 kafka-console-producer --bootstrap-server kafka1:9092"
alias run_consumer="docker exec -it lab1-kafka1-1 kafka-console-consumer --bootstrap-server kafka1:9092 --from-beginning"
