# StudyHub - Project Documentation for AI Assistant

## Project Overview
StudyHub is a full-stack web application for managing academic modules, resources, and study materials. Users can upload files and links organized by modules, weeks, and academic terms.

**Tech Stack:**
- **Backend:** Go (Golang) with Chi router, PostgreSQL database, AWS S3 for file storage
- **Frontend:** React + TypeScript, React Router v6, Tailwind CSS, shadcn/ui components

---

## Project Structure

```
StudyHub/
├── backend/
│   ├── internal/
│   │   ├── http/              # HTTP handlers and routing
│   │   │   ├── http.go        # Main router setup
│   │   │   ├── resource_handlers.go
│   │   │   ├── module_handlers.go
│   │   │   └── user_handlers.go
│   │   ├── resources/         # Resource business logic
│   │   │   ├── type.go        # Resource types and structs
│   │   │   ├── repository.go  # Database queries
│   │   │   └── service.go     # Business logic
│   │   ├── modules/           # Module business logic
│   │   ├── users/             # User management
│   │   └── auth/              # Authentication logic
│   └── cmd/
│       └── server/
│           └── main.go        # Application entry point
│
└── frontend/
    ├── src/
    │   ├── pages/             # Page components
    │   │   ├── ModulesPage.tsx
    │   │   ├── ModuleDetailPage.tsx
    │   │   ├── WeekDetailPage.tsx
    │   │   ├── UserProfilePage.tsx
    │   │   ├── AcademicTermsPage.tsx
    │   │   ├── LoginPage.tsx
    │   │   └── RegisterPage.tsx
    │   ├── components/
    │   │   ├── layout/
    │   │   │   ├── Layout.tsx        # Main layout wrapper
    │   │   │   ├── Header.tsx        # Top navigation bar
    │   │   │   └── Sidebar.tsx       # Side navigation menu
    │   │   ├── resources/
    │   │   │   ├── ResourceCard.tsx         # Resource card for week view
    │   │   │   ├── ResourceList.tsx         # Grid of resource cards
    │   │   │   ├── UserResourceCard.tsx     # Resource card for user profile
    │   │   │   ├── UserResourceList.tsx     # Grid for user resources
    │   │   │   ├── ResourceUploadDialog.tsx # File upload modal
    │   │   │   └── ResourceLinkDialog.tsx   # Link upload modal
    │   │   ├── ui/            # shadcn/ui components
    │   │   └── auth/          # Auth-related components
    │   ├── api/               # API client functions
    │   │   ├── client.ts      # Axios instance with auth interceptor
    │   │   ├── resources.ts   # Resource API calls
    │   │   ├── modules.ts     # Module API calls
    │   │   └── auth.ts        # Auth API calls
    │   ├── context/
    │   │   └── AuthContext.tsx # Authentication state management
    │   ├── types/
    │   │   └── index.ts       # TypeScript type definitions
    │   └── App.tsx            # Main app component with routing
    └── package.json
```

---

## Key Backend Endpoints

### Resources
- `GET /api/v1/resources/weeks/{week_id}` - Get all resources for a week
  - Returns: `[]ResourceWithUser` (includes UserName and UserID)
- `GET /api/v1/resources/users/{user_id}` - Get all resources uploaded by a user
  - Returns: `[]UserResources` (includes module context: ModuleName, Semester, Year, WeekNumber)
- `GET /api/v1/resources/{object_id}` - Get presigned S3 URL for file download
  - Returns: `{"url": "presigned-s3-url"}`
- `POST /api/v1/resources/file/{week_id}` - Upload a file resource (multipart/form-data)
- `POST /api/v1/resources/link/{week_id}` - Upload a link resource (JSON: name, url)

### Modules
- `GET /api/v1/modules` - List all modules
- `GET /api/v1/modules/{id}` - Get module with runs and weeks
- `POST /api/v1/modules` - Create a new module

### Users
- `GET /api/v1/users/{user_id}` - Get user by ID
- `POST /api/v1/users` - Register a new user

### Auth
- `POST /api/v1/auth/login` - Login (returns JWT token)

---

## Backend Type Definitions

