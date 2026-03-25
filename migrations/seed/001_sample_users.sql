-- Sample Users Seed Data
-- This script creates sample users for testing and development

-- Insert sample users with bcrypt hashed passwords (password: "password123")
-- Bcrypt hash generated with cost factor 12
INSERT INTO users (id, email, password_hash, first_name, last_name, date_of_birth, profile_image_url, preferences, created_at, updated_at)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'john.doe@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIr.8VfGHO', 'John', 'Doe', '1995-03-15', 'https://i.pravatar.cc/150?img=12', '{"currency": "USD", "notifications_enabled": true, "email_notifications": true, "push_notifications": true, "savings_reminders": true, "reminder_time": "09:00", "theme": "light"}', NOW() - INTERVAL '90 days', NOW()),
    ('22222222-2222-2222-2222-222222222222', 'jane.smith@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIr.8VfGHO', 'Jane', 'Smith', '1998-07-22', 'https://i.pravatar.cc/150?img=45', '{"currency": "USD", "notifications_enabled": true, "email_notifications": true, "push_notifications": false, "savings_reminders": true, "reminder_time": "08:00", "theme": "dark"}', NOW() - INTERVAL '60 days', NOW()),
    ('33333333-3333-3333-3333-333333333333', 'mike.johnson@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIr.8VfGHO', 'Mike', 'Johnson', '1997-11-08', 'https://i.pravatar.cc/150?img=33', '{"currency": "USD", "notifications_enabled": false, "email_notifications": false, "push_notifications": false, "savings_reminders": false, "reminder_time": "10:00", "theme": "light"}', NOW() - INTERVAL '45 days', NOW()),
    ('44444444-4444-4444-4444-444444444444', 'sarah.williams@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIr.8VfGHO', 'Sarah', 'Williams', '1996-05-30', 'https://i.pravatar.cc/150?img=47', '{"currency": "USD", "notifications_enabled": true, "email_notifications": true, "push_notifications": true, "savings_reminders": true, "reminder_time": "07:30", "theme": "dark"}', NOW() - INTERVAL '30 days', NOW()),
    ('55555555-5555-5555-5555-555555555555', 'alex.brown@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIr.8VfGHO', 'Alex', 'Brown', '1999-01-12', 'https://i.pravatar.cc/150?img=68', '{"currency": "USD", "notifications_enabled": true, "email_notifications": false, "push_notifications": true, "savings_reminders": true, "reminder_time": "09:30", "theme": "light"}', NOW() - INTERVAL '15 days', NOW());

-- Note: All users have the same password for testing: "password123"
-- In production, users would set their own unique passwords during registration
