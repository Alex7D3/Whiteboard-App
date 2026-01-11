"use client";

import { useActionState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import { API_URL } from "@/src/constants/url";

export default function Register() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const returnTo = searchParams.get("returnTo") || "/";

  async function submit(_: string | null, formInput: FormData) {
    if (formInput.get("password") !== formInput.get("passwordconfirm")) {
      return "passwords do not match";
    }
    formInput.delete("passwordconfirm");
    const payload = Object.fromEntries(formInput);
    const res = await fetch(`${API_URL}/api/signup`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    });
    const data = await res.json();

    if (!res.ok) {
      return data.error;
    }
      
    router.push(returnTo);
    return null;
    
  }

  const [formError, action, pending] = useActionState(submit, null);
  
  return (
    <div className="flex flex-col items-center justify-center min-w-full min-h-screen">
      <h1 className="text-3xl text-blue">Welcome!</h1>
      <form action={action} className="flex flex-col md:w-1/5 ">
        <input className="p-3 mt-8 rounded-md border border-grey focus:outline-none focus:border-blue"
          placeholder="username"
          name="username"
          type="text"
          minLength={6}
          required 
        />
        <input className="p-3 mt-4 rounded-md border border-grey focus:outline-none focus:border-blue"
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

      <Link href={{ pathname: "/login", query: returnTo ? { returnTo } : {}}} className="p-3 mt-4 text-blue">
        Sign In
      </Link>
    </div>
  );
}