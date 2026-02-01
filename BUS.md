# Bus Transport in Sweden

## Overview

Sweden has three categories of bus transport:
1. **Regional/Local buses** (Länstrafik) - Operated by 21 regional authorities
2. **Commercial long-distance buses** - FlixBus, Vy Bus4You
3. **Airport buses** - Flygbussarna (Vy)

## 1. Regional Bus Operators (Länstrafik)

Sweden has 21 regions, each with a Public Transport Authority (PTA) responsible for local and regional buses.

### Major Metropolitan Operators

| Region | Operator | Website | Coverage |
|--------|----------|---------|----------|
| Stockholm | **SL** | sl.se | Stockholm County |
| Västra Götaland | **Västtrafik** | vasttrafik.se | Gothenburg region |
| Skåne | **Skånetrafiken** | skanetrafiken.se | Malmö, Lund, Helsingborg |
| Uppsala | **UL** | ul.se | Uppsala County |

### All 21 Regional Operators

| Region | Operator | Website |
|--------|----------|---------|
| Blekinge | Blekingetrafiken | blekingetrafiken.se |
| Dalarna | Dalatrafik | dalatrafik.se |
| Gotland | Gotlandstrafiken | gotland.se/kollektivtrafik |
| Gävleborg | X-trafik | xtrafik.se |
| Halland | Hallandstrafiken | hallandstrafiken.se |
| Jämtland | Länstrafiken Jämtland | dintur.se |
| Jönköping | JLT | jlt.se |
| Kalmar | Kalmar länstrafik | klt.se |
| Kronoberg | Länstrafiken Kronoberg | lanstrafikenkron.se |
| Norrbotten | Länstrafiken Norrbotten | ltnbd.se |
| Skåne | Skånetrafiken | skanetrafiken.se |
| Stockholm | SL | sl.se |
| Södermanland | Sörmlandstrafiken | sormlandstrafiken.se |
| Uppsala | UL | ul.se |
| Värmland | Värmlandstrafik | varmlandstrafik.se |
| Västerbotten | Länstrafiken Västerbotten | tabussen.nu |
| Västernorrland | Din Tur | dintur.se |
| Västmanland | VL | vl.se |
| Västra Götaland | Västtrafik | vasttrafik.se |
| Örebro | Länstrafiken Örebro | lanstrafikenorebro.se |
| Östergötland | Östgötatrafiken | ostgotatrafiken.se |

### Cross-Regional Travel

Tickets from regional operators are often valid for connections:
- Skånetrafiken tickets work with Hallandstrafiken, Västtrafik, Blekingetrafiken
- UL and SL have joint lines across county borders
- Mälartåg connects Stockholm, Uppsala, Örebro, Västmanland regions

## 2. Commercial Long-Distance Buses

### FlixBus ✓ MAJOR OPERATOR

**Coverage:** 35+ European countries, extensive Swedish network

**Swedish Routes:**
- Stockholm ↔ Gothenburg (4-5 hours)
- Stockholm ↔ Malmö (6-8 hours)
- Gothenburg ↔ Malmö (3-4 hours)
- Stockholm ↔ Oslo
- Stockholm ↔ Copenhagen
- Plus many more regional connections

**Prices:**
- Starting from 29 SEK (~€2.70)
- Stockholm → Gothenburg: from ~180 SEK
- Stockholm → Malmö: from ~280 SEK
- Gothenburg → Malmö: from ~150 SEK

**Amenities:**
- Free WiFi
- Power outlets
- Extra legroom
- Free luggage (1 checked + 1 carry-on)

**Booking:**
- Website: flixbus.se
- App: FlixBus (iOS/Android)

**Booking URL Pattern:**
```
https://shop.flixbus.se/search?departureCity=CITY_ID&arrivalCity=CITY_ID&departureDate=YYYY-MM-DD
```

### Vy Bus4You ✓ PREMIUM OPERATOR

**Background:** Formerly Nettbuss Express, owned by Norwegian Vy Group. "Sweden's most satisfied customers" 2011-2023.

**Main Routes:**
1. Stockholm ↔ Gothenburg (via Norrköping, Linköping, Jönköping, Borås, Landvetter)
2. Stockholm ↔ Oslo (via Karlstad, Örebro)
3. Oslo ↔ Copenhagen (via Swedish west coast)
4. Kalmar ↔ Arlanda Airport

**Prices:**
- Stockholm → Gothenburg: from ~220 SEK
- Stockholm → Oslo: from ~250 SEK

**Amenities:**
- Free WiFi
- Power outlets
- Comfortable seats
- Onboard toilet

**Booking:**
- Website: vybus4you.se / vy.se
- App: Vy (iOS/Android)

### Swebus (Historical Note)

Swebus Express was acquired and no longer operates as a separate brand. Routes absorbed by FlixBus and Vy.

## 3. Airport Buses (Flygbussarna)

### Vy Flygbussarna

