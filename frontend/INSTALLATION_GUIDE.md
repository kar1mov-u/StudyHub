# StudyHub Frontend - Installation & Usage Guide

## Prerequisites

Before you begin, ensure you have:
- **Node.js** 18.0 or higher
- **npm** 9.0 or higher (comes with Node.js)
- Backend API running on `http://localhost:8080`

Check your versions:
```bash
node --version
npm --version
```

## Installation Steps

### Step 1: Navigate to Frontend Directory
```bash
cd frontend
```

### Step 2: Install Dependencies
```bash
npm install
```

This will install all required packages:
- React 18
- TypeScript
- Vite
- Tailwind CSS
- React Router
- Axios
- Lucide React (icons)
- And more...

### Step 3: Verify Installation
```bash
npm list --depth=0
```

You should see all dependencies listed without errors.

## Running the Application

### Development Mode

Start the development server with hot reload:
```bash
npm run dev
```

You should see output like:
```
  VITE v5.0.8  ready in 500 ms

  ➜  Local:   http://localhost:3000/
  ➜  Network: use --host to expose
  ➜  press h to show help
```

Open your browser and navigate to `http://localhost:3000`

### Production Build

Build for production:
```bash
npm run build
```

This creates an optimized build in the `dist/` directory.

Preview the production build:
```bash
npm run preview
```

## First Time Usage

### 1. Check Backend Connection

Make sure your backend is running on `http://localhost:8080`. You can verify by visiting:
```
http://localhost:8080/api/v1/modules
```

If you see a JSON response, the backend is working!

### 2. Access the Frontend

Open `http://localhost:3000` in your browser.

### 3. Create Your First Academic Term

1. Click "Academic Terms" in the sidebar
2. Click "Create Term" button
3. Fill in:
   - Year: 2024 (or current year)
   - Semester: "fall" or "spring"
4. Click "Create"
5. Click "Activate" on the newly created term

### 4. Create Your First Module

1. Click "Modules" in the sidebar
2. Click "Create Module" button
3. Fill in:
   - Code: e.g., "CS101"
   - Name: e.g., "Introduction to Computer Science"
   - Department: e.g., "Computer Science"
4. Click "Create"

### 5. Add a Run to Your Module

1. Click on the module card you just created
2. Click "Add Run" button
3. Fill in:
   - Year: 2024
   - Semester: "fall"
   - Check "Set as active run" if desired
4. Click "Create"

### 6. Explore the Dashboard

1. Click "Dashboard" in the sidebar
2. You should see:
   - Active academic term
   - Total modules count
   - Recent modules list

Congratulations! You're now using StudyHub!

## Common Commands

### Development
```bash
npm run dev          # Start dev server
npm run build        # Build for production
npm run preview      # Preview production build
npm run lint         # Run ESLint
```

### Cleaning Up
```bash
rm -rf node_modules  # Remove dependencies
rm -rf dist          # Remove build output
npm install          # Reinstall dependencies
```

## Troubleshooting

### Problem: "npm: command not found"

**Solution**: Install Node.js from https://nodejs.org/

### Problem: Dependencies not installing

**Solution**: 
```bash
# Clear npm cache
npm cache clean --force

# Delete node_modules and package-lock.json
rm -rf node_modules package-lock.json

# Reinstall
npm install
```

### Problem: "Cannot connect to backend"

**Solution**:
1. Verify backend is running: `curl http://localhost:8080/api/v1/modules`
2. Check console for CORS errors
3. Ensure Vite proxy is configured (it should be)

### Problem: Port 3000 already in use

**Solution**:
```bash
# Find process using port 3000
lsof -i :3000

# Kill the process (replace PID with actual process ID)
kill -9 PID

# Or use a different port
npm run dev -- --port 3001
```

### Problem: TypeScript errors

**Solution**:
```bash
# Rebuild TypeScript
npm run build

# Check for type errors
npx tsc --noEmit
```

### Problem: Tailwind styles not working

**Solution**:
1. Ensure `index.css` is imported in `main.tsx`
2. Check `tailwind.config.js` content paths
3. Restart dev server

## Usage Guide

### Navigation

