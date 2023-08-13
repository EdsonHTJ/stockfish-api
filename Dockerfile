# Start from the latest Go image
FROM golang:latest

# Set the current working directory inside the container
WORKDIR /home/app

# Copy local content to the working directory
COPY . /home/app

# Download all dependencies.
RUN go mod download

# Download and extract Stockfish
RUN wget https://github.com/official-stockfish/Stockfish/releases/download/sf_16/stockfish-ubuntu-x86-64-avx2.tar -O /home/app/stockfish-internal.tar
RUN tar -xvf /home/app/stockfish-internal.tar -C /home/app/

# Build the Go app to a binary named "main"
RUN go build -o main .

# Command to run the application
CMD ["./main"]