**Airports Served:**
| Airport | Code | Routes |
|---------|------|--------|
| Stockholm Arlanda | ARN | Cityterminalen, Liljeholmen, Brommaplan |
| Stockholm Bromma | BMA | Cityterminalen |
| Gothenburg Landvetter | GOT | Gothenburg C, Lindholmen |
| Malmö Sturup | MMX | Malmö C |
| Stockholm Skavsta | NYO | Stockholm C |
| Västerås | VST | Stockholm C |
| Visby | VBY | Visby centrum |

**Prices (2026):**
- Arlanda ↔ Stockholm C: ~119-139 SEK (online discount)
- Landvetter ↔ Gothenburg C: ~119 SEK
- +99 SEK for companion ticket

**Travel Times:**
- Arlanda → Stockholm C: ~45 min
- Landvetter → Gothenburg C: ~35 min
- Bromma → Stockholm C: ~20 min

**Booking:**
- Website: flygbussarna.se
- App: Flygbussarna (iOS/Android)
- Tickets valid for 3 months, any departure

**Booking URL:**
```
https://www.flygbussarna.se/en/arlanda
https://www.flygbussarna.se/en/landvetter
```

## API & Data Access

### Trafiklab (Official Swedish Transport Data)

**GTFS Regional** - Per operator, high quality
```
https://www.trafiklab.se/api/gtfs-datasets/gtfs-regional/
```
- Covers all 21 regions
- Static + real-time data
- Vehicle positions (GPS)
- Free with API key (CC0 license)

**GTFS Sweden 3** - Aggregated national feed
```
https://www.trafiklab.se/api/gtfs-datasets/gtfs-sweden/
```
- All public transport in Sweden
- Single feed for all operators

### Commercial Bus APIs

| Operator | Official API | Alternative |
|----------|-------------|-------------|
| FlixBus | No public API | RapidAPI, Lyko.tech |
| Vy Bus4You | No public API | Omio, Busbud |
| Flygbussarna | No public API | Website scraping |

### Aggregator Platforms

These platforms search multiple operators:
- **Omio** (omio.com) - FlixBus, Vy, trains
- **Busbud** (busbud.com) - All major bus operators
- **Rome2rio** (rome2rio.com) - Multi-modal search

## Booking URL Patterns

### FlixBus
```
https://shop.flixbus.se/search?
  departureCity=88&           # Stockholm
  arrivalCity=96&             # Gothenburg
  departureDate=2026-02-15&
  adult=1
```

City IDs (examples):
- Stockholm: 88
- Gothenburg: 96
- Malmö: 1234
- Oslo: 1374
- Copenhagen: 1162

### Vy Bus4You
```
https://www.vy.se/en/traffic-and-routes/buses
```
(No deep link with pre-filled route)

### Flygbussarna
```
https://www.flygbussarna.se/en/arlanda    # Stockholm Arlanda
https://www.flygbussarna.se/en/landvetter # Gothenburg
https://www.flygbussarna.se/en/bromma     # Stockholm Bromma
```

### Omio (Aggregator)
```
https://www.omio.com/search-frontend/results?
  departurePosition=Stockholm&
  arrivalPosition=Gothenburg&
  outboundDate=2026-02-15&
  passengers=A1
```

## Price Comparison (Stockholm → Gothenburg)

| Operator | Price Range | Duration | Frequency |
|----------|-------------|----------|-----------|
| FlixBus | 99-299 SEK | 4.5-5.5h | 8-12/day |
| Vy Bus4You | 199-399 SEK | 4-5h | 6-8/day |
| SJ Train | 299-899 SEK | 3h | Hourly |

## Implementation Ideas

### Phase 1: Long-Distance Bus Search
```bash
transport buss stockholm göteborg
# → Show FlixBus, Vy Bus4You options with prices
# → Generate booking URLs
```

### Phase 2: Regional Bus Integration
```bash
transport buss -r Odenplan Kista    # Regional (SL)
transport buss -l Malmö Lund        # Local (Skånetrafiken)
```

### Phase 3: Airport Bus
```bash
transport flygbuss Arlanda
# → Show Flygbussarna times and prices
# → Integrated with flight arrivals?
```

## Useful Links

### Operators
- FlixBus Sweden: https://www.flixbus.se/
- Vy Bus4You: https://www.vy.se/en/traffic-and-routes/buses
- Flygbussarna: https://www.flygbussarna.se/en
- SL (Stockholm): https://sl.se/
- Västtrafik (Gothenburg): https://www.vasttrafik.se/
- Skånetrafiken (Malmö): https://www.skanetrafiken.se/

### Booking Aggregators
- Omio: https://www.omio.com/buses
- Busbud: https://www.busbud.com/en/country/se
- Rome2rio: https://www.rome2rio.com/

### Data/API
- Trafiklab: https://www.trafiklab.se/
- GTFS Regional: https://www.trafiklab.se/api/gtfs-datasets/gtfs-regional/

## Notes

- **Swebus** no longer exists as a brand (absorbed by FlixBus)
- **Nettbuss** rebranded to **Vy** in 2019
- FlixBus acquired many European bus operators
- Regional operators (länstrafik) focus on local/commuter service
- Commercial operators (FlixBus, Vy) focus on intercity express service
- Airport buses (Flygbussarna) owned by Vy since 2020
