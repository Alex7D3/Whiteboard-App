import { useActionState } from "react";
import { Link, useLocation, useNavigate } from "react-router";
import { useAuth } from "../context/auth";

export default function Register() {
  const { register } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const returnTo = location.state?.from?.pathname || "/";

  async function submit(_: string | null, formInput: FormData) {
    if (formInput.get("password") !== formInput.get("passwordconfirm")) {
      return "passwords do not match";
    }
    formInput.delete("passwordconfirm");
    try {
      await register(formInput);
      navigate("/login", { replace: true, state: returnTo });
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
      <h1 className="text-3xl text-blue">Register your account.</h1>
      <form action={action} className="flex flex-col md:w-1/5 ">
        <input className="p-3 mt-8 rounded-md border border-grey focus:outline-none focus:border-blue"
          placeholder="username"
          name="username"
          type="text"
          required
        />

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

        <input className="p-3 mt-4 rounded-md border border-grey focus:outline-none focus:border-blue"
          placeholder="confirm password"
          name="passwordconfirm"
          type="password"
          minLength={6}
          required
        />

        <button
          className="p-3 mt-4 rounded-md bg-blue font-bold text-white hover:opacity-80" type="submit"
          disabled={pending}
        >
          {pending ? "Loading..." : "Sign Up"}
        </button>
        {formError && <span className="text-red ">{formError}</span>}

      </form>

      <Link
        to="/signin"
        state={{ from: location.state?.from }}
        className="p-3 mt-4 text-blue"
      >
        Sign In
      </Link>
    </div>
  );
}
