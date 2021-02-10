run:
	go run *.go

gomod-exp:
	export GO111MODULE=on

gobuild:
	GOOS=linux GOARCH=amd64 go build -o update
dockerbuild:
	docker build -t update .
dockerbuildandpush:
	docker build -t update .
	docker tag update americanwonton/update
	docker push americanwonton/update
dockerrun:
	docker run -it -p 80:80 update
dockerrundetached:
	docker run -d -p 80:80 update
dockerrunitvolume:
	docker run -it -p 80:80 -v photo-images:/static/images update
dockerrundetvolume:
	docker run -d -p 80:80 -v photo-images:/static/images update
dockertagimage:
	docker tag update americanwonton/update
dockerimagepush:
	docker push americanwonton/update
dockerallpush:
	docker tag update americanwonton/update
	docker push americanwonton/update
dockerseeshell:
	docker run -it update sh