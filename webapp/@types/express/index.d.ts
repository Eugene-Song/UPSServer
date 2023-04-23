export {}

declare global {
  namespace Express {
    interface User {
      email: string
        id: string
        name: string
        username: string
    }
  }
}