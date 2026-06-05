# QuestHub

![QuestHub Banner](https://via.placeholder.com/1200x400?text=QuestHub+Preview)

> **A modern, powerful platform for managing tabletop RPG campaigns with real-time capabilities, powered by SvelteKit and Supabase.**

![SvelteKit](https://img.shields.io/badge/SvelteKit-2.0-FF3E00?style=for-the-badge&logo=svelte)
![Supabase](https://img.shields.io/badge/Supabase-Enabled-3ECF8E?style=for-the-badge&logo=supabase)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

## 📖 Overview

**QuestHub** is a comprehensive Virtual Tabletop (VTT) and campaign management companion designed to streamline the TTRPG experience for Game Masters and Players alike. Built with performance and real-time interaction in mind, it provides a seamless campaign organizer and real-time gaming environment.

---

## 🛠 Tech Stack

### Frontend
-   **Framework**: [SvelteKit](https://kit.svelte.dev/) - Svelte 5, full-stack Svelte framework.
-   **Styling**: [TailwindCSS v4](https://tailwindcss.com/) - Utility-first CSS.
-   **Type Safety**: TypeScript.

### Backend (Serverless / Managed)
-   **Database**: [PostgreSQL](https://www.postgresql.org/) hosted on Supabase.
-   **Authentication**: Supabase Auth (Email + Google Social login).
-   **Real-time**: Supabase Realtime Channels (heartbeats, Broadcast, and Change Data Capture).
-   **Storage**: Supabase Storage (for campaign assets and avatars).
-   **Access Control**: Row Level Security (RLS) policies defined in database schema.

---

## 🚀 Getting Started

Follow these instructions to set up the project locally for development and testing purposes.

### Prerequisites

-   [Bun](https://bun.sh/) (highly recommended) or Node.js.
-   A [Supabase Cloud](https://supabase.com/) project.

### Installation & Database Setup

1.  **Clone the repository**
    ```bash
    git clone https://github.com/yourusername/questhub.git
    cd questhub
    ```

2.  **Initialize local Supabase CLI**
    If you haven't linked the CLI to your project yet:
    ```bash
    bunx supabase link --project-ref <your-project-ref>
    ```

3.  **Apply database migrations**
    Push the SQL schema (tables, RLS policies, views, and RPCs) to your Supabase Cloud instance:
    ```bash
    bunx supabase db push
    ```

### Frontend Configuration

1.  **Configure environment variables**
    Create a `.env` file inside the `frontend` folder:
    ```bash
    cd frontend
    touch .env
    ```
    Populate it with your Supabase credentials:
    ```env
    PUBLIC_SUPABASE_URL=https://<your-project-ref>.supabase.co
    PUBLIC_SUPABASE_ANON_KEY=<your-anon-key>
    ```

2.  **Install dependencies**
    ```bash
    bun install
    ```

3.  **Run SvelteKit dev server**
    ```bash
    bun run dev
    ```
    Your client is running on `http://localhost:5173`.

---

## 📂 Project Structure

```bash
questhub/
├── frontend/           # SvelteKit application
│   ├── src/
│   │   ├── lib/        # Reusable components, Supabase client & store wrappers
│   │   └── routes/     # App pages & layouts
│   └── package.json
└── supabase/           # Supabase migrations & configurations
    └── migrations/     # PostgreSQL schema migrations
```

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1.  Fork the project
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`)
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4.  Push to the branch (`git push origin feature/AmazingFeature`)
5.  Open a Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.
