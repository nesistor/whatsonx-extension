
# Google Calendar Meeting Scheduler API for WatsonX

This project provides a RESTful API that integrates with Google Calendar to enable meeting scheduling and calendar management for WatsonX. With this API, WatsonX can offer functionalities like checking availability across calendars, proposing meetings, sending invitations with links, and managing user groups for efficient scheduling.


## Features

- **User Authorization**: Users authorize the application through Google OAuth2.
- **Check Availability**: Query available slots in a user’s calendar for meeting proposals.
- **Propose Meetings**: Based on available slots, propose and schedule meetings.
- **Group Management**: Add users to specific groups for team-based scheduling.
- **Invitation and Meeting Scheduling**: Send Google Meet invitations and automatically add events to participants' calendars.

## Endpoints Overview

| Endpoint                 | Method | Description                                  |
|--------------------------|--------|----------------------------------------------|
| `/add-user`              | `POST` | Initiates user authorization process.        |
| `/oauth2callback`        | `GET`  | Handles Google OAuth2 callback.             |
| `/check-availability`    | `GET`  | Retrieves availability data for a user.     |
| `/add-user-to-group`     | `POST` | Adds a user to a specific group.            |
| `/list-users`            | `GET`  | Lists all registered users.                 |
| `/list-groups`           | `GET`  | Lists all available groups.                 |

## How It Works

1.User Authorization
Users authorize the application via Google OAuth2. The API stores the access and refresh tokens securely for calendar operations.

Check Availability
Use the /check-availability endpoint to query free slots in a user’s calendar within a specified time range.

Propose Meetings
Based on availability, WatsonX can propose a meeting time and use this API to schedule the meeting and send invites.

Group Management
Users can be grouped together using /add-user-to-group, enabling meeting proposals for entire teams or departments.

Invitation and Meeting Scheduling
Automatically send Google Meet invitations and add the event to participants’ calendars.
## Instalation

1. Clone the repository:

   ```bash
   git clone git@github.com:nesistor/whatsonx-meeting-scheduler.git
   cd whatsonx-meeting-scheduler