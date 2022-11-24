module WorkerGobees

go 1.19

require (
	WorkerGobees/endpoints/data v0.0.0-00010101000000-000000000000
	WorkerGobees/endpoints/home v0.0.0-00010101000000-000000000000
	WorkerGobees/endpoints/jobs v0.0.0-00010101000000-000000000000
	WorkerGobees/endpoints/lifestatus v0.0.0-00010101000000-000000000000
	WorkerGobees/globals v0.0.0-00010101000000-000000000000
	github.com/TwiN/go-color v1.2.0
)

require (
	WorkerGobees/utils v0.0.0-00010101000000-000000000000 // indirect
	github.com/theTardigrade/golang-hash v1.4.2 // indirect
	github.com/twotwotwo/sorts v0.0.0-20160814051341-bf5c1f2b8553 // indirect
)

replace WorkerGobees/globals => ./globals

replace WorkerGobees/utils => ./utils

replace WorkerGobees/endpoints/data => ./endpoints/data

replace WorkerGobees/endpoints/jobs => ./endpoints/jobs

replace WorkerGobees/endpoints/home => ./endpoints/home

replace WorkerGobees/endpoints/lifestatus => ./endpoints/lifestatus
