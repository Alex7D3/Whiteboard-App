import { useActionState } from "react";
import { Link, useLocation, useNavigate } from "react-router";
import { useAuth } from "../context/auth";

export default function Login() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const returnTo = location.state?.from?.pathname || "/";

  async function submit(_: string | null, formInput: FormData) {
    try {
      await login(formInput); 
      navigate(returnTo, { replace: true });
      return null;
    } catch (err) {
      return (err instanceof Error)
        ? err.message
        : "An error has occured. Please try again later";
    }
    
  }

  const [formError, action, pending] = useActionState(submit, null);
  
  return (
    <div className="flex flex-col items-center justify-center min-w-full min-h-screen">
      <h1 className="text-3xl text-blue">Welcome back!</h1>
      <form action={action} className="flex flex-col md:w-1/5 ">
        <input className="p-3 mt-8 rounded-md border border-grey focus:outline-none focus:border-blue"
          placeholder="email"
          autoComplete="email"
          name="email"
          type="email"
          required 
        />

        <input className="p-3 mt-4 rounded-md border border-grey focus:outline-none focus:border-blue"
          placeholder="password"
          name="password"
          type="password"
          minLength={6}
          required
        />

        <button 
          className="p-3 mt-4 rounded-md bg-blue font-bold text-white hover:opacity-80" type="submit"
          disabled={pending}
        >
          {pending ? "Loading..." : "Sign In"}
        </button>
        {formError && <span className="text-red ">{formError}</span>}

      </form>

      <Link
        to="/signup"
        state={{ from: location.state?.from }}
        className="p-3 mt-4 text-blue"
      >
        Sign Up
      </Link>
    </div>
  );
}
