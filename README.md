# GoBees

Anywhere you see `SS` in this project, it's the short form for Shared Storage.

## Preface
GoBees is yet another Map Reduce framework, written and powered by GO! This project was never built or meant for production cases, but given how well we have written it, it might just pass off in prod :)). During the development of this project few assumptions were made, these assumptions were solely to meet the submission deadlines. Some of these are : 
```
> Output from mapper file HASSSSS to be <key,value> => comma is a must!

> NOTE : Giving custom partition function will RADICALLY slow down map reduce, this is due to the limitations of golang not havin generic implementations at runtime like Rust or Java :,(

> Note : if the  custom partition function is not following the fixed template, it may lead to infifite job.
```

## Features and Behind the hood : 

- Powered by go
- Written entirely from scratch, Uses almost 0 external libraries. Only uses a fast string hasher and hyper-fast in memory quick sort.
- Mapper and Reducer files for now have to be provided in python, but extending it to support other languages should be less than 10 lines of code per language!
- Supports custom partition/shuffle function written in go. Note: It has to follow the temple present [here](./MasterNodeServices/test/shuffle_streamer.go)
- Can support local and network workers at the same time!
- We have an entirely dockerized setup, it can even run on windows!!

## Getting started | Usage

- Clone repo and move into project root : 
```
git clone https://github.com/NavinShrinivas/gobees ~/gobees && cd ~/gobees
```

### Manual setup

- First, start the master node :
```
cd MasterNodeServices
go run .
```
- Now open a new terminal and start 1 worker node : 
```
cd WorkerNodeService
go run .
# Start worker in 5000
```

- To add more workers : 
> Note : Each worker node runs in its own terminal
```
# If in the same machine (i.e localhost) :
cd WorkerNodeService
go run . -port=5001
```

```
# If not in localhost 
cd WorkerNodeService 
go run . -master="http://ip.of.master.node:3000" 
```

### Docker setup


