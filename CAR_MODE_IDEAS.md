# Car Mode - Future Ideas

## Current Implementation
- VW Tiguan Allspace 2018 (Diesel)
- Fuel consumption: 9.0 L/100km (<20km), 7.0 L/100km (≥20km)
- Tank: 58L, range ~700km
- Google Maps URL generation
- Fuel stop calculation based on distance and tank level

## Routing API Integration

### OSRM (Open Source Routing Machine) ✓ FREE
- **Endpoint:** `https://router.project-osrm.org/route/v1/driving/{lon1},{lat1};{lon2},{lat2}`
- **Returns:** Distance (meters), duration (seconds)
- **Rate limit:** Light use only (public demo server)
- **Example:**
  ```bash
  curl "https://router.project-osrm.org/route/v1/driving/18.0686,59.3293;14.5208,63.1792"
  # Returns: 564 km, ~7.4 hours (Stockholm → Åre)
  ```

### OpenRouteService
- **Free tier:** 2,000 requests/day
- **Requires:** API key (free registration)
- **Max distance:** 6,000 km
- **URL:** https://openrouteservice.org/

### Trafikverket/NVDB
- Swedish national road database
- Free account at https://lastkajen.trafikverket.se
- More for road data than routing

## Fuel Price APIs

### Current Status: No fully free real-time Swedish API

| Source | Data | Free? | Notes |
|--------|------|-------|-------|
| bensinpriser.nu | Crowdsourced prices | Scraping only | No official API |
| GlobalPetrolPrices | National averages | 2-week trial | Then paid |
| fuel_prices_sweden | Home Assistant | Free | Scrapes bensinpriser.nu |

### Current Average (Jan 2026)
- **Diesel:** ~16.39 SEK/liter
- **Bensin 95:** ~17.50 SEK/liter

### Workarounds
1. Hardcode national average price (update periodically)
2. Link to bensinpriser.nu with location filter:
   ```
   https://bensinpriser.nu/stationer/diesel/{region}/alla
   ```
3. Calculate estimated fuel cost using average price

## TODO: Future Enhancements

### Phase 1: Auto-distance with OSRM
- [ ] Integrate OSRM API to auto-fetch distance
- [ ] Geocode addresses to coordinates (need geocoding API)
- [ ] Cache common routes

### Phase 2: Fuel Cost Estimation
- [ ] Add `-p` flag for fuel price (default to current average)
- [ ] Show estimated trip cost in SEK
- [ ] Option to fetch current average from GlobalPetrolPrices (if API available)

### Phase 3: Fuel Stop Suggestions
- [ ] Generate Google Maps search URLs for diesel stations at calculated stop points
- [ ] Link to bensinpriser.nu for price comparison at stop locations

### Phase 4: Multiple Vehicles
- [ ] Config file for vehicle profiles
- [ ] Support for different fuel types (bensin, diesel, el, hybrid)
- [ ] Electric vehicle support with charging station APIs

## Useful Links
- OSRM API Docs: https://project-osrm.org/docs/v5.5.1/api/
- OpenRouteService: https://openrouteservice.org/
- bensinpriser.nu: https://bensinpriser.nu/
- GlobalPetrolPrices Sweden: https://www.globalpetrolprices.com/Sweden/diesel_prices/
- fuel_prices_sweden GitHub: https://github.com/deler-aziz/fuel_prices_sweden
- Trafikverket NVDB: https://www.nvdb.se/sv/about-nvdb/