Use the sidebar to navigate between pages:
- **Dashboard**: Overview and quick stats
- **Modules**: Manage all modules
- **Academic Terms**: Manage terms

### Creating Items

All create actions open a dialog/modal:
1. Click "Create" button
2. Fill in the form
3. Click "Create" in dialog
4. Success toast appears
5. List refreshes automatically

### Deleting Items

All delete actions require confirmation:
1. Click trash icon or "Delete" button
2. Confirm in dialog
3. Success toast appears
4. List refreshes automatically

### Viewing Details

Click on a module card to view:
- Module information
- Active run
- Weeks associated with the run
- Option to add new runs

### Managing Academic Terms

Only one term can be active at a time:
1. Create multiple terms
2. Click "Activate" on desired term
3. Previously active term becomes inactive
4. Click "Deactivate" to deactivate current term

## Features Overview

### Dashboard
- **Active Term Card**: Shows current academic term
- **Modules Count**: Total number of modules
- **Quick Actions**: Fast access to create items
- **Recent Modules**: Last 5 modules created

### Modules Page
- **Grid Layout**: Visual card display
- **Create Button**: Add new modules
- **Module Cards**: Show code, name, department
- **Delete Option**: Remove modules with confirmation
- **Click to View**: Navigate to module details

### Module Detail Page
- **Back Navigation**: Return to modules list
- **Module Header**: Code, name, department
- **Active Run Display**: Current run with details
- **Week Badges**: Visual representation of weeks
- **Add Run**: Create new runs
- **Delete Run**: Remove runs with confirmation

### Academic Terms Page
- **Grid Layout**: Visual card display
- **Create Button**: Add new terms
- **Term Cards**: Show semester and year
- **Active Badge**: Visual indicator
- **Activate/Deactivate**: Switch active term

## Keyboard Shortcuts

- **Esc**: Close dialogs/modals
- **Enter**: Submit forms (when focused)
- **Tab**: Navigate between form fields

## Best Practices

### When Creating Modules
- Use clear, consistent naming
- Include department for organization
- Use standard code format (e.g., DEPT###)

### When Managing Terms
- Keep only one term active
- Create terms in advance
- Follow semester naming: "spring" or "fall"

### When Creating Runs
- Match run to academic term
- Set active run for current semester
- One active run per module recommended

## Data Management

### Importing Data
Currently manual entry only. Future versions may support:
- CSV import
- Bulk creation
- API integration

### Exporting Data
Data is stored in backend. Access via API:
- GET endpoints return all data
- Can be saved/processed externally

### Backup
Data is in the backend database:
- Backend handles persistence
- Frontend is stateless
- Safe to rebuild/redeploy frontend

## Performance Tips

### For Best Performance
1. **Use Modern Browser**: Chrome, Firefox, Safari, Edge
2. **Clear Cache**: If experiencing issues
3. **Stable Network**: For API calls
4. **Backend Performance**: Ensure backend is optimized

### Loading States
All async operations show loading indicators:
- Spinner while fetching data
- "Creating..." text on buttons
- Disabled states during operations

## Browser Compatibility

Tested and verified on:
- ✓ Chrome 100+
- ✓ Firefox 100+
- ✓ Safari 15+
- ✓ Edge 100+

Not supported:
- ✗ Internet Explorer
- ✗ Very old browser versions

## Mobile Support

The interface is responsive and works on:
- Mobile phones (portrait/landscape)
- Tablets
- Desktop screens
- Large monitors

## Accessibility

Built with accessibility in mind:
- Semantic HTML
- Proper ARIA labels
- Keyboard navigation
- Color contrast compliance
- Screen reader friendly

## Getting Help

### Documentation
- `README.md`: Overview and setup
- `IMPLEMENTATION_SUMMARY.md`: Technical details
- `FILE_STRUCTURE.md`: File organization
- This guide: Installation and usage

### Common Issues
Check the Troubleshooting section above.

### Code Examples
Look at existing components in `src/components/` for patterns.

## Next Steps

Now that you have the frontend running:

1. ✓ Create academic terms
2. ✓ Add your modules
3. ✓ Create module runs
4. ✓ Explore the dashboard
5. ✓ Customize as needed

Enjoy using StudyHub!
