module MasterGobees

go 1.19

require (
	MasterGobees/configuration v0.0.0-00010101000000-000000000000
	MasterGobees/endpoints/home v0.0.0-00010101000000-000000000000
	MasterGobees/endpoints/node v0.0.0-00010101000000-000000000000
	MasterGobees/shell v0.0.0-00010101000000-000000000000
	MasterGobees/utils v0.0.0-00010101000000-000000000000
	github.com/TwiN/go-color v1.2.0
)

require (
	MasterGobees/endpoints/data v0.0.0-00010101000000-000000000000 // indirect
	MasterGobees/endpoints/jobs v0.0.0-00010101000000-000000000000 // indirect
	MasterGobees/globals v0.0.0-00010101000000-000000000000
)

replace MasterGobees/globals => ./globals

replace MasterGobees/configuration => ./configuration

replace MasterGobees/shell => ./shell

replace MasterGobees/utils => ./utils

replace MasterGobees/endpoints/data => ./endpoints/data

replace MasterGobees/endpoints/jobs => ./endpoints/jobs

replace MasterGobees/endpoints/home => ./endpoints/home

replace MasterGobees/endpoints/node => ./endpoints/node
