import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: HomePage,
})

function HomePage() {
  return (
    <div className="container mx-auto p-8">
      <h1 className="text-4xl font-bold mb-4">Oauth Client</h1>
      <p className="text-lg">Welcome to Oauth Client SPA</p>
    </div>
  )
}
