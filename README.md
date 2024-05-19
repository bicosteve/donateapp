### Donate App

#### Background 
This application is used to make donations of any kind.
User can create a profile then make donations where they can do all the CRUD
operations on the donations.

#### 1. Prerequisites

To be able to run this application in your local machine, you will need the following;

- Golang
- MySQL

#### 2. Installation

1. Clone the repository to local machine:

- git clone git@github.com:bicosteve/donateapp.git

2. Navigate to the project dir

- cd donateapp

3. Install dependencies

- go mod download
- go mod verify

4. Run the application
   I have provided an environment variables examples on .env.example.
   Use the `sql` dir to find the scripts of created tables for the application.
   After creating the tables and providing the environment vars;
   run **make dev** this will start the server on your localhost

#### 3. Usage

Once the application is running your provided port, register by going to http://localhost:2003/api/v1/users/auth/register

#### 4. Endpoints
These are the end points on the application:

    "/api/v1/users/auth/register" -> register
    "/api/v1/users/auth/login"
    "/server/v1/users/profile"

	"/api/v1/donations/donate"
	"/api/v1/donations/donation/{id}"
	"/api/v1/donations/"
	"/api/v1/donations/donation/{id}"
	"/api/v1/donations/donation/{id}"
    
