# User authentification on Go

Testing web site written on Go to test authentification. 

## Architecture

Just a simple microservice architecture, that can be seen on that image:

![Project architecture](./images/architecture.png)

### Frontend 

Just a simple web application

### API Gateway

Authentification and inner API expose

### User service

User API to manage user profiles. Interacts with the database.

### Database

Just a simple PostgreSQL database.
