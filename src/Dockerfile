FROM golang:1.20.3

# Set the working directory inside the container
WORKDIR /app

## Copy the Go module and sum files to the container
#COPY go.mod ./



# Copy the source code to the container
COPY . /app

#RUN source .bashrc
# Download dependencies
RUN go mod tidy

RUN cd /app/ups/cmd

RUN go build -o main .
## Expose the port that the application will run on
#EXPOSE 8080

# Run the compiled binary
CMD ["go run main.go"]



