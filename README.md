# soft-exam-2

## How to run the project

1. Clone the repository `gh repo clone rasm445f/soft-exam-2`
2. Run in your terminal `chmod +x install-tools.sh`
3. Run the install-tools.sh scripts like so `./install-tools.sh` if you don't already have the tools. (Ensure your PATH is set up correctly).
4. Rename the `.example.env` file to `.env` and populate the fields as needed.
5. Make sure you are in the ROOT of the project, Run `docker compose up` to start the PostgreSQL Docker container or `docker compose up -d` to run it in the background as a daemon on every boot.
6. Make a new terminal in the ROOT of the project
7. Run `make migrate-up` to setup the database with the tables etc. specified in the `db/migration/` folder.
8. Run `make run` to start the server.
9. Check the server is running by visiting `http://localhost:8080/` in your browser.
10. you can now test the endpoints using the swagger documentation at `http://localhost:8080/swagger/index.html`

## Technology Stack

### Version Control Platform:

- Git - Github

### Text Editing and Development Environment:

- VSCode / Neovim
- DBeaver
- Swagger

### General Online Research Tools:

- Stack Overflow
- MDN Web Docs
- Golang Docs

## Development Stack

### Backend Development:

- Golang

### Database Management:

- PostgreSQL

### Development Tools:

- Docker
- Docker Compose

### CI/CD Pipeline:

- GitHub Actions
