import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"
import { CoinData, TrendingCoin, DetailedCoinApiResult} from "@/lib/types"
import axios from "axios"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

async function fetchCoinPrice(coinID: string, days: number): Promise<Array<{
  stamp: number,
  price: number
}>> {
  let response = await axios.get(`${import.meta.env.VITE_API_URL}/crypto/price?id=${coinID}&duration=${days}`)
  let result = []
  for (let i = 0; i < response.data["prices"].length; i++) {
    result.push({
      stamp: response.data["prices"][i][0],
      price: response.data["prices"][i][1],
    })
  }
  return result
}

async function fetchCoinOHLC(coinID: string, days: number): Promise<Array<{
  stamp: number,
  o: number,
  h: number,
  l: number,
  c: number
}>> {
  let response = await axios.get(`${import.meta.env.VITE_API_URL}/crypto/ohlc?id=${coinID}&duration=${days}`)
  let result = []
  for (let i = 0; i < response.data["ohlc"].length; i++) {
    result.push({
      stamp: response.data["ohlc"][i][0],
      o: response.data["ohlc"][i][1],
      h: response.data["ohlc"][i][2],
      l: response.data["ohlc"][i][3],
      c: response.data["ohlc"][i][4]
    })
  }
  return result
}

async function fetchConversionData(coinID: string): Promise<{[key:string]:number}> {
  let response = await axios.get(`${import.meta.env.VITE_API_URL}/crypto/conversion?id=${coinID}`) 
  let data = response.data["conversions"]
  if (!data) {
    return { "MYR": 0.0, "CNY": 0.0, "EUR": 0.0, "SGD": 0.0, "VND": 0.0}
  }
  let result: { [key: string]: number } = {}
  for (let i = 0; i < data.length; i++) {
    result[data[i]["Symbol"]] = data[i]["Rate"]
  }
  return result
}

async function fetchDetailedData(coinID: string): Promise<DetailedCoinApiResult> {
  let response = await axios.get(`${import.meta.env.VITE_API_URL}/crypto/detailed?id=${coinID}`)
  let result: DetailedCoinApiResult = {
    name: response.data.Name,
    symbol: response.data.Symbol,
    price: response.data.Price,
    rank: response.data.Rank,
    supply: response.data.Supply,
    marketCap: response.data.Market
  }
  return result
}

export async function fetchCoinData(coinID: string, days: number): Promise<CoinData> {
  await new Promise(resolve => setTimeout(resolve, 1000)) 

  // const generateDynamicData = (days: number) => {
  //   const data = []
  //   const now = new Date()
  //   let lastClose = Math.random() * 10000 + 100
  //   for (let i = days; i >= 0; i--) {
  //     const date = new Date(now.getTime() - i * 24 * 60 * 60 * 1000)
  //     const open = lastClose + (Math.random() - 0.5) * (lastClose * 0.05)
  //     const high = Math.max(open, open + Math.random() * (open * 0.03))
  //     const low = Math.min(open, open - Math.random() * (open * 0.03))
  //     const close = (open + high + low) / 3 + (Math.random() - 0.5) * (open * 0.02)
      
  //     data.push({
  //       x: date.getTime(),
  //       y: [open, high, low, close].map(val => parseFloat(val.toFixed(2))),
  //       price: parseFloat(close.toFixed(2))
  //     })
      
  //     lastClose = close
  //   }
  //   return data
  // }

  // const dynamicData = generateDynamicData(days)
  // const latestPrice = dynamicData[dynamicData.length - 1].price
  // const latestPrice = priceData.chartData[priceData.chartData.length - 1].price
  let chart_data = []
  const priceData = await fetchCoinPrice(coinID, days)
  const ohlcData = await fetchCoinOHLC(coinID, days)
  const conversionData = await fetchConversionData(coinID)
  const detailedData = await fetchDetailedData(coinID)

  // x: date, y: [open, high, low, close], price: close
  for (let i = 0; i < ohlcData.length; i++) {
    let stamp = priceData[i].stamp
    let open = ohlcData[i].o
    let high = ohlcData[i].h
    let low = ohlcData[i].l
    let close = ohlcData[i].c
    chart_data.push({
      x: stamp,
      y: [open, high, low, close],
      price: close
    })
  }
  return {
    name: detailedData["name"],
    symbol: detailedData["symbol"],
    price: detailedData["price"],
    rank: detailedData["rank"],
    supply: detailedData["supply"],
    marketCap: detailedData["marketCap"],
    conversions: {
      MYR: conversionData["myr"],
      CNY: conversionData["cny"],
      EUR: conversionData["eur"],
      SGD: conversionData["sgd"],
      VND: conversionData["vnd"]
    },
    chartData: chart_data
  }
}

export async function fetchTrendingCoins(): Promise<TrendingCoin[]> {
  try {
    const response = await axios.get(`${import.meta.env.VITE_API_URL}/trending`)
    if (!response) {
      throw new Error('Failed to fetch trending coins')
    }
    let allCoins = []
    for (let coin of response.data) {
      allCoins.push({
        id: coin.ID,
        name: coin.Name,
        symbol: coin.Symbol
      })
    } 
    return allCoins }
  catch (error) {
    console.error('Error fetching trending coins:', error)
    return []
  }
}

export async function fetchCoinSuggestions(query: string): Promise<TrendingCoin[]> {
  try {
    const response = await axios.get(`${import.meta.env.VITE_API_URL}/crypto/search?name=${query}`)
    if (!response) {
      throw new Error('Failed to fetch coin suggestions')
    }
    let allCoins = []
    for (let coin of response.data) {
      allCoins.push({
        id: coin.ID,
        name: coin.Name,
        symbol: coin.Symbol
      })
    }
    return allCoins
  } catch (error) {
    console.error('Error fetching coin suggestions:', error)
    return []
  }
}
// export async function fetchTrendingCoins(): Promise<TrendingCoin[]> {
//   await new Promise(resolve => setTimeout(resolve, 500))
//   return [
//     { id: 'bitcoin', name: 'Bitcoin', symbol: 'BTC' },
//     { id: 'ethereum', name: 'Ethereum', symbol: 'ETH' },
//     { id: 'cardano', name: 'Cardano', symbol: 'ADA' },
//     { id: 'dogecoin', name: 'Dogecoin', symbol: 'DOGE' },
//     { id: 'polkadot', name: 'Polkadot', symbol: 'DOT' },
//   ]
// }

// export async function fetchCoinSuggestions(query: string): Promise<TrendingCoin[]> {
//   await new Promise(resolve => setTimeout(resolve, 300))
//   const allCoins = [
//     { id: 'bitcoin', name: 'Bitcoin', symbol: 'BTC' },
//     { id: 'ethereum', name: 'Ethereum', symbol: 'ETH' },
//     { id: 'cardano', name: 'Cardano', symbol: 'ADA' },
//     { id: 'dogecoin', name: 'Dogecoin', symbol: 'DOGE' },
//     { id: 'polkadot', name: 'Polkadot', symbol: 'DOT' },
//     { id: 'ripple', name: 'Ripple', symbol: 'XRP' },
//     { id: 'litecoin', name: 'Litecoin', symbol: 'LTC' },
//   ]
//   return allCoins.filter(coin => 
//     coin.name.toLowerCase().includes(query.toLowerCase()) || 
//     coin.symbol.toLowerCase().includes(query.toLowerCase())
//   )
// }