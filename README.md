<p align="center">
  <img src="https://em-content.zobj.net/source/apple/391/guitar_1f3b8.png" width="120" alt="Buskr Logo" />
</p>

<h1 align="center">Buskr</h1>

<p align="center">
  <b>Smart booking platform for street musicians</b><br/>
  <sub>Telegram Bot · Slot Management · Interactive Map · Karma System</sub>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.25" />
  <img src="https://img.shields.io/badge/Telegram-Bot%20API-26A5E4?style=flat-square&logo=telegram&logoColor=white" alt="Telegram Bot API" />
  <img src="https://img.shields.io/badge/PostgreSQL-PostGIS-4169E1?style=flat-square&logo=postgresql&logoColor=white" alt="PostgreSQL + PostGIS" />
  <img src="https://img.shields.io/badge/Redis-State%20Store-DC382D?style=flat-square&logo=redis&logoColor=white" alt="Redis" />
</p>

---

## What is Buskr?

**Buskr** is a Telegram-based platform that helps cities organize street music. Musicians discover open performance spots on an interactive map, reserve time slots, and check in — all through a conversational bot interface. Administrators manage locations, approve artists, and monitor activity in real time.

The project was born out of a real need: coordinating dozens of performers across multiple public spaces with limited schedules, noise restrictions, and no central system. Buskr replaces spreadsheets, group chats, and manual coordination with an automated, fair, and transparent booking flow.

---

## Core Concepts

### 🎸 For Musicians

- **Discover locations** on an interactive Google Maps–powered WebApp, or browse a text-based list filtered by noise compatibility.
- **Book time slots** up to several days in advance, with automatic collision detection — no double-bookings, no schedule conflicts.
- **Check in** when arriving at a spot to confirm presence. Missed check-ins are handled automatically (see Karma below).
- **Track performance stats** — total bookings, successful check-ins, and no-shows — all visible in a personal profile.

### 🛡️ For Administrators

- **Manage locations** — create, edit, activate/deactivate performance spots with geo-coordinates and noise limits.
- **Moderate musicians** — review applications (with optional video/audio portfolios), approve or reject, ban/unban, promote to admin or demote.
- **Browse participants** with cyclic sorting (by date, karma, role, or name) and full-text search.
- **Configure everything** via environment variables — booking limits, karma thresholds, advance booking window, and more.

### 🎯 Karma System

Every musician starts with a perfect karma score. The system rewards reliability and penalizes no-shows:

| Event | Effect |
|---|---|
| Successful check-in | Karma increases |
| Missed check-in (no-show) | Karma decreases |

Low karma triggers a visible warning in the musician's profile and a status indicator in admin views. All karma parameters — rewards, penalties, thresholds — are fully configurable through environment variables.

### 🔥 Hot Slots

When a musician misses their check-in and the slot gets cancelled, Buskr can automatically broadcast the freed-up time as a **hot slot** to all noise-compatible active musicians. First come, first served — ensuring performance spots never go to waste.

---

## Key Features

| Feature | Description |
|---|---|
| **Interactive Map** | Telegram WebApp with Google Maps integration for browsing locations and viewing schedules |
| **Smart Booking** | Collision detection, noise compatibility checks, adjacency radius enforcement for loud performers |
| **Automated Lifecycle** | Background scheduler handles reminders, check-in timeouts, booking completions, and hot slot broadcasts |
| **Noise Matching** | Three-tier noise system (Light / Medium / Hard) for both musicians and locations ensures compatibility |
| **Onboarding Flow** | Multi-step guided registration: name → noise level → optional media portfolio → admin approval |
| **Invite System** | Generate single-use invite links with preset noise levels and expiration times |
| **Multilingual** | Full i18n support with Russian and English locales |
| **Admin Panel** | In-bot admin interface for user management, location management, and real-time moderation |
| **Graceful Shutdown** | Signal handling with clean resource cleanup |

---

## Architecture

Buskr follows a **Clean Architecture** pattern with clear separation of concerns:

```
cmd/bot/              → Application entry point
internal/
├── config/           → Environment-based configuration (cleanenv)
├── domain/           → Core business entities and services
│   ├── user/         → User, roles, karma, invites
│   ├── booking/      → Booking lifecycle, noise compatibility
│   └── location/     → Locations, coordinates, noise limits
├── infrastructure/   → External adapters
│   ├── postgres/     → Repository implementations (lib/pq)
│   └── redis/        → State management for conversational flows
├── transport/        → Telegram bot layer
│   └── telegram/
│       ├── handlers/ → Command and callback handlers
│       ├── middleware/ → Auth middleware
│       ├── notifier/ → Admin notification service
│       └── render/   → Response rendering engine
├── usecase/          → Application-specific business rules
│   ├── auth/         → Authentication and main menu
│   ├── booking/      → Booking creation, schedule viewing
│   ├── onboarding/   → Registration flow
│   ├── profile/      → User profile management
│   ├── admin/        → Admin panel entry
│   ├── adminloc/     → Location administration
│   └── adminuser/    → User administration
├── worker/           → Background scheduler (reminders, completions)
├── i18n/             → Internationalization (go-i18n)
└── mapimg/           → Static map image generation
web/                  → Telegram WebApp (interactive map)
migrations/           → SQL migrations (goose-compatible)
```

**Key architectural decisions:**
- Domain layer has zero external dependencies — all configuration is injected through constructors
- Conversational state (multi-step flows) is managed via Redis, keeping the bot stateless between requests
- The response layer uses an abstract `Reply` struct, decoupling business logic from Telegram API specifics
- PostGIS is used for geospatial queries (location proximity, adjacency radius enforcement)

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| Bot Framework | [telebot v3](https://github.com/tucnak/telebot) |
| Database | PostgreSQL 15 + PostGIS |
| State Store | Redis 7 |
| Configuration | [cleanenv](https://github.com/ilyakaznacheev/cleanenv) + `.env` files |
| Localization | [go-i18n v2](https://github.com/nicksnyder/go-i18n) |
| Maps | Google Maps JavaScript API (Telegram WebApp) |
| Containerization | Docker Compose |

---

## Getting Started

```bash
# Clone the repository
git clone https://github.com/k4sper1love/buskr.git
cd buskr

# Start infrastructure
docker compose up -d

# Apply migrations
psql $DATABASE_DSN -f migrations/000001_init_schema.sql
psql $DATABASE_DSN -f migrations/000002_create_invites.sql

# Configure environment
cp .env.example .env
# Edit .env with your Telegram bot token and other settings

# Run
go run cmd/bot/main.go
```

---

## Configuration

All settings are managed through environment variables. See [`.env.example`](.env.example) for the full list.

Key configuration groups:

- **Telegram** — Bot token, admin chat ID, notification threads
- **Booking** — Max active bookings, advance days, hot slots toggle, no-show cancellation
- **Karma** — Max/min values, reward/penalty amounts, warning thresholds
- **Infrastructure** — PostgreSQL DSN, Redis connection, timezone

---

## License

This project is proprietary. All rights reserved.

---

<p align="center">
  <sub>Built with ❤️ for the street music community</sub>
</p>