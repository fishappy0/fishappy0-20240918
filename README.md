## CryptWatch

<p style="font-size:23px"> A scuffed wrapper for the CoinGecko API /w Vite and Go Gin<p>

## Deployment instructions

    Requirement:
        - Require a CoinGecko API Key

    On local machines:
    1. Copy these global environment variables either to docker compose yml or make a separate file
    and pass the path to docker compose.
    ```
        POSTGRES_DB={YOUR DB NAME HERE, I used 'cryptwatchbe'}
        POSTGRES_USER={YOUR PG USER, I used 'postgres'}
        POSTGRES_PASSWORD={YOUR DB PASSWORD GOES HERE}
        POSTGRES_HOST=db
        POSTGRES_PORT=5434
        CG_API_KEY={YOUR KEY GOES HERE, STARTS WITH CG-abcxyzblabla}
        JWT_SECRET={SELECT YOUR PASS PHRASE}
        APP_MODE=production (development for running on local)

    ```
    2. Run `docker compose up`
    The Frontend should be accessible via localhost:3000 and the backend should be
    usable via postman at localhost:4008

    Configurations:

    + Due to vite requiring its own method of parsing .env, this file can be located in frontend/CryptWatchFE/.env if you wish to host the backend separately

## Quirks

    + "Sometimes clicking the search result won't work", "Changing the time duration", "Random 404 errors, 500 errors": These are all caused by CoinGecko's nasty ""30 requests per minute"" limit,

    Specifically, if you send requests back to back, which this app needed for OHLC, price (more in a minute about caching this), and conversions to fiats rate. And the cooldown for this is 30 seconds if it somehow feels like giving you a rate limit.

    Regarding caching: Due to the odd design of CG's API regarding fetching fiats prices for coins, it's not possible to send 1 request to grab all the fiats at once, same for getting coin details, it's heavily limited. Therefore, in order to cache a large chunk of coins, we would have to consume a couple hundred requests per hour (The limit for this is 10,000 per month).
    Similarly, there is no way to fetch all of the ohlc or prices for a coin at once, combined with the sheer number of coins there are, we would have to send at least 3 requests for the 3 automated granularity for N number of coins, let's say 20000, then it would be 60k request per caches. Furthermore, caches have to be recent, can't have a cache that is more than an hour old, since a lot can happen in the crypto land in 1 hour.
