import { AuthForm } from '@/components/auth-form'

export const metadata = {
  title: 'Login - Gym Dashboard',
  description: 'Sign in to your gym dashboard',
}

export default function LoginPage() {
  return <AuthForm mode="login" />
}
