generate-example-env:
	@awk -F'=' '{print $$1"="}' .env > example.env