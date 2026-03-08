-- Create lessons table
CREATE TABLE lessons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100) NOT NULL,
    duration_minutes INT NOT NULL,
    difficulty VARCHAR(20) CHECK (difficulty IN ('beginner', 'intermediate', 'advanced')),
    content TEXT NOT NULL,
    video_url TEXT,
    resources JSONB DEFAULT '[]',
    "order" INT NOT NULL
);

-- Create indexes for lessons
CREATE INDEX idx_lessons_category ON lessons(category, "order");
CREATE INDEX idx_lessons_difficulty ON lessons(difficulty);
CREATE INDEX idx_lessons_order ON lessons("order");

-- Add comments
COMMENT ON TABLE lessons IS 'Stores financial education lesson content';
COMMENT ON COLUMN lessons.difficulty IS 'Lesson difficulty level: beginner, intermediate, or advanced';
COMMENT ON COLUMN lessons.resources IS 'JSONB array of additional resources (links, documents, etc.)';
COMMENT ON COLUMN lessons."order" IS 'Display order of lessons within a category';
