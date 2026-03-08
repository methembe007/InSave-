-- Create education progress table
CREATE TABLE education_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    is_completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(user_id, lesson_id)
);

-- Create indexes for education progress
CREATE INDEX idx_education_progress_user ON education_progress(user_id, is_completed);
CREATE INDEX idx_education_progress_lesson ON education_progress(lesson_id);
CREATE INDEX idx_education_progress_completed ON education_progress(user_id, completed_at DESC);

-- Add comments
COMMENT ON TABLE education_progress IS 'Tracks user progress through financial education lessons';
COMMENT ON COLUMN education_progress.is_completed IS 'Whether the user has completed the lesson';
COMMENT ON COLUMN education_progress.completed_at IS 'Timestamp when the lesson was completed';
