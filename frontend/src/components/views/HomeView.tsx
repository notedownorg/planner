import { Button } from '@/components/ui/button'
import React from 'react'

export function HomeView() {
    const [count, setCount] = React.useState(0)

    return (
        <div className="flex-1 grid place-items-center mx-auto py-8">
            <div className="text-foreground text-2xl font-bold flex flex-col items-center space-y-4">
                <h1>Vite + React + TS + Tailwind + shadcn/ui</h1>
                <Button onClick={() => setCount(count + 1)}>Count up ({count})</Button>
            </div>
        </div>
    )
}
