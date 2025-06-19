export const config = {
    Config: {
        createFrom: (data) => data,
    },
}

export const habits = {
    WeeklyHabits: class WeeklyHabits {
        constructor(data) {
            Object.assign(this, data)
        }

        convertValues(a) {
            return a
        }
    },
    Habit: class Habit {
        constructor(data) {
            Object.assign(this, data)
        }
    },
}