### ResourceWithUser (returned by GET /resources/weeks/{id})
```go
type ResourceWithUser struct {
    ID           uuid.UUID
    WeekID       uuid.UUID
    UserID       uuid.UUID
    UserName     string       // User's first name
    ObjectID     *uuid.UUID   // For files
    ExternalLink *string      // For links
    ResourceType ResourceType // "file" | "link" | "note"
    Name         string
    CreatedAt    time.Time
}
```

### UserResources (returned by GET /resources/users/{id})
```go
type UserResources struct {
    ID           uuid.UUID
    WeekID       uuid.UUID
    UserID       uuid.UUID
    ModuleName   string       // Module name for context
    Semester     string       // e.g., "Spring", "Fall"
    Year         int          // e.g., 2024
    WeekNumber   int          // Week number in module
    ObjectID     *uuid.UUID   // For files (nullable)
    ExternalLink *string      // For links (nullable)
    ResourceType ResourceType // "file" | "link" | "note"
    Name         string
    CreatedAt    time.Time
}
```

---

## Frontend Type Definitions

Located in: `/frontend/src/types/index.ts`

### Key Types
```typescript
// Resource types
export type ResourceType = 'file' | 'link' | 'note'

// Resource from week view
export interface Resource {
  ID: string
  WeekID: string
  UserID: string
  UserName: string        // Added for uploader display
  ResourceType: ResourceType
  Hash: string
  Name: string
  Url: string            // For links
  ObjectID: string       // For files
  CreatedAt: string
  UpdatedAt: string
}

// Resource from user profile view
export interface UserResource {
  ID: string
  WeekID: string
  UserID: string
  ModuleName: string
  Semester: string
  Year: number
  WeekNumber: number
  ObjectID: string | null
  ExternalLink: string | null
  ResourceType: ResourceType
  Name: string
  CreatedAt: string
}

// User
export interface User {
  ID: string
  Email: string
  FirstName: string
  LastName: string
  IsAdmin: boolean
  CreatedAt: string
  UpdatedAt: string
}

// Module
export interface Module {
  ID: string
  Code: string
  Name: string
  DepartmentName: string
  CreatedAt: string
  UpdatedAt: string
}

// Module Run
export interface ModuleRun {
  ID: string
  ModuleID: string
  Year: number
  Semester: string
  IsActive: boolean
  CreatedAt: string
}

// Week
export interface Week {
  ID: string
  ModuleRunID: string
  Number: number
}
```

---

## Frontend Routing

Located in: `/frontend/src/App.tsx`

### Public Routes
- `/login` - Login page
- `/register` - Registration page

### Protected Routes (require authentication)
All wrapped in `<ProtectedRoute>` and `<Layout>`:
- `/` - Redirects to `/modules`
- `/home` - Home/Dashboard page
- `/modules` - List of all modules
- `/modules/:id` - Module detail page with runs and weeks
- `/modules/:moduleId/weeks/:weekId` - Week detail page with resources
- `/academic-terms` - Academic terms management
- `/users/:userId` - User profile page with uploaded resources

---

## Authentication

### How It Works
1. User logs in via `/login`
2. Backend returns JWT token
3. Token stored in `localStorage` as `auth_token`
4. User data cached in `localStorage` as `auth_user`
5. `AuthContext` provides current user to all components
6. `apiClient` automatically adds token to all API requests via interceptor

### AuthContext Hook
```typescript
const { user, isAuthenticated, login, logout } = useAuth()
```

### API Client
Located in: `/frontend/src/api/client.ts`
- Axios instance with base URL
- Request interceptor adds `Authorization: Bearer <token>` header
- Response interceptor handles 401 errors (redirects to login)

---

## Common Patterns

### Page Components
All page components follow this pattern:
```typescript
import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Loader2 } from 'lucide-react'
import { useToast } from '@/components/ui/toast'

const MyPage: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { showToast } = useToast()
  
  const [data, setData] = useState<MyType[]>([])
  const [loading, setLoading] = useState(true)

  const loadData = async () => {
    try {
      setLoading(true)
      const result = await myApi.getData(id)
      setData(result)
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to load data',
        'error'
      )
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadData()
  }, [id])

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Page content */}
    </div>
  )
}

export default MyPage
```

### API Client Pattern
Located in: `/frontend/src/api/`
```typescript
import apiClient from './client'
import type { MyType } from '@/types'

export const myApi = {
  getItems: async (): Promise<MyType[]> => {
    const response = await apiClient.get('/items')
    return response.data
  },

  createItem: async (data: CreateItemRequest): Promise<void> => {
    await apiClient.post('/items', data)
  },
}
```

