# InsiderBackendStudy

League Simulation API Project for Insider's internship program.

# Using Public Server
You can try the API without setting up an environment yourself.

The server will be up from 25/05/2025 to 24/06/2025, and is running on a Compute Engine VM Instance of Google Cloud.

There is no authorization.

Public IP: `35.239.226.80`
Port: `8080`

Refer to the endpoints for testing, and use the given public IP as the hostname instead of localhost.

# Setting Up
1- Prepare the .env file according to the given .env.example

2- Run main.go

3- Refer to the hostname accordingly.


# Endpoints

| Method | Path                   | Description                                                               |
| ------ | ---------------------- | ------------------------------------------------------------------------- |
| POST   | `/reset`               | Reset the simulation: clears all played matches and resets the week to 1. |
| POST   | `/simulate_week`       | Simulate the next unplayed week’s matches.                                |
| POST   | `/simulate_tournament` | Simulate all remaining weeks until the tournament is complete.            |
| GET    | `/teams`               | Retrieve the list of all teams with their current stats.                  |
| GET    | `/all_matches`         | Retrieve every match in the database (played and unplayed).               |
| GET    | `/matches/:week`       | Retrieve matches scheduled for a specific week.                           |
| GET    | `/last_matches`        | Retrieve matches played in the most recently simulated week.              |
| GET    | `/next_matches`        | Retrieve matches scheduled for the upcoming week.                         |
| GET    | `/predict_chances`     | Estimate winning probabilities for all remaining matches.                 |

---

### 1. POST `/reset`

Resets the entire simulation back to week 1 and clears all match results and team statistics.

**Request:**

```http
POST /reset HTTP/1.1
Host: localhost:8080
```

**Response:**

* **Status:** `200 OK`
* **Body:**

  ```json
  {
    "status": "ok",
    "message": "Simulation has been reset. Played matches are deleted and week is reset to week 1."
  }
  ```

---

### 2. POST `/simulate_week`

Simulates matches for the next unplayed week in sequence.

**Request:**

```http
POST /simulate_week HTTP/1.1
Host: localhost:8080
```

**Responses:**

* **Status:** `200 OK`
* First simulation:

  ```json
  { "status": "ok", "message": "Simulation for week 2 is done." }
  ```
* Subsequent weeks:

  ```json
  { "status": "ok", "message": "Simulation for week 3 is done." }
  ```
* After completion:

  ```json
  { "status": "ok", "message": "Simulation for tournament is already done. Reset to simulate again." }
  ```

---

### 3. POST `/simulate_tournament`

Simulates all remaining weeks until the tournament is complete.

**Request:**

```http
POST /simulate_tournament HTTP/1.1
Host: localhost:8080
```

**Responses:**

* **Status:** `200 OK`
* On first completion:

  ```json
  { "status": "ok", "message": "Simulation for tournament is done." }
  ```
* If already complete:

  ```json
  { "status": "ok", "message": "Simulation for tournament is already done. Reset to simulate again." }
  ```

---

### 4. GET `/teams`

Fetch the current list of teams and their statistics.

**Request:**

```http
GET /teams HTTP/1.1
Host: localhost:8080
```

**Response:**

* **Status:** `200 OK`
* **Body:**

  ```json
  [
    {
      "name": "Team1",
      "points": 6,
      "strength": 100,
      "played": 3,
      "won": 2,
      "drawn": 0,
      "lost": 1,
      "goalsFor": 5,
      "goalsAgainst": 3
    },
    { /* ... */ }
  ]
  ```

---

### 5. GET `/all_matches`

Fetch every match record, both played and unplayed.

**Request:**

```http
GET /all_matches HTTP/1.1
Host: localhost:8080
```

**Response:**

* **Status:** `200 OK`
* **Body:**

  ```json
  [
    {
      "homeTeam": <Team1>,
      "awayTeam": <Team2>,
      "isPlayed": True,
      "homeGoals": 2,
      "awayGoals": 1,
      "week": 1,
    },
    { /* ... */ }
  ]
  ```

---

### 6. GET `/matches/:week`

Fetch matches scheduled for a specific week.

**Request:**

```http
GET /matches/2 HTTP/1.1
Host: localhost:8080
```

**Response:**

* **Status:** `200 OK`
* **Body:**

  ```json
  [
    {
      "homeTeam": <Team1>,
      "awayTeam": <Team2>,
      "isPlayed": True,
      "homeGoals": 2,
      "awayGoals": 1,
      "week": 1,
    },
    { /* ... */ }
  ]
  ```

---

### 7. GET `/last_matches`

Fetch matches from the most recently simulated week.

**Request:**

```http
GET /last_matches HTTP/1.1
Host: localhost:8080
```

**Responses:**

* **If at least one week played (`200 OK`):**

  ```json
  [
    {
      "homeTeam": <Team1>,
      "awayTeam": <Team2>,
      "isPlayed": True,
      "homeGoals": 2,
      "awayGoals": 1,
      "week": 1,
    },
    { /* ... */ }
  ]
  ```
* **If no week played (`400 Bad Request`):**

  ```json
  { "status": "error", "message": "No matches have been played yet." }
  ```

---

### 8. GET `/next_matches`

Fetch matches scheduled for the upcoming week.

**Request:**

```http
GET /next_matches HTTP/1.1
Host: localhost:8080
```

**Responses:**

* **If upcoming matches exist (`200 OK`):** array of match objects (same shape as above).
* **If tournament complete (`200 OK`):**

  ```json
  { "status": "ok", "message": "Simulation for tournament is already done. Reset to simulate again." }
  ```

---

### 9. GET `/predict_chances`

Estimate winning probabilities for remaining matches.

**Request:**

```http
GET /predict_chances HTTP/1.1
Host: localhost:8080
```

**Responses:**

* **If unplayed matches remain (`200 OK`):**

  ```json
  [
    {
      "Team1": 26,124,
      "Team2": 21,008,
      "Team3": 29,27,
      "Team4": 23,598,
    },
    { /* ... */ }
  ]
  ```
* **If tournament complete (`200 OK`):**

  ```json
  { "status": "ok", "message": "Simulation for tournament is already done. Reset to simulate again." }
  ```

---

## Error Handling & Notes

* **500 Internal Server Error** on unexpected failures (database or logic).
* **Idempotent operations:**

  * `/reset` always returns to the initial state.
  * `/simulate_tournament` has no effect once complete.

Ensure the `"server"` remains running to preserve state between calls.
