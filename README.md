# Transport

A command-line tool for planning trips in Sweden. Supports public transport, flights, taxis, long-distance buses, and car trips with fuel calculations.

## Features

- **Public Transport** - Trip planning via SL (Stockholm) or ResRobot (nationwide)
- **Next Departures** - Real-time departures for bus, metro, train, tram, and ferry
- **Flight Search** - Generate booking links for flights (Skyscanner, Google Flights, Norwegian, etc.)
- **Nearby Airports** - Find airports near any Swedish city
- **Taxi** - Fare estimates for Taxi Stockholm, Taxi Kurir, Uber, and Bolt
- **Long-distance Bus** - FlixBus, Vy Bus4You, and Flygbussarna airport buses
- **Car Directions** - Route planning with fuel consumption and gas station suggestions

Supports both English and Swedish commands.

## Installation

### Download Pre-built Installers

| Platform | Download | Notes |
|----------|----------|-------|
| **Windows** | [transport_0.1.0_amd64_setup.exe](https://github.com/rAgamemnon74/transport/releases/download/v0.1.0/transport_0.1.0_amd64_setup.exe) | Adds to PATH automatically |
| **macOS** | [transport_0.1.0_macos.dmg](https://github.com/rAgamemnon74/transport/releases/download/v0.1.0/transport_0.1.0_macos.dmg) | Intel & Apple Silicon |
| **Linux (Debian/Ubuntu)** | [transport_0.1.0_amd64.deb](https://github.com/rAgamemnon74/transport/releases/download/v0.1.0/transport_0.1.0_amd64.deb) | `sudo dpkg -i transport_0.1.0_amd64.deb` |

### Build from Source

Requires Go 1.22+

```bash
git clone https://github.com/rAgamemnon74/transport.git
cd transport
make build
make install  # Installs to ~/.local/bin/
```

### Build Installers

```bash
make installers  # Builds all installers (DEB, Windows, DMG)
```

## Usage

### Public Transport (Stockholm)

Plan a trip within Stockholm using SL:

```bash
# Trip from Slussen to Odenplan
transport Slussen Odenplan

# Trip to a destination (uses TRANSPORT_DEFAULT_LOCATION as origin)
transport Odenplan

# Depart at a specific time
transport -t 08:30 Slussen T-Centralen

# Arrive by a specific time
transport -a -t 09:00 "Spånga station" "Stockholm Central"

# Limit number of changes
transport -c 1 Slussen Kista
```

### Public Transport (Nationwide)

Search all of Sweden using ResRobot (requires API key):

```bash
# Gothenburg to Stockholm
transport -se Göteborg "Stockholm Central"

# Sundsvall to Ånge
transport -se Sundsvall Ånge

# Malmö to Lund
transport -se Malmö Lund
```

**Setup:** Get a free API key at [Trafiklab](https://www.trafiklab.se/api/trafiklab-apis/resrobot-v21/) and set it:

```bash
export RESROBOT_API_KEY="your-key-here"
```

### Next Departures

Show real-time departures from a stop:

```bash
# Next buses from Odenplan
transport nästa buss Odenplan
transport next bus Odenplan

# Next metro from Slussen
transport nästa tunnelbana Slussen
transport next metro Slussen

# Next trains from Stockholm Central
transport nästa tåg "Stockholm Central"
transport next train "Stockholm Central"

# Filter by direction
transport nästa buss "Spånga station" Brommaplan

# Show more departures
transport next -n 10 bus Odenplan
```

**Supported modes:** `bus`/`buss`, `metro`/`tunnelbana`/`t-bana`, `train`/`tåg`, `tram`/`spårvagn`, `ship`/`båt`/`färja`

### Flight Search

Generate booking links for flights:

```bash
# Stockholm to Vilnius
transport fly from stockholm to vilnius
transport flyga från stockholm till vilnius

# With specific dates
transport fly -d 2026-03-15 from göteborg to barcelona

# Round trip
transport fly -d 2026-06-01 -r 2026-06-08 ARN BCN

# Include private jet options
transport fly -p from bromma to visby
```

### Nearby Airports

Find airports near a location:

```bash
# Airports near Stockholm (default 100km radius)
transport flyg
transport flight

# Airports near Gothenburg with 150km radius
transport flyg -r 150 Göteborg

# Only airports with scheduled service
transport flyg -s Malmö

# Airports near Kiruna
transport flight -r 200 Kiruna
```

### Taxi

Get fare estimates and booking links:

```bash
# Slussen to Arlanda
transport taxi Slussen Arlanda

# With addresses
transport taxi "Kungsgatan 1" "Arlanda Terminal 5"

# Natural language
transport taxi from T-Centralen to Bromma Airport
```

Shows estimates for:
- Taxi Stockholm
- Taxi Kurir
- Uber
- Bolt

### Long-distance Bus

Search FlixBus, Vy Bus4You, and Flygbussarna:

```bash
# Stockholm to Gothenburg
transport buss Stockholm Göteborg

# Malmö to Oslo
transport buss from Malmö to Oslo

# With specific date
transport buss -d 2026-03-15 Stockholm Arlanda

# Airport buses
transport buss Stockholm Arlanda
transport buss Göteborg Landvetter
```

### Car Directions

Get driving directions with fuel calculations:

```bash
# Basic route
transport bil Stockholm Göteborg
transport car Stockholm Gothenburg

# With known distance
transport bil -d 620 Stockholm Åre

# Starting with half tank
transport bil -d 620 -f 50 Stockholm Åre

# Local trip
transport bil Bromma Arlanda
```

Shows:
- Google Maps link
- Estimated fuel consumption
- Gas station suggestions along the route

## Configuration

### Default Location

Set a default origin to skip typing it every time:

```bash
export TRANSPORT_DEFAULT_LOCATION="Slussen"

# Now you can just type the destination
transport Odenplan  # Same as: transport Slussen Odenplan
```

### API Key for Nationwide Search

For searching outside Stockholm, you need a ResRobot API key:

1. Go to [Trafiklab](https://www.trafiklab.se/api/trafiklab-apis/resrobot-v21/)
2. Create a free account
3. Get an API key for ResRobot v2.1
4. Set the environment variable:

```bash
export RESROBOT_API_KEY="your-key-here"
```

## Options Reference

### Trip Planning

| Option | Description |
|--------|-------------|
| `-t`, `--time` | Departure time (HH:MM) |
| `-d`, `--date` | Departure date (YYYY-MM-DD) |
| `-a`, `--arrive` | Search by arrival time |
| `-c`, `--changes` | Maximum number of changes (0-9) |
| `-n`, `--results` | Number of results (1-6) |
| `-se`, `--sweden` | Search nationwide (ResRobot) |
| `-l`, `--lang` | Language (sv/en) |
| `-v`, `--version` | Show version |

### Flight Search

| Option | Description |
|--------|-------------|
| `-d`, `--date` | Departure date (YYYY-MM-DD) |
| `-r`, `--return` | Return date for round-trip |
| `-p`, `--private` | Show private jet & helicopter options |

### Nearby Airports

| Option | Description |
|--------|-------------|
| `-r`, `--radius` | Search radius in km (default: 100) |
| `-s`, `--scheduled` | Only show airports with scheduled service |

### Car

| Option | Description |
|--------|-------------|
| `-d`, `--distance` | Distance in km (if known) |
| `-f`, `--fuel` | Starting fuel level in % (default: 100) |

## License

MIT
