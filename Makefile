include ./.env

web-gen:
	templ generate
	tailwindcss -i ./assets/input.css -o ./assets/tailwind.css

web-gen-run : web-gen
	go run ./cmd/werbsrv/main.go

run-web: 
	go run ./cmd/werbsrv/main.go

air : web-gen
	air
