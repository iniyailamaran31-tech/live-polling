# Live Polling Tool

A real-time web-based polling application where users can create polls, share them with others, vote, and view live results without refreshing the page.

## 🚀 Live Demo

Frontend: https://live-polling-wheat-theta.vercel.app/

Backend: https://live-polling-oh5m.onrender.com/

## ✨ Features

- User registration and login
- JWT-based authentication
- Create polls with multiple options
- Share polls using a public link
- Real-time voting results
- Results update automatically without page refresh
- WebSocket-based live updates
- Redis-powered real-time vote counting
- MongoDB for persistent application data
- Backend validation for poll creation and voting
- Responsive and modern user interface

## 🛠️ Technology Stack

### Frontend
- React.js
- Vite
- JavaScript
- CSS

### Backend
- Go
- Gin Framework
- JWT Authentication
- WebSocket

### Database
- MongoDB

### Real-Time Layer
- Redis
- Redis Pub/Sub
- WebSocket

### Deployment
- Vercel – Frontend
- Render – Backend
- MongoDB Atlas – Database
- Upstash Redis – Redis

## 🏗️ Architecture

```text
                    ┌─────────────────────┐
                    │      React UI       │
                    │   Vite + JavaScript │
                    └──────────┬──────────┘
                               │
                         HTTP / WebSocket
                               │
                               ▼
                    ┌─────────────────────┐
                    │    Go + Gin API     │
                    │ Authentication      │
                    │ Poll Management     │
                    │ Vote Validation     │
                    └──────┬─────────┬────┘
                           │         │
                           ▼         ▼
                  ┌────────────┐  ┌────────────┐
                  │  MongoDB   │  │   Redis    │
                  │ Persistent │  │ Live Votes │
                  │   Data     │  │ Pub/Sub    │
                  └────────────┘  └──────┬─────┘
                                          │
                                     WebSocket
                                          │
                                          ▼
                                  Connected Clients
                                  Live Results
🔄 How It Works
1. User Authentication

A user registers or logs in through the React frontend.

The Go backend validates the request and uses JWT authentication to protect poll creation and management operations.

2. Poll Creation

An authenticated user creates a poll by providing:

Poll question
Two to six answer options

The Go backend validates the data and stores the poll in MongoDB.

3. Sharing

After creating a poll, the application provides a public poll link.

Anyone with the link can open the poll and vote.

4. Voting

When a user selects an option:

React
  ↓
Go/Gin API
  ↓
Validate poll + option
  ↓
Redis HINCRBY
  ↓
Redis Pub/Sub
  ↓
WebSocket
  ↓
Connected clients receive updated results
5. Real-Time Results

Redis stores the live vote counts and publishes an update whenever a vote is submitted.

The Go WebSocket handler subscribes to the Redis update channel and sends the new results to connected clients.

Therefore, users watching the poll see updated results without refreshing the page.

🗄️ Data Storage
MongoDB

MongoDB stores persistent application data including:

User accounts
Poll questions
Poll options
Poll ownership
Poll creation information
Redis

Redis is used for the real-time voting layer.

For each poll, vote counts are maintained using Redis and updates are published through Redis Pub/Sub.

This allows multiple connected clients to receive live changes efficiently.

🔐 Authentication & Validation

The backend includes:

User registration
User login
Password hashing using bcrypt
JWT authentication
Protected poll creation
Input validation
Poll and option validation

Authentication tokens are sent using the Authorization header.

📁 Project Structure
live-polling/
│
├── backend/
│   ├── config/
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── poll.go
│   │   └── realtime.go
│   ├── middleware/
│   │   └── auth.go
│   ├── models/
│   │   ├── poll.go
│   │   └── user.go
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   └── .gitignore
│
├── frontend/
│   ├── src/
│   │   ├── App.jsx
│   │   ├── App.css
│   │   ├── index.css
│   │   └── main.jsx
│   ├── public/
│   ├── package.json
│   └── vite.config.js
│
├── README.md
└── .gitignore
⚙️ Local Setup
Prerequisites

Install:

Node.js
Go
MongoDB
Redis
Clone the repository
git clone https://github.com/iniyailamaran31-tech/live-polling.git
cd live-polling
Backend
cd backend
go mod download
go run main.go

The backend runs locally on:

http://localhost:8080
Frontend

Open another terminal:

cd frontend
npm install
npm run dev

The frontend will be available at:

http://localhost:5173
🔑 Environment Variables

Create a .env file inside the backend directory.

Example:

MONGODB_URI=your_mongodb_connection_string
MONGODB_DATABASE=livepoll
REDIS_URL=your_redis_connection_string
JWT_SECRET=your_jwt_secret

For the frontend, configure:

VITE_API_URL=http://localhost:8080

Never commit real credentials or secret keys to GitHub.

🌐 Deployment
Frontend

The React application is deployed using Vercel.

Backend

The Go/Gin API is deployed using Render.

Database

MongoDB Atlas is used for cloud database storage.

Real-Time Infrastructure

Upstash Redis is used for Redis-based real-time vote counting and Pub/Sub.

🧪 Real-Time Test

The application can be tested using two browser windows:

Create a poll.
Copy the public poll link.
Open the link in another browser window.
Vote from the second window.
Observe the results update in the first window without refreshing.
🎯 Project Goal

The goal of this project is to demonstrate a complete real-time polling system using a React frontend, Go backend, MongoDB database, and Redis-powered real-time communication.

👩‍💻 Project

Live Polling Tool

Built as part of a developer assignment.