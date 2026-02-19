import { AuthForm } from "@/components/auth-form";

export const metadata = {
  title: "Sign Up - S.P.A.R.T.A",
  description: "Create your S.P.A.R.T.A account",
};

export default function RegisterPage() {
  return <AuthForm mode="register" />;
}
