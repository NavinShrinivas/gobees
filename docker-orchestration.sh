docker rm -f $(docker ps -a -q)
docker build -t master_image ./MasterNodeServices/
docker build -t worker_image ./WorkerNodeServices/
docker compose up
