# Taxi Mode - Stockholm Region ✓ IMPLEMENTED

## Quick Usage

```bash
transport taxi Slussen Arlanda
transport taxi from T-Centralen to Bromma Airport
transport taxi "Kungsgatan 1" "Globen"
```

**Features:**
- Geocodes addresses via OpenStreetMap Nominatim
- Calculates route via OSRM (distance & duration)
- Estimates fares for Taxi Stockholm, Taxi Kurir, Uber, Bolt
- Detects airport routes → shows fixed prices
- Generates Uber deep link (opens app with destination)
- Generates booking URLs for all companies
- Shows Google Maps route preview

## Reputable Companies

### Traditional Taxi Companies

| Company | Phone | Cars | App | Website |
|---------|-------|------|-----|---------|
| **Taxi Stockholm** | 020-93 93 93 | 1,500 | iOS/Android | taxistockholm.se |
| **Taxi Kurir** | 08-30 00 00 | 1,000+ | iOS/Android | taxikurir.se |
| **Sverigetaxi** | 020-20 20 20 | 1,000 | iOS/Android | sverigetaxi.se |
| **Taxi 020** | 020-20 20 20 | - | Via Cabonline | taxi020.se |

### App-Based Services

| Service | Type | Website |
|---------|------|---------|
| **Uber** | Ride-hailing | uber.com |
| **Bolt** | Ride-hailing | bolt.eu |
| **Cabonline** | Aggregator | cabonline.com |

## Pricing (January 2026)

### Taxi Stockholm

**Regular Taxi (1-4 passengers)**

| Tariff | When | Base | Per km | Per hour | Comparison* |
|--------|------|------|--------|----------|-------------|
| 1 | Weekdays | 59 kr | 14.90 kr | 565 kr | 349 kr |
| 2 | Fri 15:00 - Mon 06:00, holidays | 59 kr | 14.90 kr | 595 kr | 357 kr |
| F | Waiting time | 59 kr | 0 kr | 565 kr | 200 kr |

**Large Taxi (5-7 passengers)**

| Tariff | When | Base | Per km | Per hour | Comparison* |
|--------|------|------|--------|----------|-------------|
| 6 | All times | 90 kr | 23.00 kr | 673 kr | 494 kr |
| 5 | Waiting time | 90 kr | 0 kr | 695 kr | 264 kr |

**Airport/Hub Fees**
- Arlanda: 50 kr + 30 kr booking fee
- Bromma: 24 kr + 30 kr pre-order
- Skavsta: 43 kr
- Central Station: 25 kr
- Climate fee (large cars): 64 kr

### Taxi Kurir

**Standard Tariffs**

| Tariff | Passengers | Base | Per km | Per hour | Comparison* |
|--------|------------|------|--------|----------|-------------|
| 1 | 1-4 | 55 kr | 14.60 kr | 576 kr | 345 kr |
| 6 | 5-8 | 85 kr | 23.75 kr | 690 kr | 495 kr |
| F | Waiting | 58 kr | 0 kr | 565 kr | 199 kr |

**Fixed Prices (Arlanda ↔ Stockholm City)**
- Standard (1-4 pax): 695 kr
- Large (5-7 pax): 1,095 kr

**Airport/Hub Fees**
- Bromma/Central Station: 30 kr
- Arlanda pre-booked: 100 kr
- Child seat: 150 kr
- Environmental fee (large): 75 kr

### Uber (Approximate, dynamic pricing)

| Service | Base | Per km | Per min | Min fare |
|---------|------|--------|---------|----------|
| UberX | 20-35 kr | 8-12 kr | 1-2 kr | ~100 kr |
| Comfort | +20-40% | | | |
| Black | 50-70 kr | 15-25 kr | 2-3.5 kr | ~175 kr |

Note: Surge pricing applies during rush hours, events, bad weather (1.3x-2x typical).

### Sverigetaxi

No public tariff available. Contact: 020-20 20 20

## Fare Calculation Formula

```
fare = base_fee + (distance_km * per_km_rate) + (time_hours * hourly_rate) + hub_fees
```

For a typical 10 km, 15 min trip:
- Taxi Stockholm: 59 + (10 × 14.90) + (0.25 × 565) = 349 kr
- Taxi Kurir: 55 + (10 × 14.60) + (0.25 × 576) = 345 kr

