# StudyHub - Quick Start Guide

A full-stack application for managing academic modules, module runs, and academic terms.

## Project Structure

```
StudyHub/
├── backend/           # Go backend API
└── frontend/          # React TypeScript frontend
```

## Getting Started

### Backend Setup

The backend is a Go application with PostgreSQL database.

1. Navigate to the backend directory and follow the setup instructions
2. Ensure the backend is running on `http://localhost:8080`

### Frontend Setup

The frontend is a modern React application built with Vite and TypeScript.

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

4. Open your browser and navigate to `http://localhost:3000`

## Features Overview

### Dashboard
- View active academic term
- See total modules count
- Quick access to create modules and terms
- Recent modules overview

### Modules
- Create, view, and delete modules
- Each module has: Code, Name, Department
- Click on a module to view details and runs

### Module Runs
- Create runs for modules (e.g., Fall 2024, Spring 2025)
- View run details including weeks
- Mark runs as active/inactive
- Delete runs

### Academic Terms
- Manage academic terms (Spring/Fall + Year)
- Only one term can be active at a time
- Activate/deactivate terms
- Create new terms

## Technology Stack

### Frontend
- React 18 + TypeScript
- Vite (build tool)
- Tailwind CSS + shadcn/ui
- React Router v6
- Axios for API calls

### Backend
- Go
- Chi router
- PostgreSQL database

## API Endpoints

All endpoints are prefixed with `/api/v1`

**Modules:**
- `GET /modules` - List all modules
- `POST /modules` - Create module
- `GET /modules/{id}` - Get module details
- `DELETE /modules/{id}` - Delete module
- `GET /modules/{moduleID}/runs` - List runs for a module
- `POST /modules/{moduleID}/runs` - Create module run

**Module Runs:**
- `GET /module-runs/{id}` - Get run details
- `DELETE /module-runs/{id}` - Delete module run

**Academic Terms:**
- `GET /academic-terms` - List all terms
- `GET /academic-terms/active` - Get active term
- `POST /academic-terms` - Create term
- `PATCH /academic-terms/{id}/activate` - Activate term
- `PATCH /academic-terms/{id}/deactivate` - Deactivate term

## User Guide

### Creating Your First Module

1. Go to the Modules page (sidebar navigation)
2. Click "Create Module"
3. Fill in:
   - Module Code (e.g., CS101)
   - Module Name (e.g., Introduction to Computer Science)
   - Department Name (e.g., Computer Science)
4. Click "Create"

### Creating a Module Run

1. Click on a module to view its details
2. Click "Add Run"
3. Fill in:
   - Year (e.g., 2024)
   - Semester (Spring or Fall)
   - Set as active (optional)
4. Click "Create"

### Managing Academic Terms

1. Go to Academic Terms page
2. Click "Create Term"
3. Fill in:
   - Year
   - Semester (Spring or Fall)
4. Click "Create"
5. Click "Activate" on a term to make it the active term

## Development

### Frontend Development

```bash
cd frontend
npm run dev        # Start dev server
npm run build      # Build for production
npm run preview    # Preview production build
npm run lint       # Run ESLint
```

### Making Changes

1. All components are in `frontend/src/components/`
2. Pages are in `frontend/src/pages/`
3. API functions are in `frontend/src/api/`
4. Types are in `frontend/src/types/`

## Troubleshooting

### Frontend can't connect to backend

- Ensure the backend is running on `http://localhost:8080`
- Check that the proxy configuration in `vite.config.ts` is correct
- The frontend dev server proxies `/api` requests to the backend

### CORS errors

- The Vite dev server handles CORS through the proxy configuration
- No additional CORS setup should be needed

### TypeScript errors

- Run `npm install` to ensure all dependencies are installed
- Check that `tsconfig.json` is properly configured

## Next Steps

1. Start the backend server
2. Start the frontend development server
3. Create an academic term
4. Create some modules
5. Add runs to your modules
6. Explore the dashboard and features

For more detailed information, see the README files in the backend and frontend directories.
