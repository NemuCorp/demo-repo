# Update README

**Issue:** [NemuCorp/demo-repo#14](https://github.com/NemuCorp/demo-repo/issues/14)

## Problem

The current README is a stub — it only contains the project name and a one-line description (`Create a Basic Eccommerce Website`). Visitors and potential contributors have no introduction to the project, no visibility into its features, no screenshots to understand the UI, and no guidance on how to set up or run it locally.

## Goals

- Write a clear project introduction that explains what "Demo Store" is and its purpose.
- Showcase the application with screenshots of key pages (home page, product detail, cart, login/register, and admin panel).
- List the features the application provides and the use cases it supports.
- Document the prerequisites and step-by-step instructions to set up and run both the server and client.
- Keep the tone friendly, structured, and appropriate for a public GitHub README.

## Approach

### README Structure

The updated README will follow a common open-source project layout:

```
# Demo Store
> Short tagline or project summary

## Introduction
Brief description of the project — a full-stack ecommerce demo application.

## Screenshots
Grid or list of images showing each major page:
- Home page (product grid)
- Product detail page
- Shopping cart
- Login / Register
- Admin dashboard (product management)

## Features
Bulleted list of implemented features grouped by domain:
- **Products**: browse, view details, images, pricing, stock
- **Auth**: register, login, logout, bcrypt-hashed passwords
- **Cart**: add, update quantity, remove items, view total
- **Admin**: dashboard, create products, product inventory table

## Use Cases
Short list of user stories or use cases (e.g., "As a shopper, I can browse products and add them to my cart").

## Tech Stack
Two-column or badge-based listing of technologies: Go, Gin, PostgreSQL, React, TypeScript, React Router.

## Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL

## Setup & Run

### Backend
1. Set environment variables (DATABASE_URL, PORT, LOG_MODE)
2. Run database migrations
3. Start the server

### Frontend
1. Install dependencies
2. Start the dev server

### Using Docker (if applicable) or manual steps only

Each step will include concrete commands and expected output.

## Project Structure
A trimmed-down tree of the repo (optional but helpful for contributors).
```

### Image Handling

Page screenshots will be stored in a `docs/images/` directory and referenced via relative paths in the README. Image filenames will be descriptive (e.g., `home-page.png`, `product-detail.png`). If actual screenshots are not available yet, placeholders or instructions for capturing them will be noted.

### Constraints

- The README update should only add the sections described above — it should not change any code, configuration, or existing documentation beyond `README.md` and `docs/images/`.
- References to specific technologies and versions must match what is actually present in `client/package.json` and `server/go.mod`.
- All setup instructions must be verifiable against the existing `main.go` environment variable names (DATABASE_URL, PORT, LOG_MODE) and `client/package.json` scripts.

## Risks

| Risk | Mitigation |
|------|------------|
| Screenshots become stale as UI evolves | Store images with version in filename or update alongside UI changes; note that images are illustrative |
| Setup instructions drift from actual behavior | Derive all commands from existing entrypoints (`main.go`, `package.json scripts`); cross-check before finalizing |
| README becomes too long or verbose | Use collapsible sections or keep language concise; link to separate docs if needed |

## Validation

- [ ] README contains a project introduction paragraph.
- [ ] README includes screenshots of at least 4 key pages (home, product detail, cart, admin).
- [ ] README lists features grouped by domain.
- [ ] README includes use cases or user stories.
- [ ] README documents prerequisites (Go, Node.js, PostgreSQL).
- [ ] README provides step-by-step setup/run instructions for backend and frontend that match actual source code entrypoints.
- [ ] README mentions the tech stack.
- [ ] Images are stored in `docs/images/` and referenced with working relative paths.
