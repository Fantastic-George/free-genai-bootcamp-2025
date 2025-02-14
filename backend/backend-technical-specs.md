# Backend server technical specs

## Business Goal:

A language learning school wants to build a prototype of a learning portal which will act as three things:
- Inventory of possible vocabulary that can be learned
- Act as a learning record score (LRS), providing correct and wrong score on practice vocabulary
- A unified launchpad to launch different learning apps

## Technical requirements:

- The backend will be built using Go
- The database will be SQLite3
- The API will be built using Gin
- The API will always return JSON
- There will be no authentication

## Database schema:
- words - stored vocabulary words
    - id (integer)
    - Chinese character (string)
    - Pinyin (string)
    - English translation (string)
    - Radicals (string)
    - Parts JSON (json)
- word_groups - join table for words and groups many-to-many
    - id (integer)
    - word_id (integer)
    - group_id (integer)
- groups - thematic groups of words
    - id (integer)
    - name (string)
    - description (string)
- study_sessions - record of a study session grouping word_review_items
    - id (integer)
    - group_id (integer)
    - created_at (datetime)
    - study_activity_id (integer)
- study_activities - a specific study activity, linking a study session to group
    - id (integer)
    - name (string)
    - description (string)
    - created_at (datetime)
    - updated_at (datetime)
    - group_id (integer)
- word_review_items -a record of word practice, determining if the word was correct or not
    - word_id (integer)
    - study_session_id (integer)
    - correct (boolean)
    - created_at (datetime)

## API Endpoints:

- GET /words
- GET /words/:id
- GET /groups
    - pagination with 100 items per page
- GET /groups/:id
- GET /study_sessions
- GET /study_activities
- GET /word_review_items
- GET /api/study_activities/:id
- GET /api/study_activities/:id/study_sessions
- GET /api/dashboard/last_study_session
- GET /api/dashboard/study_progress
- GET /api/dashboard/quick_stats
- GET /api/words
    - pagination with 100 items per page
- GET /api/groups/:id/study_sessions
- GET /api/study_sessions
    - pagination with 100 items per page
- GET /api/study_sessions/:id
- GET /api/study_sessions/:id/words


- POST /api/study_activities/:id/launch
    - params:
        - group_id (integer)
        - study_activity_id (integer)
- POST /api/dashboard/start_studying
- POST /api/settings/theme
- POST /api/settings/reset_history
- POST /api/settings/full_reset
- POST /api/study_sessions/:id/words/:word_id/review

### Get all words:

