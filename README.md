# Pocketbase x HTMX

This is a small demo on how to use [pocketbase](https://pocketbase.io/) as a framework, with [templ](https://templ.guide/) and [htmx](https://htmx.org/).

## Prerequisites
 - templ installed
 - node/npx installed (only for dev)

## Setup
If you have node installed you can start the dev server with `make dev`, else run `make build` and `make run`

After setting up pocketbase import the pocketbase.json file and add a user in the user collection.

Then you can checkout the small demo up under http://localhost:8090/auth/login and login with your newly created user.