### UI Components
Using shadcn/ui components with consistent styling:
- `Card`, `CardHeader`, `CardTitle`, `CardContent`, `CardFooter`
- `Button` with variants: `default`, `outline`, `ghost`
- `Badge` with variants: `default`, `outline`, `secondary`
- `Tabs`, `TabsList`, `TabsTrigger`, `TabsContent`
- Icons from `lucide-react`

---

## Recent Features Implemented

### 1. Copy Link Functionality (ResourceCard)
- Copy button next to "Open Link" for link resources
- Uses `navigator.clipboard.writeText()`
- Shows checkmark icon for 2 seconds after successful copy

### 2. User Profile Pages
- Users can view their own profile from sidebar ("My Profile")
- Clicking on uploader names navigates to their profile (`/users/{userId}`)
- Profile shows:
  - User name, email, resource count
  - Tabs for Files and Links
  - Each resource displays module context (name, semester, year, week)
  - Action buttons (Download, Open Link, Copy)

### 3. Uploader Attribution (ResourceCard)
- Footer in each resource card shows "Uploaded by [Name]"
- Name is clickable, links to user's profile
- Uses `UserName` field from backend

---

## Important Notes

### Backend
- **DO NOT modify backend code** unless explicitly requested
- Backend returns `UserName` (first name) in week resources
- Backend returns full context (module, semester, year, week) in user resources
- Week navigation requires `ModuleID` which is NOT included in user resources response

### Frontend
- All file paths should use absolute paths (not relative)
- Components use TypeScript with strict typing
- CSS uses Tailwind utility classes
- State management via React hooks and Context API
- No Redux or other state management libraries

### File Operations
- File uploads handled via multipart/form-data
- File downloads use presigned S3 URLs (backend generates, frontend opens in new tab)
- Links open in new tab with `window.open(url, '_blank', 'noopener,noreferrer')`

### Code Style
- Use functional components with hooks (not class components)
- Use `const` for component declarations
- Use TypeScript interfaces from `/types/index.ts`
- Follow existing naming conventions (PascalCase for components, camelCase for functions)
- Use lucide-react for icons
- Use shadcn/ui components for UI elements

---

## Development Commands

### Backend
```bash
cd backend
go run cmd/server/main.go
```

### Frontend
```bash
cd frontend
npm install
npm run dev
```

---

## Future Considerations

### Not Yet Implemented
1. Week navigation from user profile resources (needs ModuleID in response or client-side lookup)
2. Resource editing/deletion
3. Admin-only features
4. Resource notes (type exists but not implemented)
5. Search functionality
6. Pagination for large resource lists

### Known Limitations
- User resources endpoint doesn't include ModuleID (only WeekID)
- No error boundaries in React components
- No offline support
- No real-time updates (websockets)

---

## Quick Reference: Where to Find Things

**Need to modify routing?** → `/frontend/src/App.tsx`

**Need to add a new API endpoint call?** → `/frontend/src/api/`

**Need to add a new type?** → `/frontend/src/types/index.ts`

**Need to modify sidebar navigation?** → `/frontend/src/components/layout/Sidebar.tsx`

**Need to modify authentication logic?** → `/frontend/src/context/AuthContext.tsx`

**Need to see backend endpoint handlers?** → `/backend/internal/http/`

**Need to see database queries?** → `/backend/internal/*/repository.go`

**Need to see backend type definitions?** → `/backend/internal/*/type.go`

---

## Common Tasks

### Adding a New Page
1. Create page component in `/frontend/src/pages/MyPage.tsx`
2. Add route in `/frontend/src/App.tsx`
3. Optionally add navigation link in `/frontend/src/components/layout/Sidebar.tsx`

### Adding a New API Endpoint (Frontend)
1. Define types in `/frontend/src/types/index.ts`
2. Add API method in `/frontend/src/api/myApi.ts`
3. Use in page component with `useState` and `useEffect`

### Adding a New Resource Type Display
1. Update `getResourceIcon()` and `getResourceTypeBadge()` functions
2. Add handling in `ResourceCard` or `UserResourceCard` components
3. Update type checking logic if needed

---

**Last Updated:** 2026-02-03
**Maintained By:** AI Assistant for efficient context loading
