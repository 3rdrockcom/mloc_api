# The name of the executable
TARGET=mloc
.DEFAULT_GOAL: $(TARGET)

# Command arguments and flags
OPTIONS=""

# These are the values we want to pass for VERSION and BUILD
VERSION=1.0.2-dev
BUILD=`git rev-parse HEAD`

# Setup the -ldflags option for go build here, interpolate the variable values
LD_FLAGS=-ldflags="-X github.com/epointpayment/mloc_api_go/app/config.Version=$(VERSION) -X github.com/epointpayment/mloc_api_go/app/config.Build=$(BUILD)"

# Ignore phony targets
.PHONY: build install clean deps vendor run run-development run-production run-watch

# Builds project
$(TARGET):
	go build $(LD_FLAGS) -o $(TARGET)

build: $(TARGET)
	@true

# Installs project: copies binary
install:
	go install $(LD_FLAGS) -o $(TARGET)

# Cleans project: deletes binary
clean:
	if [ -f $(TARGET) ] ; then rm $(TARGET) ; fi

# Get project dependencies
deps:
	go get
	go get github.com/cespare/reflex
	go get -u github.com/golang/dep/cmd/dep

# Vendor project dependencies
vendor:
	dep ensure

# Runs project: executes binary
run:
	./$(TARGET) ${OPTIONS}

# Runs project: executes binary with development settings
run-development:
	godotenv -f .env.development ./$(TARGET) ${OPTIONS}

# Runs project: executes binary with production settings
run-production:
	godotenv -f .env.production ./$(TARGET) ${OPTIONS}

# Runs project: watches directory for changes and executes binary with development settings
run-watch:
	reflex -d 'none' -R 'vendor.' -r '\.go$\' -s -- sh -c 'make clean && make build && make run-development OPTIONS=serve'
