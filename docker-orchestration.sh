sudo docker rm -f $(sudo docker ps -a -q)
sudo docker build -t master_image ./MasterNodeServices/ 
sudo docker build -t worker_image ./WorkerNodeServices/ 
sudo docker compose up
