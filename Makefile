# The name of the executable
TARGET=mloc
.DEFAULT_GOAL: $(TARGET)

# Command arguments and flags
OPTIONS=""

# These are the values we want to pass for VERSION and BUILD
VERSION=0.0.1
BUILD=`git rev-parse HEAD`

# Setup the -ldflags option for go build here, interpolate the variable values
LD_FLAGS=-ldflags="-X github.com/epointpayment/mloc_api_go/app/config.Version=$(VERSION) -X github.com/epointpayment/mloc_api_go/app/config.Build=$(BUILD)"

# Ignore phony targets
.PHONY: build clean install run run-development run-production

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

# Runs project: executes binary
run:
	./$(TARGET) ${OPTIONS}

# Runs project: executes binary with development settings
run-development:
	godotenv -f .env.development ./$(TARGET) ${OPTIONS}

# Runs project: executes binary with production settings
run-production:
	godotenv -f .env.production ./$(TARGET) ${OPTIONS}




