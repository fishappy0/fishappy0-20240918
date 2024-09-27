import { CoinData } from "@/lib/types"
import { Suspense, lazy } from "react"

const ApexChart = lazy(() => import('react-apexcharts'))

export function CombinedChart({ data }: { data: CoinData['chartData'] }) {
  const options = {
    chart: {
      type: 'candlestick',
      height: 350
    },
    title: {
      text: 'Price and OHLC Chart',
      align: 'left'
    },
    xaxis: {
      type: 'datetime'
    },
    yaxis: [
      {
        tooltip: {
          enabled: true
        }
      },
      {
        opposite: true,
        tooltip: {
          enabled: true
        }
      }
    ],
    tooltip: {
      shared: true,
      custom: ({ dataPointIndex, w }: any) => {
        if (dataPointIndex === -1) return ''
        
        const candlestickData = w.globals.seriesCandleO[0][dataPointIndex] !== undefined
          ? {
              o: w.globals.seriesCandleO[0][dataPointIndex],
              h: w.globals.seriesCandleH[0][dataPointIndex],
              l: w.globals.seriesCandleL[0][dataPointIndex],
              c: w.globals.seriesCandleC[0][dataPointIndex]
            }
          : null

        const price = w.globals.series[1][dataPointIndex]
        const date = new Date(w.globals.seriesX[0][dataPointIndex])

        let tooltipContent = `<div class="apexcharts-tooltip-box">
          <div>Date: ${date.toLocaleDateString()}</div>`

        if (candlestickData) {
          tooltipContent += `
            <div>Open: <span class="value">${candlestickData.o.toFixed(2)}</span></div>
            <div>High: <span class="value">${candlestickData.h.toFixed(2)}</span></div>
            <div>Low: <span class="value">${candlestickData.l.toFixed(2)}</span></div>
            <div>Close: <span class="value">${candlestickData.c.toFixed(2)}</span></div>`
        }

        tooltipContent += `<div>Price: <span class="value">${price.toFixed(2)}</span></div>
        </div>`

        return tooltipContent
      }
    },
    legend: {
      show: false
    }
  }

  const series = [
    {
      name: 'OHLC',
      type: 'candlestick',
      data: data.map(item => ({
        x: item.x,
        y: item.y
      }))
    },
    {
      name: 'Price',
      type: 'line',
      data: data.map(item => ({
        x: item.x,
        y: item.price
      })),
      color: '#1E90FF'
    }
  ]

  return (
    <Suspense fallback={<div>Loading chart...</div>}>
      <ApexChart
        options={options}
        series={series}
        type="line"
        height={350}
      />
    </Suspense>
  )
}