import { useState, useEffect, useRef, Suspense } from 'react'
import { SearchIcon, TrendingUp } from 'lucide-react'
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { CoinData, TrendingCoin } from "@/lib/types"
import { fetchCoinData, fetchTrendingCoins, fetchCoinSuggestions } from "@/lib/utils"
import { CombinedChart } from './charts/stock-charts'
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
// import {terminal} from 'virtual:terminal'

export function CryptoSearch() {
  const [searchTerm, setSearchTerm] = useState('')
  const [coinData, setCoinData] = useState<CoinData | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [trendingCoins, setTrendingCoins] = useState<TrendingCoin[]>([])
  const [showTrending, setShowTrending] = useState(false)
  const [suggestions, setSuggestions] = useState<TrendingCoin[]>([])
  const [typedText, setTypedText] = useState('')
  const [showCursor, setShowCursor] = useState(true)
  const [shouldAnimate, setShouldAnimate] = useState(true)
  const [timeSpan, setTimeSpan] = useState('30')
  const searchRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    const loadTrendingCoins = async () => {
      try {
        const trending = await fetchTrendingCoins()
        setTrendingCoins(trending)
      } catch (error) {
        console.error('Error fetching trending coins:', error)
        setTrendingCoins([])
      }
    }
    loadTrendingCoins()
  }, [])

  useEffect(() => {
    const getSuggestions = async () => {
      if (searchTerm.length > 1) {
        try {
          const fetchedSuggestions = await fetchCoinSuggestions(searchTerm)
          setSuggestions(fetchedSuggestions)
        } catch (error) {
          console.error('Error fetching suggestions:', error)
          setSuggestions([])
        }
      } else {
        setSuggestions([])
      }
    }
    getSuggestions()
  }, [searchTerm])

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (searchRef.current && !searchRef.current.contains(event.target as Node)) {
        setShowTrending(false)
        setSuggestions([])
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [])

  useEffect(() => {
    if (shouldAnimate) {
      const text = "CryptWatch"
      let i = 0
      setTypedText('')
      setShowCursor(true)
      const typingInterval = setInterval(() => {
        if (i < text.length) {
          setTypedText(_ => text.slice(0, i + 1))
          i++
        } else {
          clearInterval(typingInterval)
          setShowCursor(false)
          setShouldAnimate(false)
        }
      }, 150)

      return () => clearInterval(typingInterval)
    }
  }, [shouldAnimate])

  useEffect(() => {
    if (coinData) {
      handleSearch(coinData.name)
    }
  }, [timeSpan])

  const handleSearch = async (coinID: string = searchTerm) => {
    if (!coinID) return
    setIsLoading(true)
    try {
      const data = await fetchCoinData(coinID, parseInt(timeSpan))
      setCoinData(data)
      // terminal.log('Fetched data:', data)
      setShouldAnimate(false)
    } catch (error) {
      console.error('Error fetching coin data:', error)
      setCoinData(null)
    } finally {
      setIsLoading(false)
    }
  }

  const handleCoinClick = (coin: TrendingCoin) => {
    setSearchTerm(coin.name)
    setShowTrending(false)
    setSuggestions([])
    handleSearch(coin.id)
  }

  const handleReturnToSearch = () => {
    setCoinData(null)
    setSearchTerm('')
    setTypedText('')
    setShowCursor(true)
    setShouldAnimate(true)
    if (inputRef.current) {
      inputRef.current.focus()
    }
  }

  return (
    <div className="container mx-auto p-4 min-h-screen flex flex-col">
      <div className={`transition-all duration-300 ease-in-out ${coinData ? 'mb-8' : 'flex-grow flex items-center justify-center'}`}>
        <div className="w-full max-w-md mx-auto">
          <h1 
            className="text-4xl sm:text-5xl md:text-6xl font-bold mb-8 text-center bg-gradient-to-r from-black via-orange-800 to-orange-500 text-transparent bg-clip-text cursor-pointer"
            onClick={handleReturnToSearch}
          >
            {typedText}
            {showCursor && <span className="animate-blink">|</span>}
          </h1>
          <div className="relative mb-4" ref={searchRef}>
            <div className="flex gap-2">
              <Input
                type="text"
                placeholder="Search for a coin..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                onFocus={() => setShowTrending(true)}
                className="flex-grow"
                ref={inputRef}
              />
              <Button onClick={() => handleSearch()} disabled={isLoading}>
                {isLoading ? 'Searching...' : <SearchIcon className="h-4 w-4" />}
              </Button>
            </div>
            {((suggestions && suggestions.length > 0) || (showTrending && trendingCoins && trendingCoins.length > 0)) && (
              <Card className="absolute z-10 w-full mt-1">
                <CardContent>
                  <ul>
                    {suggestions && suggestions.length > 0 ? (
                      suggestions.map((coin) => (
                        <li
                          key={coin.id}
                          className="cursor-pointer hover:bg-muted p-2 rounded"
                          onClick={() => handleCoinClick(coin)}
                        >
                          {coin.name} ({coin.symbol})
                        </li>
                      ))
                    ) : trendingCoins && trendingCoins.length > 0 ? (
                      <>
                        <CardHeader>
                          <CardTitle className="text-sm flex items-center">
                            <TrendingUp className="h-4 w-4 mr-2" />
                            Trending Coins
                          </CardTitle>
                        </CardHeader>
                        {trendingCoins.map((coin) => (
                          <li
                            key={coin.id}
                            className="cursor-pointer hover:bg-muted p-2 rounded"
                            onClick={() => handleCoinClick(coin)}
                          >
                            {coin.name} ({coin.symbol})
                          </li>
                        ))}
                      </>
                    ) : null}
                  </ul>
                </CardContent>
              </Card>
            )}
          </div>
        </div>
      </div>
      {coinData && (
        <div className="flex flex-wrap gap-4">
          <Card className="flex-grow basis-full md:basis-[calc(50%-0.5rem)] xl:basis-[calc(66.66%-0.5rem)]">
            <CardHeader>
              <CardTitle>{coinData.name} Details</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex flex-wrap">
                {Object.entries({
                  Name: coinData.name,
                  Symbol: coinData.symbol,
                  'Price (USD)': `$${coinData.price}`,
                  Rank: coinData.rank,
                  Supply: coinData.supply,
                  'Market Cap': `$${coinData.marketCap}`
                }).map(([key, value]) => (
                  <div key={key} className="w-full sm:w-1/2 p-2">
                    <p className="font-semibold">{key}:</p>
                    <p className="break-words">{value}</p>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
          <Card className="flex-grow basis-full md:basis-[calc(50%-0.5rem)] xl:basis-[calc(33.33%-0.5rem)]">
            <CardHeader>
              <CardTitle>{coinData.name} Price Conversions</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex flex-wrap">
                {Object.entries(coinData.conversions).map(([currency, value]) => (
                  <div key={currency} className="w-full sm:w-1/2 p-2">
                    <p className="font-semibold">{currency}:</p>
                    <p>{value}</p>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
          <Card className="w-full">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle>{coinData.name} Price and OHLC Chart</CardTitle>
              <Select
                value={timeSpan}
                onValueChange={(value) => setTimeSpan(value)}
              >
                <SelectTrigger className="w-[180px]">
                  <SelectValue placeholder="Select time span" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="7">7 days</SelectItem>
                  <SelectItem value="30">30 days</SelectItem>
                  <SelectItem value="90">90 days</SelectItem>
                  <SelectItem value="365">1 year</SelectItem>
                </SelectContent>
              </Select>
            </CardHeader>
            <CardContent>
              <Suspense fallback={<div>Loading chart...</div>}>
                <CombinedChart data={coinData.chartData} />
              </Suspense>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  )
}