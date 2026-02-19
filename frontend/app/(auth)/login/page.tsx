import { AuthForm } from "@/components/auth-form";

export const metadata = {
  title: "Login - S.P.A.R.T.A",
  description: "Sign in to S.P.A.R.T.A",
};

export default function LoginPage() {
  return <AuthForm mode="login" />;
}
