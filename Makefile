.PHONY: run
run:
	@echo "Running the program..."
	go run .

.PHONY: tidy
tidy:
	@echo "Tidying up the go modules..."
	go mod tidy