# User Deck Cards Feature - Implementation Summary

## Overview
This feature allows users to create personal flashcard decks for each week. Users can:
- Add auto-generated flashcards (from PDF files) to their personal deck
- Create custom flashcards from scratch
- Edit flashcards in their deck
- Track study progress (review count, difficulty rating)
- View deck statistics

## Architecture

### Database Schema
**Table: `user_deck_cards`**
- `id` (UUID, PK) - Unique card identifier
- `user_id` (UUID, FK -> users) - Card owner
- `week_id` (UUID, FK -> weeks) - Associated week
- `source_flashcard_id` (UUID, FK -> flashcards, nullable) - Link to auto-generated card
- `front` (TEXT) - Question text
- `back` (TEXT) - Answer text
- `is_custom` (BOOLEAN) - True if user-created, false if copied from auto-generated
- `last_reviewed_at` (TIMESTAMP, nullable) - Last study session
- `review_count` (INT, default 0) - Number of times reviewed
- `difficulty_rating` (INT, 1-5, nullable) - User's difficulty assessment
- `created_at` (TIMESTAMP) - Creation timestamp
- `updated_at` (TIMESTAMP) - Last update timestamp
- **Unique constraint**: (user_id, week_id, source_flashcard_id) - Prevents duplicate additions

**Indexes:**
- `idx_user_deck_cards_user_week` on (user_id, week_id)
- `idx_user_deck_cards_source` on (source_flashcard_id)

### Backend Structure

#### 1. Types (`internal/content/types.go`)
- `UserDeckCard` - Core card struct with all fields
- `AddCardToDeckRequest` - Request to add auto-generated card
- `CreateCustomCardRequest` - Request to create custom card
- `UpdateCardRequest` - Request to update card content
- `RecordReviewRequest` - Request to record study session
- `DeckStats` - Statistics response struct

#### 2. Repository Layer (`internal/content/repository.go`)
Database operations:
- `AddCardToUserDeck` - Copy auto-generated card to user deck
- `CreateCustomCardInDeck` - Create new custom card
- `RemoveCardFromUserDeck` - Delete card from deck
- `GetUserDeckForWeek` - Retrieve all cards for a week
- `GetUserDeckCard` - Get single card by ID
- `UpdateUserDeckCard` - Modify card content
- `RecordCardReview` - Update review statistics
- `GetDeckStatistics` - Aggregate deck stats

#### 3. Service Layer (`internal/content/service.go`)
Business logic:
- Validates flashcard existence before adding
- Prevents empty front/back for custom cards
- Enforces authorization (users can only manage their own cards)
- Validates difficulty ratings (1-5 scale)
- Handles duplicate prevention via DB constraints

#### 4. HTTP Handlers (`internal/http/content_handlers.go`)
Request handling:
- Extracts user ID from JWT token
- Parses URL parameters
- Validates request bodies
- Returns appropriate status codes
- Handles errors with meaningful messages

## API Endpoints

All endpoints require authentication via `Authorization: Bearer <token>` header.

### 1. Add Auto-Generated Card to Deck
```
POST /api/v1/decks/weeks/{week_id}/cards
Content-Type: application/json

Body:
{
  "flashcard_id": "uuid-of-auto-generated-card"
}

Response: 201 Created
{
  "data": null
}

Errors:
- 400: Invalid week_id, flashcard_id, or request body
- 404: Flashcard not found
- 409: Card already in deck
- 500: Server error
```

### 2. Create Custom Card
```
POST /api/v1/decks/weeks/{week_id}/cards/custom
Content-Type: application/json

Body:
{
  "front": "What is the capital of France?",
  "back": "Paris"
}

Response: 201 Created
{
  "data": {
    "ID": "uuid",
    "UserID": "uuid",
    "WeekID": "uuid",
    "SourceFlashcardID": null,
    "Front": "What is the capital of France?",
    "Back": "Paris",
    "IsCustom": true,
    "LastReviewedAt": null,
    "ReviewCount": 0,
    "DifficultyRating": null,
    "CreatedAt": "2026-03-30T10:30:00Z",
    "UpdatedAt": "2026-03-30T10:30:00Z"
  }
}

Errors:
- 400: Invalid week_id, empty front/back, or invalid request body
- 500: Server error
```

### 3. Get User Deck for Week
```
GET /api/v1/decks/weeks/{week_id}/cards

Response: 200 OK
{
  "data": [
    {
      "ID": "uuid",
      "UserID": "uuid",
      "WeekID": "uuid",
      "SourceFlashcardID": "uuid",
      "Front": "Question text",
      "Back": "Answer text",
      "IsCustom": false,
      "LastReviewedAt": "2026-03-30T10:00:00Z",
      "ReviewCount": 3,
      "DifficultyRating": 4,
      "CreatedAt": "2026-03-28T10:30:00Z",
      "UpdatedAt": "2026-03-30T10:00:00Z"
    },
    ...
  ]
}

Errors:
- 400: Invalid week_id
- 500: Server error
```

