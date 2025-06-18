import { Button } from "@/components/ui/button"
import { NavBar } from "@/components/NavBar"
import React from "react"

function App() {
  const [count, setCount] = React.useState(0)

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <NavBar />
      <div className="flex-1 grid place-items-center mx-auto py-8">
        <div className="text-foreground text-2xl font-bold flex flex-col items-center space-y-4">
          <h1>Vite + React + TS + Tailwind + shadcn/ui</h1>
          <Button onClick={() => setCount(count + 1)}>Count up ({count})</Button>
        </div>
      </div>
    </div>
  )
}

export default App
