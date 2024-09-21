# MiniTwitter - A Basic Twitter-Like Application

## Project Overview
MiniTwitter is a lightweight, Twitter-like application that allows users to post tweets, follow other users, like/retweet tweets, and perform basic CRUD operations on user profiles and tweets. It focuses on functionality and scalability, leveraging modern backend technologies like Golang, PostgreSQL, MongoDB, and message brokers for efficient processing. 

## Features

### 1. User Authentication
- Sign up, log in, and log out functionality.
- Basic profile management (name, username, bio, profile picture).
  
### 2. Posting Tweets
- Authenticated users can post text-based tweets with optional images or videos.
  
### 3. Following/Unfollowing Users
- Users can follow or unfollow other users.
- Users can view their following list and follower list.
  
### 4. Likes and Retweets
- Users can like or unlike tweets.
- Retweets appear on the retweeter's followers' timelines.

### 5. Timeline and Search
- Users can view a timeline of tweets from the people they follow.
- Basic search functionality to find users or tweets.

## Technical Details

### 1. Backend
- **Language**: Golang
- **Framework**: Gin for HTTP routing and gRPC for inter-service communication.
- **Database**: PostgreSQL for user data and MongoDB for storing tweets, likes, and comments.
  
### 2. Message Broker
- **Kafka** (or RabbitMQ) for handling asynchronous tasks like notifications, timeline updates, etc.

### 3. Authentication
- **JWT** token-based authentication for secure API communication.
- **bcrypt** for password hashing to enhance security.

### 4. Caching
- **Redis** for caching frequently accessed data like timelines and user profiles to optimize performance.
  
### 5. Storage
- **Amazon S3** or **Google Cloud Storage** for media file storage (profile pictures, tweet images, videos).

### 6. API Documentation
- **Swagger** for documenting API endpoints, making it easier for developers to interact with the app's RESTful services.

### 7. DevOps and Deployment
- **Docker**: The entire application is containerized using Docker for easy setup and deployment.
- **CI/CD**: Continuous integration and deployment using **GitHub Actions**.
- **Cloud**: Deployable to cloud providers like AWS or Google Cloud using **Terraform** for infrastructure automation.

### 8. Security
- Implemented **role-based access control (RBAC)** for permissions and access management.
- **Input validation and sanitization** to prevent SQL injection and XSS attacks.

### 9. Performance
- Database queries optimized with indexes.
- **Load testing** using **k6** to ensure the app handles up to 100k RPS on GET endpoints and 5k TPS for POST/PUT operations.
  
## Project Structure


MiniTwitter/
│
├── cmd/
│   └── main.go           # Entry point for the application
│
├── api/
│   ├── user.go           # Handlers for user authentication, profile
│   └── tweet.go          # Handlers for tweets, likes, retweets
│
├── internal/
│   ├── db/               # Database migrations and models
│   └── broker/           # Kafka/RabbitMQ implementation
│
├── config/
│   └── config.go         # Configuration handling (env, .yaml)
│
├── docs/
│   └── swagger.yaml      # API documentation using Swagger
│
├── tests/                # Unit and integration tests
└── README.md             # Project setup, instructions


## Setup Instructions

1. **Clone the Repository**:
   
  git clone https://github.com/your-repo/MiniTwitter.git
   

2. **Environment Variables**: Set up the `.env` file with required environment variables like `DB_URL`, `JWT_SECRET_KEY`, `REDIS_URL`, `KAFKA_BROKER`.

3. **Docker**: Start all services using Docker Compose:
   
  docker-compose up --build
   

4. **Database Migration**: Run migrations using golang-migrate:
   
  migrate -path db/migrations -database "postgres://user:password@localhost:5432/minitwitter?sslmode=disable" up
   

5. **Running the Application**: Start the Go server:
   
  go run cmd/main.go
   

6. **Access the API**: Open `http://localhost:8080/swagger/index.html` to explore API endpoints.

## Challenges and Future Enhancements

- **Scalability**: Implementing horizontal scaling and database sharding to handle a high volume of users and tweets.
- **Real-time Features**: Integrating WebSockets for real-time notifications and updates.
- **Advanced Analytics**: Adding data pipeline support for analyzing trending tweets and user activity using Apache Spark.