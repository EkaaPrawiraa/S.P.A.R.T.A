import { AuthForm } from '@/components/auth-form'

export const metadata = {
  title: 'Sign Up - Gym Dashboard',
  description: 'Create your gym dashboard account',
}

export default function RegisterPage() {
  return <AuthForm mode="register" />
}
