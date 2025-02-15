package service

import (
	"database/sql"
	"fmt"
	"sync"

	"pengyou-chinese/backend/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

// DBService handles all database operations
type DBService struct {
	db *sql.DB
	mu sync.Mutex
}

// NewDBService creates a new database service instance
func NewDBService(dbPath string) (*DBService, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return &DBService{
		db: db,
	}, nil
}

// Close closes the database connection
func (s *DBService) Close() error {
	return s.db.Close()
}

// GetLastStudySession retrieves the most recent study session
func (s *DBService) GetLastStudySession() (*models.StudySession, error) {
	query := `
		SELECT s.id, s.group_id, s.created_at, s.study_activity_id, g.name as group_name
		FROM study_sessions s
		JOIN groups g ON s.group_id = g.id
		ORDER BY s.created_at DESC
		LIMIT 1
	`

	var session models.StudySession
	err := s.db.QueryRow(query).Scan(
		&session.ID,
		&session.GroupID,
		&session.CreatedAt,
		&session.StudyActivityID,
		&session.GroupName,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting last study session: %v", err)
	}

	return &session, nil
}

// GetStudyProgress retrieves study progress statistics
func (s *DBService) GetStudyProgress() (*models.StudyProgress, error) {
	query := `
		WITH studied_words AS (
			SELECT DISTINCT word_id
			FROM word_review_items
		)
		SELECT 
			(SELECT COUNT(*) FROM studied_words) as total_words_studied,
			(SELECT COUNT(*) FROM words) as total_available_words
	`

	var progress models.StudyProgress
	err := s.db.QueryRow(query).Scan(
		&progress.TotalWordsStudied,
		&progress.TotalAvailableWords,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting study progress: %v", err)
	}

	return &progress, nil
}

// GetQuickStats retrieves dashboard statistics
func (s *DBService) GetQuickStats() (*models.QuickStats, error) {
	query := `
		WITH review_stats AS (
			SELECT 
				COUNT(*) as total_reviews,
				SUM(CASE WHEN correct = 1 THEN 1 ELSE 0 END) as correct_reviews
			FROM word_review_items
		),
		active_groups AS (
			SELECT COUNT(DISTINCT group_id) as count
			FROM study_sessions
			WHERE created_at >= datetime('now', '-30 days')
		),
		streak_days AS (
			SELECT COUNT(DISTINCT date(created_at)) as count
			FROM study_sessions
			WHERE created_at >= datetime('now', '-30 days')
		)
		SELECT 
			COALESCE(CAST(correct_reviews AS FLOAT) / NULLIF(total_reviews, 0) * 100, 0) as success_rate,
			(SELECT COUNT(*) FROM study_sessions) as total_sessions,
			(SELECT count FROM active_groups) as active_groups,
			(SELECT count FROM streak_days) as streak_days
		FROM review_stats
	`

	var stats models.QuickStats
	err := s.db.QueryRow(query).Scan(
		&stats.SuccessRate,
		&stats.TotalStudySessions,
		&stats.TotalActiveGroups,
		&stats.StudyStreakDays,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting quick stats: %v", err)
	}

	return &stats, nil
}

// GetWords retrieves a paginated list of words with their statistics
func (s *DBService) GetWords(page, pageSize int) ([]models.WordWithStats, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var totalItems int
	countQuery := "SELECT COUNT(*) FROM words"
	if err := s.db.QueryRow(countQuery).Scan(&totalItems); err != nil {
		return nil, 0, fmt.Errorf("error counting words: %v", err)
	}

	// Get words with stats
	query := `
		SELECT 
			w.id, w.japanese, w.romaji, w.english, w.parts,
			COALESCE(SUM(CASE WHEN wr.correct = 1 THEN 1 ELSE 0 END), 0) as correct_count,
			COALESCE(SUM(CASE WHEN wr.correct = 0 THEN 1 ELSE 0 END), 0) as wrong_count
		FROM words w
		LEFT JOIN word_review_items wr ON w.id = wr.word_id
		GROUP BY w.id
		ORDER BY w.id
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying words: %v", err)
	}
	defer rows.Close()

	var words []models.WordWithStats
	for rows.Next() {
		var word models.WordWithStats
		err := rows.Scan(
			&word.ID,
			&word.Japanese,
			&word.Romaji,
			&word.English,
			&word.Parts,
			&word.CorrectCount,
			&word.WrongCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning word row: %v", err)
		}
		words = append(words, word)
	}

	return words, totalItems, nil
}

// AddWordReview adds a new word review record
func (s *DBService) AddWordReview(wordID, studySessionID int64, correct bool) error {
	query := `
		INSERT INTO word_review_items (word_id, study_session_id, correct)
		VALUES (?, ?, ?)
	`

	_, err := s.db.Exec(query, wordID, studySessionID, correct)
	if err != nil {
		return fmt.Errorf("error adding word review: %v", err)
	}

	return nil
}

// CreateStudySession creates a new study session
func (s *DBService) CreateStudySession(groupID, studyActivityID int64) (*models.StudySession, error) {
	query := `
		INSERT INTO study_sessions (group_id, study_activity_id)
		VALUES (?, ?)
		RETURNING id, group_id, created_at, study_activity_id
	`

	var session models.StudySession
	err := s.db.QueryRow(query, groupID, studyActivityID).Scan(
		&session.ID,
		&session.GroupID,
		&session.CreatedAt,
		&session.StudyActivityID,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating study session: %v", err)
	}

	return &session, nil
}

// GetWord retrieves a single word by ID with its statistics
func (s *DBService) GetWord(id int64) (*models.WordWithStats, error) {
	query := `
		SELECT 
			w.id, w.japanese, w.romaji, w.english, w.parts,
			COALESCE(SUM(CASE WHEN wr.correct = 1 THEN 1 ELSE 0 END), 0) as correct_count,
			COALESCE(SUM(CASE WHEN wr.correct = 0 THEN 1 ELSE 0 END), 0) as wrong_count
		FROM words w
		LEFT JOIN word_review_items wr ON w.id = wr.word_id
		WHERE w.id = ?
		GROUP BY w.id
	`

	var word models.WordWithStats
	err := s.db.QueryRow(query, id).Scan(
		&word.ID,
		&word.Japanese,
		&word.Romaji,
		&word.English,
		&word.Parts,
		&word.CorrectCount,
		&word.WrongCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting word: %v", err)
	}

	// Get groups for this word
	groupsQuery := `
		SELECT g.id, g.name
		FROM groups g
		JOIN words_groups wg ON g.id = wg.group_id
		WHERE wg.word_id = ?
	`

	rows, err := s.db.Query(groupsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("error getting word groups: %v", err)
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.ID, &group.Name); err != nil {
			return nil, fmt.Errorf("error scanning group: %v", err)
		}
		groups = append(groups, group)
	}

	return &word, nil
}

// GetGroups retrieves a paginated list of groups
func (s *DBService) GetGroups(page, pageSize int) ([]models.Group, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var totalItems int
	countQuery := "SELECT COUNT(*) FROM groups"
	if err := s.db.QueryRow(countQuery).Scan(&totalItems); err != nil {
		return nil, 0, fmt.Errorf("error counting groups: %v", err)
	}

	// Get groups with word count
	query := `
		SELECT 
			g.id, 
			g.name,
			COUNT(DISTINCT wg.word_id) as word_count
		FROM groups g
		LEFT JOIN words_groups wg ON g.id = wg.group_id
		GROUP BY g.id
		ORDER BY g.id
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying groups: %v", err)
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.WordCount); err != nil {
			return nil, 0, fmt.Errorf("error scanning group: %v", err)
		}
		groups = append(groups, group)
	}

	return groups, totalItems, nil
}