## Booking URLs

### Deep Links

```
# Taxi Stockholm (web booking)
https://www.taxistockholm.se/en/booking/

# Taxi Kurir (web booking)
https://www.taxikurir.se/boka

# Uber (with destination)
https://m.uber.com/ul/?action=setPickup&pickup=my_location&dropoff[formatted_address]=DESTINATION

# Bolt
https://bolt.eu/

# Sverigetaxi
https://www.sverigetaxi.se/boka-taxi
```

## Routing API Options (OpenStreetMap-based)

### OSRM - Routing ✓ TESTED & WORKING

**Open Source Routing Machine** - Free, no API key required.

```bash
# Example: Slussen → Arlanda
curl "https://router.project-osrm.org/route/v1/driving/18.0722,59.3193;17.9237,59.6498"

# Response (key fields):
{
  "routes": [{
    "distance": 43791.5,    # meters (43.8 km)
    "duration": 2513.2      # seconds (42 min)
  }]
}
```

**API Format:** `/route/v1/driving/{lon1},{lat1};{lon2},{lat2}`

**Features:**
- Returns distance (meters) and duration (seconds)
- Supports waypoints
- Very fast (~1ms query time)
- Public demo server for light use

**Rate Limits:** Light use only on demo server. For production, self-host.

### Nominatim - Geocoding ✓ TESTED & WORKING

**OpenStreetMap Geocoding** - Free, no API key required.

```bash
# Example: Geocode "Arlanda Airport"
curl "https://nominatim.openstreetmap.org/search?q=Arlanda+Airport,Stockholm&format=json&limit=1" \
  -H "User-Agent: transport-cli/1.0"

# Response (key fields):
{
  "lat": "59.6467921",
  "lon": "17.9370443",
  "name": "Stockholm-Arlanda flygplats",
  "display_name": "Stockholm-Arlanda flygplats, Sigtuna kommun, Stockholms län, Sverige"
}
```

**API Format:** `/search?q={query}&format=json&limit=1`

**Requirements:**
- Must include User-Agent header
- Rate limit: 1 request/second
- Add `countrycodes=se` to limit to Sweden

### Alternative APIs

| Service | Free Tier | API Key | Notes |
|---------|-----------|---------|-------|
| **OpenRouteService** | 2,000 req/day | Required | More features |
| **GraphHopper** | 500 req/day | Required | Good for routing |
| **Photon** | Unlimited | No | Fast geocoding |
| **Valhalla** | Self-host | No | Full-featured |

## Tested Route Example: Slussen → Arlanda

```
Geocoded coordinates:
  Slussen:  59.3193, 18.0722
  Arlanda:  59.6498, 17.9237

OSRM Result:
  Distance: 43.8 km
  Duration: 42 min

Fare Estimates:
  Taxi Stockholm: 59 + (43.8 × 14.90) + (0.7 × 565) = 59 + 653 + 396 = 1,108 kr
  Taxi Kurir:     55 + (43.8 × 14.60) + (0.7 × 576) = 55 + 639 + 403 = 1,097 kr
  Fixed price:    695 kr (Taxi Kurir) ← BETTER DEAL!
```

## Implementation Plan

### Phase 1: Basic Taxi Command ✓ IMPLEMENTED
- [x] Geocode addresses via Nominatim
- [x] Calculate route via OSRM (distance/duration)
- [x] Estimate fares for Taxi Stockholm & Taxi Kurir
- [x] Generate Uber deep link with coordinates
- [x] Generate web booking URLs for Swedish taxis
- [x] Show Google Maps route preview
- [x] Detect airport routes → show fixed prices

### Phase 2: Future Enhancements
- [ ] Time-based tariff selection (weekday vs weekend)
- [ ] Large group option (-l flag for 5-8 passengers)
- [ ] Cache geocoding results for common locations
- [ ] Add more taxi companies (Sverigetaxi, Cabonline)

## Android App Deep Links

### Uber ✓ DOCUMENTED

Uber has official deep link support with pickup/dropoff coordinates.

```
https://m.uber.com/ul/?action=setPickup
  &pickup[latitude]=59.3193
  &pickup[longitude]=18.0722
  &pickup[nickname]=Slussen
  &dropoff[latitude]=59.6498
  &dropoff[longitude]=17.9237
  &dropoff[nickname]=Arlanda
  &dropoff[formatted_address]=Stockholm-Arlanda+Airport
```

