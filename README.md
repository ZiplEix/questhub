# QuestHub

![QuestHub Banner](https://via.placeholder.com/1200x400?text=QuestHub+Preview)

> **A modern, powerful platform for managing tabletop RPG campaigns with real-time capabilities.**

![Go](https://img.shields.io/badge/Go-1.24-00ADD8?style=for-the-badge&logo=go)
![SvelteKit](https://img.shields.io/badge/SvelteKit-2.0-FF3E00?style=for-the-badge&logo=svelte)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?style=for-the-badge&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=for-the-badge&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

## 📖 Overview

**QuestHub** is a comprehensive Virtual Tabletop (VTT) and campaign management companion designed to streamline the TTRPG experience for Game Masters and Players alike. Built with performance and real-time interaction in mind, it bridges the gap between complex campaign organization and fluid in-game action.

Whether you're tracking character stats, managing intricate inventory systems, or rolling dice in real-time with your party, QuestHub provides a seamless, responsive interface powered by a robust backend.

## ✨ Key Features

-   **Real-Time Interactivity**: Instant updates for dice rolls, chat messages, and game state changes using WebSockets.
-   **Campaign Management**: Centralized hub for campaign notes, NPCs, locations, and lore.
-   **Dynamic Character Sheets**: Fully customizable character sheets with automated stat calculations and inventory tracking.
-   **Game Master Tools**: robust suite of GM tools including initiative tracking, secret rolls, and player management.
-   **Modern UI/UX**: A sleek, accessible, and responsive interface built with SvelteKit and TailwindCSS v4.
-   **Secure Authentication**: Integrated user management system.

## 🛠 Tech Stack

### Backend
-   **Language**: [Go (Golang)](https://go.dev/)
-   **Framework**: [Echo](https://echo.labstack.com/) - High performance, extensible web framework.
-   **Database**: [PostgreSQL](https://www.postgresql.org/) - Robust relational database with `pgx` driver.
-   **Real-time**: [Gorilla WebSocket](https://github.com/gorilla/websocket) - Efficient WebSocket implementation for Go.

### Frontend
-   **Framework**: [SvelteKit](https://kit.svelte.dev/) - Full-stack framework for building performant web apps.
-   **Language**: TypeScript - For type-safe, maintainable code.
-   **Styling**: [TailwindCSS v4](https://tailwindcss.com/) - Utility-first CSS framework for rapid UI development.
-   **Auth**: [Better Auth](https://www.better-auth.com/) - Comprehensive authentication solution.

### DevOps & Infrastructure
-   **Containerization**: Docker & Docker Compose for consistent environments.
-   **Tooling**: Makefiles for automated tasks (migrations, builds, dev server).

## 🚀 Getting Started

Follow these instructions to set up the project locally for development and testing purposes.

### Prerequisites

-   [Docker](https://www.docker.com/) & Docker Compose
-   [Go](https://go.dev/dl/) (v1.24+)
-   [Node.js](https://nodejs.org/) or [Bun](https://bun.sh/) (for frontend)

### Installation

1.  **Clone the repository**
    ```bash
    git clone https://github.com/yourusername/questhub.git
    cd questhub
    ```

2.  **Environment Setup**
    Ensure you have the necessary environment variables set up. Check `.env.example` if available, or confirm `backend/.env` and `frontend/.env` configurations.

### Running with Docker (Recommended)

The easiest way to get everything running is using Docker Compose.

```bash
make db
```
*This command starts the PostgreSQL database container.*

To run the full stack:
```bash
docker-compose up --build
```

### Manual Development Setup

If you prefer to run services individually:

1.  **Start Database**
    ```bash
    make db
    ```

2.  **Run Backend**
    ```bash
    make air
    # Runs the Go backend with live reload (Air)
    ```

3.  **Run Frontend**
    ```bash
    make front
    # Starts the SvelteKit development server
    ```

4.  **Run Migrations**
    ```bash
    make migrate-up
    # Applies database migrations
    ```

## 📂 Project Structure

```bash
questhub/
├── backend/            # Go API server
│   ├── config/         # App configuration
│   ├── controller/     # HTTP handlers
│   ├── models/         # Database models & structs
│   └── routes/         # API route definitions
├── frontend/           # SvelteKit application
│   ├── src/
│   │   ├── lib/        # Reusable components & utilities
│   │   └── routes/     # App pages & layouts
├── migrations/         # SQL migration files
└── docker-compose.yml  # Container orchestration
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1.  Fork the project
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`)
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4.  Push to the branch (`git push origin feature/AmazingFeature`)
5.  Open a Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

---

*Built with ❤️ by [ZiplEix](https://github.com/ziplEix)*
