include .env

runserver:
	go run cmd/main.go
	
github-push:
	@echo "pushing..."
	@git push https://${GITHUB_USERNAME}:${GITHUB_PAT}@github.com/${GITHUB_USERNAME}/${GITHUB_REPO_NAME} dev

github-push-front:
	@echo "pushing..."
	@git push https://${GITHUB_USERNAME}:${GITHUB_PAT}@github.com/${GITHUB_USERNAME}/${GITHUB_REPO_NAME} impl-frontend