**Simplified format:**
```
https://m.uber.com/ul/?action=setPickup&pickup=my_location&dropoff[latitude]=59.6498&dropoff[longitude]=17.9237&dropoff[nickname]=Arlanda
```

- Opens Uber app if installed, otherwise mobile web
- Requires `dropoff[nickname]` OR `dropoff[formatted_address]`
- Source: [Uber Developer Docs](https://developer.uber.com/docs/riders/ride-requests/tutorials/deep-links/introduction)

### Bolt ✗ NO PUBLIC DEEP LINK

- Package ID: `ee.mtakso.client`
- No documented URI scheme found
- Fallback: Link to Play Store or website

**Workaround - Android Intent (requires app code):**
```
intent://ride#Intent;scheme=bolt;package=ee.mtakso.client;end
```

### Taxi Stockholm ✗ NO PUBLIC DEEP LINK

- Package ID: `se.taxistockholm`
- No documented URI scheme
- Best option: Link to web booking

**Web booking URL:**
```
https://www.taxistockholm.se/en/booking/
```

### Taxi Kurir ✗ NO PUBLIC DEEP LINK

- Package ID: `se.taxikurir.app`
- No documented URI scheme
- Best option: Link to web booking

**Web booking URL:**
```
https://www.taxikurir.se/boka
```

### Sverigetaxi ✗ NO PUBLIC DEEP LINK

- Package ID: `se.sverigetaxi.app`
- No documented URI scheme

**Web booking URL:**
```
https://www.sverigetaxi.se/boka-taxi
```

## Fallback: Google Maps Intent

For apps without deep links, use Google Maps navigation intent:

```
# Opens Google Maps with directions (works on Android)
https://www.google.com/maps/dir/?api=1
  &origin=Slussen,Stockholm
  &destination=Arlanda+Airport
  &travelmode=driving

# geo: URI (opens map app chooser on Android)
geo:59.6498,17.9237?q=Arlanda+Airport
```

## App Package IDs (Android)

| App | Package ID | Deep Link |
|-----|------------|-----------|
| Uber | `com.ubercab` | ✓ Yes |
| Bolt | `ee.mtakso.client` | ✗ No |
| Taxi Stockholm | `se.taxistockholm` | ✗ No |
| Taxi Kurir | `se.taxikurir.app` | ✗ No |
| Sverigetaxi | `se.sverigetaxi.app` | ✗ No |
| Cabonline | `se.cabonline.passenger` | ✗ No |

## Summary: What Works on Android

| Method | Opens App | Pre-fills Route |
|--------|-----------|-----------------|
| **Uber deep link** | ✓ Yes | ✓ Yes |
| **Web booking URLs** | Via browser | ✗ No |
| **Google Maps link** | ✓ Maps app | ✓ Yes (navigation only) |
| **Play Store link** | ✓ Store | ✗ No |

## Recommended Implementation

```go
// For Uber - use deep link with coordinates
func GenerateUberDeepLink(pickupLat, pickupLon, dropLat, dropLon float64, dropName string) string {
    return fmt.Sprintf(
        "https://m.uber.com/ul/?action=setPickup&pickup=my_location"+
        "&dropoff[latitude]=%f&dropoff[longitude]=%f&dropoff[nickname]=%s",
        dropLat, dropLon, url.QueryEscape(dropName))
}

// For others - use web booking + Google Maps for route preview
func GenerateTaxiStockholmURL() string {
    return "https://www.taxistockholm.se/en/booking/"
}

func GenerateGoogleMapsURL(from, to string) string {
    return fmt.Sprintf(
        "https://www.google.com/maps/dir/?api=1&origin=%s&destination=%s&travelmode=driving",
        url.QueryEscape(from), url.QueryEscape(to))
}
```

## Useful Links

- Taxi Stockholm: https://www.taxistockholm.se/en/prices/
- Taxi Kurir: https://www.taxikurir.se/stockholm
- Uber Deep Links: https://developer.uber.com/docs/riders/ride-requests/tutorials/deep-links/introduction
- OSRM Demo: https://router.project-osrm.org/
- Nominatim: https://nominatim.openstreetmap.org/
- OpenRouteService: https://openrouteservice.org/
