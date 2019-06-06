sources = $(shell find src -name *.go )

all: publy subly

publy: $(sources)
	go build -o $@ cmd/$@/main.go $(soruces)

subly: $(sources)
	go build -o $@ cmd/$@/main.go $(soruces)

clean:
	rm -rf publy subly
