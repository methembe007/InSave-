# Notification Service API Examples

This document provides examples of how to interact with the Notification Service API.

## Prerequisites

- Service running on `http://localhost:8086`
- Valid JWT token from Auth Service
- User ID in token claims

## Authentication

All API requests require a valid JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## API Endpoints

### 1. Get User Notifications

Retrieve all notifications for the authenticated user, ordered by creation date (newest first).

**Request:**
```bash
curl -X GET http://localhost:8086/api/notifications \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response:** `200 OK`
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "type": "push",
    "title": "Savings Reminder",
    "message": "Don't forget to save today! Keep your streak going.",
    "is_read": false,
    "created_at": "2024-01-15T10:30:00Z"
  },
  {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "type": "budget_alert",
    "title": "Budget Alert",
    "message": "You've used 85% of your Food budget this month.",
    "is_read": true,
    "created_at": "2024-01-14T15:20:00Z"
  },
  {
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "type": "goal_milestone",
    "title": "Goal Milestone Reached",
    "message": "Congratulations! You've reached 50% of your Emergency Fund goal.",
    "is_read": true,
    "created_at": "2024-01-13T09:15:00Z"
  }
]
```

**Empty Response:** `200 OK`
```json
[]
```

**Error Response:** `401 Unauthorized`
```json
{
  "error": "Authorization header required"
}
```

### 2. Mark Notification as Read

Mark a specific notification as read for the authenticated user.

**Request:**
```bash
curl -X PUT http://localhost:8086/api/notifications/550e8400-e29b-41d4-a716-446655440000/read \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response:** `200 OK`
```json
{
  "message": "Notification marked as read"
}
```

**Error Response:** `404 Not Found`
```json
{
  "error": "Notification not found or does not belong to user"
}
```

**Error Response:** `401 Unauthorized`
```json
{
  "error": "Invalid or expired token"
}
```

## Notification Types

The service supports various notification types:

- `push`: General push notifications
- `email`: Email notifications
- `savings_reminder`: Reminder to save money
- `budget_reminder`: Budget-related reminders
- `goal_reminder`: Goal-related reminders
- `budget_alert`: Budget threshold alerts
- `goal_milestone`: Goal milestone achievements
- `streak_achievement`: Savings streak milestones

## Integration Examples

### JavaScript/TypeScript (Frontend)

```typescript
// API client configuration
const API_BASE_URL = 'http://localhost:8086/api';

// Get notifications
async function getNotifications(token: string) {
  const response = await fetch(`${API_BASE_URL}/notifications`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });
  
  if (!response.ok) {
    throw new Error('Failed to fetch notifications');
  }
  
  return await response.json();
}

// Mark notification as read
async function markNotificationAsRead(token: string, notificationId: string) {
  const response = await fetch(`${API_BASE_URL}/notifications/${notificationId}/read`, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });
  
  if (!response.ok) {
    throw new Error('Failed to mark notification as read');
  }
  
  return await response.json();
}

// Usage example
const token = 'your-jwt-token';

// Get all notifications
const notifications = await getNotifications(token);
console.log('Notifications:', notifications);

// Mark first unread notification as read
const unreadNotification = notifications.find(n => !n.is_read);
if (unreadNotification) {
  await markNotificationAsRead(token, unreadNotification.id);
  console.log('Notification marked as read');
}
```

### React Hook Example

```typescript
import { useState, useEffect } from 'react';

interface Notification {
  id: string;
  user_id: string;
  type: string;
  title: string;
  message: string;
  is_read: boolean;
  created_at: string;
}

export function useNotifications(token: string) {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchNotifications();
  }, [token]);

  const fetchNotifications = async () => {
    try {
      setLoading(true);
      const response = await fetch('http://localhost:8086/api/notifications', {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      
      if (!response.ok) {
        throw new Error('Failed to fetch notifications');
      }
      
      const data = await response.json();
      setNotifications(data);
      setError(null);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const markAsRead = async (notificationId: string) => {
    try {
      const response = await fetch(
        `http://localhost:8086/api/notifications/${notificationId}/read`,
        {
          method: 'PUT',
          headers: {
            'Authorization': `Bearer ${token}`,
          },
        }
      );
      
      if (!response.ok) {
        throw new Error('Failed to mark notification as read');
      }
      
      // Update local state
      setNotifications(prev =>
        prev.map(n =>
          n.id === notificationId ? { ...n, is_read: true } : n
        )
      );
    } catch (err) {
      setError(err.message);
    }
  };

  const unreadCount = notifications.filter(n => !n.is_read).length;

  return {
    notifications,
    loading,
    error,
    unreadCount,
    markAsRead,
    refresh: fetchNotifications,
  };
}
```

### Go Client Example

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type Notification struct {
    ID        string `json:"id"`
    UserID    string `json:"user_id"`
    Type      string `json:"type"`
    Title     string `json:"title"`
    Message   string `json:"message"`
    IsRead    bool   `json:"is_read"`
    CreatedAt string `json:"created_at"`
}

func getNotifications(token string) ([]Notification, error) {
    req, err := http.NewRequest("GET", "http://localhost:8086/api/notifications", nil)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Authorization", "Bearer "+token)
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
    }
    
    var notifications []Notification
    if err := json.NewDecoder(resp.Body).Decode(&notifications); err != nil {
        return nil, err
    }
    
    return notifications, nil
}

func markNotificationAsRead(token, notificationID string) error {
    url := fmt.Sprintf("http://localhost:8086/api/notifications/%s/read", notificationID)
    req, err := http.NewRequest("PUT", url, nil)
    if err != nil {
        return err
    }
    
    req.Header.Set("Authorization", "Bearer "+token)
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status: %d", resp.StatusCode)
    }
    
    return nil
}
```

## Testing with curl

### Get all notifications
```bash
export TOKEN="your-jwt-token-here"

curl -X GET http://localhost:8086/api/notifications \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"
```

### Mark notification as read
```bash
export TOKEN="your-jwt-token-here"
export NOTIFICATION_ID="550e8400-e29b-41d4-a716-446655440000"

curl -X PUT http://localhost:8086/api/notifications/$NOTIFICATION_ID/read \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"
```

### Health check (no authentication required)
```bash
curl -X GET http://localhost:8086/health
```

## Error Handling

The API returns standard HTTP status codes:

- `200 OK`: Request successful
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Missing or invalid authentication token
- `403 Forbidden`: User not authorized to access resource
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

Error responses include a descriptive message:
```json
{
  "error": "Description of the error"
}
```

## Rate Limiting

Currently, no rate limiting is implemented. In production, consider:
- Rate limiting per user (e.g., 100 requests per minute)
- Rate limiting per IP address
- Implementing exponential backoff for retries

## Best Practices

1. **Cache notifications**: Cache notification data on the client to reduce API calls
2. **Poll periodically**: Poll for new notifications every 30-60 seconds
3. **Use WebSockets**: For real-time notifications, consider implementing WebSocket support
4. **Handle errors gracefully**: Always handle network errors and display user-friendly messages
5. **Optimize queries**: Use pagination for large notification lists (to be implemented)
6. **Mark as read on view**: Automatically mark notifications as read when user views them
