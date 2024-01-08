clean:
	-sudo docker compose down

dev:
	-sudo docker compose down
	npx tailwindcss -i ./src/public/dist/input.css -o ./src/public/dist/output.css
	docker compose up --build -d
	#docker exec -it simplewebsite bash

tailwind:
	npx tailwindcss -i ./src/public/dist/input.css -o ./src/public/dist/output.css --watch

prod:
	npx tailwindcss -i ./src/public/dist/input.css -o ./src/public/dist/output.css
	docker compose up --build