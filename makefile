dev: 
	npx nodemon --signal SIGTERM -e "templ go" -x "templ generate && go run main.go serve" -i "**/*_templ.go"
build:
	templ generate
run:
	go run main.go serve