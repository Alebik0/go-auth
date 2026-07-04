# User authentification on Go

Testing web site written on Go to test authentification.

## Architecture

Just a simple microservice architecture, that can be seen on that image:

![Project architecture](./diagrams/Services%20Diagram.png)

### Frontend

Just a simple web application

#### Tasks:

- [x] Learn base React
- [x] Learn base NextJS
- [x] Learn base Axios library
- [ ] Run frontend via NextJS + Vite

### API Gateway

Authentification and inner API expose

#### Tasks:

- [x] Test Redis
- [x] Test user authorization via JWT tokens
- [x] CORS support
- [ ] CSRF protection

### User service

User API to manage user profiles. Interacts with the database.

#### Tasks:

- [x] Test Gin framework
- [x] Test Swagger
- [x] Test Go unit testing

### Database

Just a simple PostgreSQL database.

#### Tasks:

- [x] Run a simple database
