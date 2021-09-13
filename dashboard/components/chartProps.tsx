// We need to define the props in a separate file from the chart due to some typescript "bug"
// https://github.com/vercel/next.js/issues/22278
export interface BarChartProps {
    data: Array<Bar>
    singleTooltip: string
    multipleTooltip: string
}

export interface Bar {
    label: string
    value: number
}
