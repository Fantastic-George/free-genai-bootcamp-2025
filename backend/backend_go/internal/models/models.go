package models

import "time"

// Word represents a vocabulary word
type Word struct {
	ID       int64  `json:"id"`
	Japanese string `json:"japanese"`
	Romaji   string `json:"romaji"`
	English  string `json:"english"`
	Parts    string `json:"parts,omitempty"` // JSON string
}

// Group represents a thematic group of words
type Group struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	WordCount int    `json:"word_count,omitempty"`
}

// StudySession represents a learning session
type StudySession struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	CreatedAt       time.Time `json:"created_at"`
	StudyActivityID int64     `json:"study_activity_id"`
	GroupName       string    `json:"group_name,omitempty"`
}

// StudyActivity represents a specific study activity
type StudyActivity struct {
	ID              int64     `json:"id"`
	StudySessionID  int64     `json:"study_session_id"`
	GroupID         int64     `json:"group_id"`
	CreatedAt       time.Time `json:"created_at"`
}

// WordReviewItem represents a practice record for a word
type WordReviewItem struct {
	ID             int64     `json:"id"`
	WordID         int64     `json:"word_id"`
	StudySessionID int64     `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
}

// WordWithStats extends Word with statistics
type WordWithStats struct {
	Word
	CorrectCount int `json:"correct_count"`
	WrongCount   int `json:"wrong_count"`
}

// StudyProgress represents study progress statistics
type StudyProgress struct {
	TotalWordsStudied    int `json:"total_words_studied"`
	TotalAvailableWords  int `json:"total_available_words"`
}

// QuickStats represents dashboard statistics
type QuickStats struct {
	SuccessRate       float64 `json:"success_rate"`
	TotalStudySessions int    `json:"total_study_sessions"`
	TotalActiveGroups  int    `json:"total_active_groups"`
	StudyStreakDays    int    `json:"study_streak_days"`
} 