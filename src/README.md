---

# Survey Go Application

This project is a **Survey application** written in Go, using **MongoDB** as the database. It is containerized using Docker and can be run locally either with individual Docker containers or using Docker Compose.

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Prerequisites](#prerequisites)
3. [Project Setup](#project-setup)
4. [Running the Application](#running-the-application)
5. [Running MongoDB Locally](#running-mongodb-locally)
6. [Accessing the Database](#accessing-the-database)
7. [Running with Docker Compose](#running-with-docker-compose)
8. [Testing Connectivity](#testing-connectivity)
9. [Testing the Survey API](#testing-the-survey-api)
10. [Cleaning Up](#cleaning-up)
11. [Important Docker Commands](#important-docker-commands)
12. [Additional Notes](#additional-notes)

---

## 1️⃣ Project Overview

The Survey app allows users to submit answers to surveys, which are stored in a MongoDB database. This README guides you through:

* Setting up the Go application
* Running MongoDB locally
* Connecting the application to the database
* Testing connectivity and API functionality

---

## 2️⃣ Prerequisites

Before running the project, ensure you have the following installed:

* [Go](https://golang.org/doc/install) (v1.21 or higher recommended)
* [Docker](https://www.docker.com/products/docker-desktop)
* [Docker Compose](https://docs.docker.com/compose/install/) (optional)
* `nc` (netcat) for connectivity testing

---

## 3️⃣ Project Setup

Initialize your Go project:

```bash
# Initialize new Go module
go mod init go-app-name

# Download dependencies
go mod tidy

# Verify all modules
go mod verify

# Download modules without updating go.mod
go mod download
```

---

## 4️⃣ Running the Application

You have two options to run your Go application:

```bash
# Option 1: Run directly
go run main.go

# Option 2: Build and run executable
go build -o survey .
./survey
```

> 💡 **Tip:** Use the executable approach for production or when testing outside your development environment.

---

## 5️⃣ Running MongoDB Locally

To start a MongoDB instance in Docker:

```bash
docker run -d \
  --name mongo-local \
  -p 27017:27017 \
  mongo:7
```

### Accessing the Database

```bash
# Enter the MongoDB shell
docker exec -it mongo-local mongosh

# Select your database
use surveyDB

# View all answers
db.answers.find().pretty()
```

> ⚠️ **Warning:** By default, MongoDB runs without authentication. For production, enable access control.

---

## 6️⃣ Running with Docker Compose

Your `docker-compose.yml` should define both the Go application and MongoDB services:

```yaml
services:
  app:
    build: .
    container_name: survey-app
    ports:
      - "8080:8080"
    environment:
      SERVICE_PORT: 8080
      MONGO_URL: mongodb://mongo:27017
    depends_on:
      - mongo
  mongo:
    image: mongo:7
    container_name: mongo-local
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db
volumes:
  mongo_data:
```

Run the application:

```bash
docker-compose up --build
```

> 🔹 This starts both the Go app and MongoDB in one command.

---

## 7️⃣ Testing Connectivity Between App and DB

```bash
# Access the survey app container
docker exec -it survey-app sh

# Test connectivity to MongoDB
nc -zv mongo-local 27017
```

Expected output:

```
mongo-local (172.19.0.2:27017) open
```

> ✅ If successful, your app can reach the database.

---

## 8️⃣ Testing the Survey API

Send a POST request to submit survey answers:

```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
    "answer1": "Go",
    "answer2": "Python",
    "answer3": "JavaScript"
  }'
```

**Response:**

```json
"698b1cdfe04340fabe326928"
```

> The response is a **MongoDB ObjectID**, a unique identifier for this survey submission.

---

### GET the Survey Question

To fetch the survey question (GET request):

```bash
curl http://localhost:8080/
```

**Response:**

```json
"What is your favorite programming language framework?"
```

---

### Verify Answers in MongoDB

Connect to MongoDB to see submitted answers:

```bash
docker exec -it mongo-local mongosh
use surveyDB
db.answers.find().pretty()
```

Sample output:

```json
{
  "_id": ObjectId("698b1cdfe04340fabe326928"),
  "Answer1": "Go",
  "Answer2": "Python",
  "Answer3": "JavaScript"
}
```

---

## 9️⃣ Cleaning Up

To remove unused Docker resources and containers:

```bash
docker system prune -a
```

> ⚠️ This will remove **all stopped containers, networks, and images**. Use with caution.

---

## 🔟 Important Docker Commands

```bash
# Run Docker Compose in the background
docker-compose up -d

# List running Docker Compose services
docker-compose ps

# View logs for services
docker-compose logs -f

# Stop and remove containers
docker-compose down
```

---

##  Additional Notes

* Go version: **1.21**
* MongoDB driver: **mongo-go-driver v1.17.8**
* Database name: `surveyDB`
* Collection name: `answers`
* Go app port: `8080`
* MongoDB port: `27017`