// GetGroup retrieves a single group by ID with statistics
func (s *DBService) GetGroup(id int64) (*models.Group, error) {
	query := `
		SELECT 
			g.id, 
			g.name,
			COUNT(DISTINCT wg.word_id) as word_count
		FROM groups g
		LEFT JOIN words_groups wg ON g.id = wg.group_id
		WHERE g.id = ?
		GROUP BY g.id
	`

	var group models.Group
	err := s.db.QueryRow(query, id).Scan(&group.ID, &group.Name, &group.WordCount)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting group: %v", err)
	}

	return &group, nil
}

// GetGroupWords retrieves words for a specific group
func (s *DBService) GetGroupWords(groupID int64, page, pageSize int) ([]models.WordWithStats, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var totalItems int
	countQuery := `
		SELECT COUNT(DISTINCT w.id)
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
	`
	if err := s.db.QueryRow(countQuery, groupID).Scan(&totalItems); err != nil {
		return nil, 0, fmt.Errorf("error counting group words: %v", err)
	}

	// Get words with stats
	query := `
		SELECT 
			w.id, w.japanese, w.romaji, w.english, w.parts,
			COALESCE(SUM(CASE WHEN wr.correct = 1 THEN 1 ELSE 0 END), 0) as correct_count,
			COALESCE(SUM(CASE WHEN wr.correct = 0 THEN 1 ELSE 0 END), 0) as wrong_count
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		LEFT JOIN word_review_items wr ON w.id = wr.word_id
		WHERE wg.group_id = ?
		GROUP BY w.id
		ORDER BY w.id
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.Query(query, groupID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying group words: %v", err)
	}
	defer rows.Close()

	var words []models.WordWithStats
	for rows.Next() {
		var word models.WordWithStats
		err := rows.Scan(
			&word.ID,
			&word.Japanese,
			&word.Romaji,
			&word.English,
			&word.Parts,
			&word.CorrectCount,
			&word.WrongCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning word row: %v", err)
		}
		words = append(words, word)
	}

	return words, totalItems, nil
}

// GetStudySessions retrieves a paginated list of study sessions
func (s *DBService) GetStudySessions(page, pageSize int) ([]models.StudySession, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var totalItems int
	countQuery := "SELECT COUNT(*) FROM study_sessions"
	if err := s.db.QueryRow(countQuery).Scan(&totalItems); err != nil {
		return nil, 0, fmt.Errorf("error counting study sessions: %v", err)
	}

	// Get study sessions with group names
	query := `
		SELECT 
			s.id, s.group_id, s.created_at, s.study_activity_id,
			g.name as group_name,
			COUNT(wr.id) as review_items_count
		FROM study_sessions s
		JOIN groups g ON s.group_id = g.id
		LEFT JOIN word_review_items wr ON s.id = wr.study_session_id
		GROUP BY s.id
		ORDER BY s.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying study sessions: %v", err)
	}
	defer rows.Close()

	var sessions []models.StudySession
	for rows.Next() {
		var session models.StudySession
		var reviewCount int
		err := rows.Scan(
			&session.ID,
			&session.GroupID,
			&session.CreatedAt,
			&session.StudyActivityID,
			&session.GroupName,
			&reviewCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning study session row: %v", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, totalItems, nil
}

// GetStudySession retrieves a single study session by ID
func (s *DBService) GetStudySession(id int64) (*models.StudySession, error) {
	query := `
		SELECT 
			s.id, s.group_id, s.created_at, s.study_activity_id,
			g.name as group_name,
			COUNT(wr.id) as review_items_count
		FROM study_sessions s
		JOIN groups g ON s.group_id = g.id
		LEFT JOIN word_review_items wr ON s.id = wr.study_session_id
		WHERE s.id = ?
		GROUP BY s.id
	`

	var session models.StudySession
	var reviewCount int
	err := s.db.QueryRow(query, id).Scan(
		&session.ID,
		&session.GroupID,
		&session.CreatedAt,
		&session.StudyActivityID,
		&session.GroupName,
		&reviewCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting study session: %v", err)
	}

	return &session, nil
}

// GetStudySessionWords retrieves words reviewed in a study session
func (s *DBService) GetStudySessionWords(sessionID int64, page, pageSize int) ([]models.WordWithStats, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var totalItems int
	countQuery := `
		SELECT COUNT(DISTINCT w.id)
		FROM words w
		JOIN word_review_items wr ON w.id = wr.word_id
		WHERE wr.study_session_id = ?
	`
	if err := s.db.QueryRow(countQuery, sessionID).Scan(&totalItems); err != nil {
		return nil, 0, fmt.Errorf("error counting session words: %v", err)
	}

	// Get words with their review status for this session
	query := `
		SELECT 
			w.id, w.japanese, w.romaji, w.english, w.parts,
			SUM(CASE WHEN wr2.correct = 1 THEN 1 ELSE 0 END) as correct_count,
			SUM(CASE WHEN wr2.correct = 0 THEN 1 ELSE 0 END) as wrong_count,
			wr1.correct as session_correct
		FROM words w
		JOIN word_review_items wr1 ON w.id = wr1.word_id AND wr1.study_session_id = ?
		LEFT JOIN word_review_items wr2 ON w.id = wr2.word_id
		GROUP BY w.id
		ORDER BY wr1.created_at
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.Query(query, sessionID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying session words: %v", err)
	}
	defer rows.Close()

	var words []models.WordWithStats
	for rows.Next() {
		var word models.WordWithStats
		var sessionCorrect bool
		err := rows.Scan(
			&word.ID,
			&word.Japanese,
			&word.Romaji,
			&word.English,
			&word.Parts,
			&word.CorrectCount,
			&word.WrongCount,
			&sessionCorrect,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning word row: %v", err)
		}
		words = append(words, word)
	}

	return words, totalItems, nil
}

// GetStudyActivity retrieves a study activity by ID
func (s *DBService) GetStudyActivity(id int64) (*models.StudyActivity, error) {
	query := `
		SELECT id, study_session_id, group_id, created_at
		FROM study_activities
		WHERE id = ?
	`

	var activity models.StudyActivity
	err := s.db.QueryRow(query, id).Scan(
		&activity.ID,
		&activity.StudySessionID,
		&activity.GroupID,
		&activity.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting study activity: %v", err)
	}

	return &activity, nil
}

// GetStudyActivitySessions retrieves study sessions for a specific activity
func (s *DBService) GetStudyActivitySessions(activityID int64, page, pageSize int) ([]models.StudySession, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var totalItems int
	countQuery := `
		SELECT COUNT(*)
		FROM study_sessions
		WHERE study_activity_id = ?
	`
	if err := s.db.QueryRow(countQuery, activityID).Scan(&totalItems); err != nil {
		return nil, 0, fmt.Errorf("error counting activity sessions: %v", err)
	}

	// Get sessions with group names
	query := `
		SELECT 
			s.id, s.group_id, s.created_at, s.study_activity_id,
			g.name as group_name,
			COUNT(wr.id) as review_items_count
		FROM study_sessions s
		JOIN groups g ON s.group_id = g.id
		LEFT JOIN word_review_items wr ON s.id = wr.study_session_id
		WHERE s.study_activity_id = ?
		GROUP BY s.id
		ORDER BY s.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.Query(query, activityID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying activity sessions: %v", err)
	}
	defer rows.Close()

	var sessions []models.StudySession
	for rows.Next() {
		var session models.StudySession
		var reviewCount int
		err := rows.Scan(
			&session.ID,
			&session.GroupID,
			&session.CreatedAt,
			&session.StudyActivityID,
			&session.GroupName,
			&reviewCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning session row: %v", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, totalItems, nil
}
