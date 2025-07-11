include .env

run:
	go run cmd/main.go
	
github-push:
	@echo "pushing..."
	@git push https://nikolat27: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@github.com/nikolat27/file_manager dev
