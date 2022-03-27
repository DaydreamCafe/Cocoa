ProjectName=Cocoa

build:
	@echo Building project...
	@if [ ! -d "./build" ];then mkdir ./build; else make clean; fi
	@echo Building...
	go build ${ProjectName}.go
	@mv ${ProjectName} ./build/${ProjectName}

clean:
	@if [ -d "./build" ];then echo Cleaning last build...; rm -rf ./build; echo Finished!; else echo No build, passing...; fi

run:
	@echo Running project: ${ProjectName}...
	go run ${ProjectName}.go
