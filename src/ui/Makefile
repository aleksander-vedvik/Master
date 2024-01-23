clean:
	-sudo docker compose down

tailwind:
	cd frontend; npx tailwindcss -i ./src/public/dist/input.css -o ./src/public/dist/output.css

ui: tailwind
	docker compose up -d

ui-rebuild: tailwind
	docker compose up --build -d