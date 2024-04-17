
.PHONY: run
run:
	@$(eval FILES_LOCATION=$(shell pwd)/gcs)
	@$(eval GOOGLE_PROJECT_ID=not-a-real-project)
	@export FILES_LOCATION=$(FILES_LOCATION) && export GOOGLE_PROJECT_ID=$(GOOGLE_PROJECT_ID) && go run main.go