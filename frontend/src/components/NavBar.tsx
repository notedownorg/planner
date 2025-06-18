import { Button } from "@/components/ui/button"
import { useTheme } from "@/components/ThemeProvider"
import { Moon, Sun } from "lucide-react"

export function NavBar() {
  const { theme, setTheme } = useTheme()

  return (
    <nav className="w-full border-b bg-background">
      <div className="flex h-16 items-center px-4 justify-between">
        <div>{/* Left side content */}</div>
        <Button
          variant="ghost"
          size="icon"
          onClick={() => setTheme(theme === "light" ? "dark" : "light")}
        >
          <Sun className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
          <Moon className="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
          <span className="sr-only">Toggle theme</span>
        </Button>
      </div>
    </nav>
  )
}