### 4. Update Card Content
```
PATCH /api/v1/decks/cards/{card_id}
Content-Type: application/json

Body:
{
  "front": "Updated question text",  // optional
  "back": "Updated answer text"      // optional
}

Response: 200 OK
{
  "data": null
}

Errors:
- 400: Invalid card_id or request body
- 404: Card not found or access denied
- 500: Server error
```

### 5. Remove Card from Deck
```
DELETE /api/v1/decks/cards/{card_id}

Response: 204 No Content

Errors:
- 400: Invalid card_id
- 404: Card not found or access denied
- 500: Server error
```

### 6. Record Card Review
```
POST /api/v1/decks/cards/{card_id}/review
Content-Type: application/json

Body:
{
  "difficulty_rating": 3  // 1-5 scale
}

Response: 200 OK
{
  "data": null
}

Errors:
- 400: Invalid card_id, rating not between 1-5, or invalid request body
- 404: Card not found or access denied
- 500: Server error
```

### 7. Get Deck Statistics
```
GET /api/v1/decks/weeks/{week_id}/stats

Response: 200 OK
{
  "data": {
    "total_cards": 25,
    "reviewed_cards": 12,
    "average_rating": 3.4,
    "last_reviewed_at": "2026-03-30T10:00:00Z"
  }
}

Errors:
- 400: Invalid week_id
- 500: Server error
```

## Security

### Authentication
- All endpoints require valid JWT token
- User ID extracted from token claims
- Middleware validates token before handler execution

### Authorization
- Users can only access/modify their own deck cards
- Repository methods filter by user_id
- Service layer validates card ownership before updates/deletes

### Data Validation
- UUID format validation for all IDs
- Non-empty validation for card front/back
- Difficulty rating range validation (1-5)
- SQL injection prevention via parameterized queries

## Database Migration

Migration files created:
- `migrations/000006_create_user_deck_cards.up.sql`
- `migrations/000006_create_user_deck_cards.down.sql`

To apply migration:
```bash
# Via docker-compose (automatic)
docker-compose up migrate

# Manual
migrate -path ./migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up
```

To rollback:
```bash
migrate -path ./migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" down 1
```

## Testing Workflow

### 1. Setup
```bash
# Start services
docker-compose up -d

# Verify migration ran
docker-compose logs migrate

# Check if table exists
docker exec -it <postgres-container> psql -U postgres -d studyhub -c "\dt user_deck_cards"
```

### 2. Get Authentication Token
```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'

# Save token
TOKEN="<jwt-token-from-response>"
```

### 3. Test Endpoints

**Create Custom Card:**
```bash
curl -X POST http://localhost:8080/api/v1/decks/weeks/<week-uuid>/cards/custom \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"front":"What is Go?","back":"A programming language"}'
```

**Get Deck:**
```bash
curl -X GET http://localhost:8080/api/v1/decks/weeks/<week-uuid>/cards \
  -H "Authorization: Bearer $TOKEN"
```

**Update Card:**
```bash
curl -X PATCH http://localhost:8080/api/v1/decks/cards/<card-uuid> \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"back":"Updated answer"}'
```

**Record Review:**
```bash
curl -X POST http://localhost:8080/api/v1/decks/cards/<card-uuid>/review \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"difficulty_rating":4}'
```

**Get Stats:**
```bash
curl -X GET http://localhost:8080/api/v1/decks/weeks/<week-uuid>/stats \
  -H "Authorization: Bearer $TOKEN"
```

**Delete Card:**
```bash
curl -X DELETE http://localhost:8080/api/v1/decks/cards/<card-uuid> \
  -H "Authorization: Bearer $TOKEN"
```

## Future Enhancements

### Not Yet Implemented
1. **Spaced Repetition Algorithm**
   - Implement SM-2 or similar algorithm
   - Calculate next review date based on difficulty
   - Sort cards by priority

2. **Bulk Operations**
   - Add multiple cards at once
   - Export/import deck

3. **Card Sharing**
   - Share decks between users
   - Public/private deck visibility

4. **Rich Content**
   - Images in flashcards
   - LaTeX math support
   - Code syntax highlighting

5. **Study Sessions**
   - Track session duration
   - Streak tracking
   - Daily goals

## Git Commits

Implementation completed in 7 commits:

1. `feat: add database migration for user deck cards`
2. `feat: add type definitions for user deck cards feature`
3. `feat: implement repository methods for user deck cards`
4. `feat: implement service layer for user deck cards`
5. `feat: add HTTP handlers for user deck cards endpoints`
6. `feat: register user deck card routes`

## Notes

- One deck per user per week (users cannot create multiple named decks)
- Auto-generated cards remain read-only in `flashcards` table
- Users get a copy when adding to their deck (can edit their copy)
- Duplicate prevention via unique constraint on (user_id, week_id, source_flashcard_id)
- Custom cards have `source_flashcard_id = NULL` and `is_custom = true`
- Review statistics are optional (can have cards without reviews)

## Contact

For questions or issues, refer to the main project documentation.
