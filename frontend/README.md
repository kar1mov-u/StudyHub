# StudyHub Frontend

A modern, intuitive web application for managing academic modules, module runs, and academic terms.

## Technology Stack

- **Framework**: React 18 with TypeScript
- **Build Tool**: Vite
- **UI Library**: Tailwind CSS + shadcn/ui components
- **Routing**: React Router v6
- **HTTP Client**: Axios
- **Icons**: Lucide React

## Features

### Dashboard
- Overview of active academic term
- Total modules count
- Quick actions to create modules and terms
- Recent modules list

### Modules Management
- List all modules with code, name, and department
- Create new modules
- View module details with runs and weeks
- Delete modules with confirmation
- Navigate to module details

### Module Runs
- Create runs for specific modules
- View run details with associated weeks
- Delete runs with confirmation
- Mark runs as active/inactive

### Academic Terms Management
- List all academic terms (Spring/Fall + Year)
- Create new terms
- Activate/deactivate terms
- Visual indicator for active term

### UI/UX Features
- Toast notifications for all actions
- Loading states for async operations
- Confirmation dialogs for destructive actions
- Responsive design
- Clean, modern interface
- Intuitive navigation with sidebar

## Getting Started

### Prerequisites

- Node.js 18+ and npm/yarn/pnpm

### Installation

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

### Building for Production

```bash
npm run build
```

The built files will be in the `dist` directory.

## API Integration

The frontend connects to the backend API at `http://localhost:8080/api/v1`. The Vite dev server is configured to proxy API requests to avoid CORS issues.

### Available API Endpoints

**Modules:**
- `GET /api/v1/modules` - List all modules
- `POST /api/v1/modules` - Create module
- `GET /api/v1/modules/{id}` - Get module details
- `DELETE /api/v1/modules/{id}` - Delete module

**Module Runs:**
- `GET /api/v1/modules/{moduleID}/runs` - List runs for a module
- `POST /api/v1/modules/{moduleID}/runs` - Create module run
- `GET /api/v1/module-runs/{id}` - Get run details
- `DELETE /api/v1/module-runs/{id}` - Delete module run

**Academic Terms:**
- `GET /api/v1/academic-terms` - List all terms
- `GET /api/v1/academic-terms/active` - Get active term
- `POST /api/v1/academic-terms` - Create term
- `PATCH /api/v1/academic-terms/{id}/activate` - Activate term
- `PATCH /api/v1/academic-terms/{id}/deactivate` - Deactivate term

## Project Structure

```
frontend/
├── src/
│   ├── api/                    # API client and service functions
│   │   ├── client.ts           # Axios instance with interceptors
│   │   ├── modules.ts          # Module API functions
│   │   ├── moduleRuns.ts       # Module Run API functions
│   │   └── academicTerms.ts    # Academic Term API functions
│   ├── components/
│   │   ├── ui/                 # shadcn/ui components
│   │   ├── layout/             # Layout components
│   │   ├── modules/            # Module-specific components
│   │   └── academic-terms/     # Academic term components
│   ├── pages/                  # Page components
│   │   ├── HomePage.tsx
│   │   ├── ModulesPage.tsx
│   │   ├── ModuleDetailPage.tsx
│   │   └── AcademicTermsPage.tsx
│   ├── types/                  # TypeScript type definitions
│   ├── lib/                    # Utility functions
│   ├── App.tsx                 # Main app component with routing
│   ├── main.tsx                # Entry point
│   └── index.css               # Global styles
├── package.json
├── tsconfig.json
├── tailwind.config.js
├── vite.config.ts
└── index.html
```

## Development Tips

1. **Type Safety**: All API responses are typed according to the backend models
2. **Error Handling**: Errors are caught and displayed as toast notifications
3. **Loading States**: All async operations show loading indicators
4. **Confirmation Dialogs**: Destructive actions require confirmation
5. **Navigation**: Use the sidebar to navigate between pages

## Common Tasks

### Adding a New Component

1. Create the component in the appropriate directory under `src/components/`
2. Use TypeScript for type safety
3. Import and use shadcn/ui components for consistency

### Adding a New API Endpoint

1. Add the function to the appropriate file in `src/api/`
2. Update types in `src/types/index.ts` if needed
3. Use the `useToast` hook for user feedback

### Styling

- Use Tailwind CSS utility classes
- Reference the design system variables defined in `index.css`
- Use shadcn/ui components for consistency

## Browser Support

- Chrome/Edge (latest)
- Firefox (latest)
- Safari (latest)

## License

This project is part of the StudyHub application.
