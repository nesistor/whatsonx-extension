
# Google Calendar Meeting Scheduler API for WatsonX

This project provides a RESTful API that integrates with Google Calendar to enable meeting scheduling and calendar management for WatsonX. With this API, WatsonX can offer functionalities like checking availability across calendars, proposing meetings, sending invitations with links, and managing user groups for efficient scheduling.


## Features

- **User Authorization**: Users authorize the application through Google OAuth2.
- **Check Availability**: Query available slots in a user’s calendar for meeting proposals.
- **Propose Meetings**: Based on available slots, propose and schedule meetings.
- **Group Management**: Add users to specific groups for team-based scheduling.
- **Invitation and Meeting Scheduling**: Send Google Meet invitations and automatically add events to participants' calendars.

## API Documentation

The API documentation is automatically generated and can be viewed through the Swagger UI. To access the documentation, start the application and visit:

[Swagger UI](http://localhost:8080/swagger)

## API Endpoints

| Endpoint                 | Method | Description                                  |
|--------------------------|--------|----------------------------------------------|
| `/add-user`              | `POST` | Initiates user authorization process.        |
| `/oauth2callback`        | `GET`  | Handles Google OAuth2 callback.             |
| `/check-availability`    | `GET`  | Retrieves availability data for a user.     |
| `/add-user-to-group`     | `POST` | Adds a user to a specific group.            |
| `/list-users`            | `GET`  | Lists all registered users.                 |
| `/list-groups`           | `GET`  | Lists all available groups.                 |
| `/swagger/*`             | `GET`  | View the Swagger documentation.             |

## How It Works

### 1. User Authorization
Users authorize the application via **Google OAuth2**. The API securely stores the access and refresh tokens for calendar operations.

### 2. Check Availability
Use the `/check-availability` endpoint to query available time slots in a user's calendar within a specified time range.

### 3. Propose Meetings
Based on the available time slots, **WatsonX** can propose a meeting time and use the API to schedule the meeting, automatically sending invites to participants.

### 4. Group Management
Users can be grouped together using the `/add-user-to-group` endpoint, enabling meeting proposals for entire teams or departments.

### 5. Invitation and Meeting Scheduling
The system automatically sends **Google Meet** invitations and adds the scheduled event to participants’ calendars.
## Installation

Follow these steps to install and set up **Meeting Scheduler** on your local machine.

### 1. Clone the repository

Start by cloning the repository to your local machine:

```bash
git clone git@github.com:nesistor/whatsonx-meeting-scheduler.git
