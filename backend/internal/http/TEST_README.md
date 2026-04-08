# Deck API Tests

This file contains unit tests for the deck-related HTTP handlers.

## Test Coverage

The test suite covers all 7 deck API endpoints:

### 1. AddCardToDeckHandler
Tests for adding auto-generated flashcards to user's deck.

**Test Cases:**
- Success: Add card to deck
- Error: Invalid week ID
- Error: Invalid flashcard ID
- Error: Flashcard not found
- Error: Duplicate card (already in deck)

### 2. CreateCustomCardHandler
Tests for creating custom flashcards.

**Test Cases:**
- Success: Create custom card with front and back
- Error: Empty front field
- Error: Empty back field
- Error: Invalid week ID

### 3. GetUserDeckHandler
Tests for retrieving user's deck for a specific week.

**Test Cases:**
- Success: Get deck with multiple cards
- Success: Get empty deck
- Error: Invalid week ID

### 4. UpdateDeckCardHandler
Tests for updating existing deck cards.

**Test Cases:**
- Success: Update both front and back
- Success: Update front only
- Success: Update back only
- Error: Card not found
- Error: Invalid card ID

### 5. RemoveDeckCardHandler
Tests for removing cards from deck.

**Test Cases:**
- Success: Remove card
- Error: Card not found
- Error: Invalid card ID

### 6. RecordCardReviewHandler
Tests for recording review sessions with difficulty ratings.

**Test Cases:**
- Success: Record review with valid difficulty (3)
- Error: Invalid difficulty rating (0)
- Error: Invalid difficulty rating (6)

### 7. GetDeckStatsHandler
Tests for retrieving deck statistics.

**Test Cases:**
- Success: Get deck stats with data
- Success: Get empty deck stats
- Error: Invalid week ID

## Running Tests

Run all tests:
```bash
cd backend
go test ./internal/http -v
```

Run specific test:
```bash
go test ./internal/http -v -run TestAddCardToDeckHandler
```

Run with coverage:
```bash
go test ./internal/http -v -cover
```

## Test Structure

Each test uses a mock service that implements the necessary methods without requiring a database connection. This allows for:
- Fast test execution
- Isolated testing of handler logic
- Easy simulation of error conditions
- No external dependencies

## Mock Service

The `mockContentService` type implements all deck-related service methods with configurable behavior through function fields. This allows each test to define custom responses and errors.

Example:
```go
mockSvc := &mockContentService{
    addCardToDeckFunc: func(ctx context.Context, userID, weekID, flashcardID uuid.UUID) error {
        return errors.New("flashcard not found")
    },
}
```
