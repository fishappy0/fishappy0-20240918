export interface CoinData {
  name: string
  symbol: string
  price: number
  rank: number
  supply: number
  marketCap: number
  conversions: {
    SGD: number
    VND: number
    MYR: number
    CNY: number
    EUR: number
  }
  chartData: Array<{
    x: number
    y: number[]
    price: number
  }>
}

export interface DetailedCoinApiResult {
    name: string
    symbol: string
    price: number
    rank: number
    supply: number
    marketCap: number
}

export interface TrendingCoin {
  id: string
  name: string
  symbol: string
}