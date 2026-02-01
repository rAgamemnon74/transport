# Transport CLI - Design Document

## Overview

A command-line transport planner for Swedish public transit, focused on Stockholm and Sweden-wide travel.

## Usage

```bash
# Trip planning - one parameter: destination only (origin = default)
transport "T-Centralen"

# Trip planning - two parameters: origin and destination
transport "Slussen" "Odenplan"

# Trip planning - with time option
transport -t 08:30 "Slussen" "Odenplan"

# Next departures - show next buses towards a destination
transport next bus "Styresman Sanders väg" Brommaplan

# Next departures - show next metro towards end station
transport next metro Slussen Ropsten

# Next departures - show more results
transport next -n 5 bus Odenplan Solna
```

## API

### SL Journey Planner v2 (Stockholm)
- **Base URL:** `https://journeyplanner.integration.sl.se/v2/`
- **No API key required**
- **Endpoints:**
  - `/stop-finder` - Search for stops/locations
  - `/trips` - Get journey proposals
  - `/system-info` - API status

## Output Format (inspired by SL Reseplanerare)

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 Slussen → Odenplan                                    Fri 31 Jan
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

 Trip 1                                          14 min │ 0 byten
───────────────────────────────────────────────────────────────────
  08:15  Slussen                              Spår 1
    │    🚇 Tunnelbana 17 → Odenplan
  08:29  Odenplan                             Spår 2

 Trip 2                                          18 min │ 1 byte
───────────────────────────────────────────────────────────────────
  08:20  Slussen                              Spår 3
    │    🚇 Tunnelbana 18 → Alvik
  08:27  T-Centralen
    │    ⏱  Byte 3 min
  08:30  T-Centralen                          Spår 1
    │    🚇 Tunnelbana 17 → Åkeshov
  08:38  Odenplan                             Spår 2

 Trip 3                                          22 min │ 0 byten
───────────────────────────────────────────────────────────────────
  08:25  Slussen                              Hållplats A
    │    🚌 Buss 3 → Södersjukhuset
  08:47  Odenplan                             Hållplats B

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## Transport Type Indicators

| Symbol | Swedish         | English       | API Code |
|--------|-----------------|---------------|----------|
| 🚇     | Tunnelbana      | Metro         | METRO    |
| 🚂     | Pendeltåg       | Commuter Rail | TRAIN    |
| 🚆     | Tåg (SJ etc)    | Train         | TRAIN    |
| 🚊     | Spårvagn        | Tram          | TRAM     |
| 🚌     | Buss            | Bus           | BUS      |
| ⛴️     | Båt/Färja       | Ferry         | FERRY    |
| 🚶     | Gång            | Walk          | WALK     |
| 🚕     | Närtrafik       | On-demand     | TAXI     |

## CLI Options (Trip Planning)

| Option            | Short | Description                              | Default        |
|-------------------|-------|------------------------------------------|----------------|
| `--time`          | `-t`  | Departure time (HH:MM)                   | Now            |
| `--date`          | `-d`  | Departure date (YYYY-MM-DD)              | Today          |
| `--arrive`        | `-a`  | Search by arrival time instead           | false          |
| `--results`       | `-n`  | Number of results (1-6)                  | 3              |
| `--changes`       | `-c`  | Max number of changes (0-9)              | unlimited      |
| `--lang`          | `-l`  | Language (sv/en)                         | sv             |
| `--json`          | `-j`  | Output raw JSON                          | false          |

## Next Command

Show next departures of a specific transport mode from a location towards a destination.

```bash
transport next <mode> <location> <towards>
```

### Arguments

| Argument    | Description                                           |
|-------------|-------------------------------------------------------|
| `mode`      | Transport type: bus, metro, train, tram, ship         |
| `location`  | Stop/station name to depart from                      |
| `towards`   | Final destination to filter by                        |

### Options

| Option | Short | Description                    | Default |
|--------|-------|--------------------------------|---------|
| `-n`   |       | Number of departures to show   | 3       |
| `-l`   |       | Language (sv/en)               | sv      |

### Example Output

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 🚌 Nästa buss från Styresman Sanders väg mot Brommaplan
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  🚌 117  Brommaplan                om  6 min (11:13)
  🚌 117  Brommaplan                om 18 min (11:24)
  🚌 117  Brommaplan                om 33 min (11:39)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## Project Structure

```
transport/
├── cmd/
│   └── transport/
│       └── main.go          # CLI entry point
├── internal/
│   ├── api/
│   │   ├── sl.go            # SL Journey Planner v2 client
│   │   └── types.go         # API response types
│   ├── planner/
│   │   └── planner.go       # Journey planning logic
│   ├── display/
│   │   ├── formatter.go     # Output formatting
│   │   └── icons.go         # Transport type symbols
│   └── config/
│       └── config.go        # Configuration
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Technology Stack

- **Language:** Go 1.21+
- **HTTP Client:** net/http (stdlib)
- **CLI Parsing:** cobra or stdlib flag
- **Output Styling:** fatih/color + custom formatting
- **Build:** Single static binary, cross-platform

## Configuration

```bash
# ~/.config/transport/config.json or environment variables
{
  "default_location": "Slussen",    # Default origin
  "language": "sv",
  "results": 3
}

# Environment variables
TRANSPORT_DEFAULT_LOCATION=Slussen
```

## API

Uses **SL Journey Planner v2** exclusively (no API key required):
- Base URL: `https://journeyplanner.integration.sl.se/v2/`
- Covers all of Stockholm County (Storstockholms Lokaltrafik)

## Location Resolution

1. **Stop names:** "Slussen", "T-Centralen" → lookup via stop-finder
2. **Addresses:** "Kungsgatan 1" → geocode to coordinates
3. **Coordinates:** "59.3293,18.0686" → use directly
4. **Current location:** Use IP geolocation or saved default

## Error Handling

```
Error: Could not find location "Slusssen"

Did you mean:
  1. Slussen (Tunnelbana)
  2. Slussen (Buss)
  3. Slussen kajen

Use: transport --select 1 "Odenplan"
```

## Real-time Information

When available, show delays:
```
  08:15  Slussen                              Spår 1
         ⚠️  Försenad, ny avgång 08:18 (+3 min)
```

## Phase 1 Implementation (MVP)

1. Basic CLI with positional arguments
2. SL API integration
3. Simple text output with icons
4. Stop name lookup with fuzzy matching

## Phase 2 Enhancements

1. Colored output
2. Real-time delay information
3. Configuration file support
4. Time/date options

## Phase 3 Features

1. Favorites/aliases ("home", "work")
2. Interactive mode with fuzzy search
3. Default location from config
