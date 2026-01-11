"use client";

import { useActionState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import { API_URL } from "@/src/constants/url";

export default function Login() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const returnTo = searchParams.get("returnTo") || "/";

  async function submit(_: string | null, formInput: FormData) {
    const payload = Object.fromEntries(formInput);
    const res = await fetch(`${API_URL}/api/login`, {
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

      <Link href={{ pathname: "/signup", query: returnTo ? { returnTo } : {}}} className="p-3 mt-4 text-blue">
        Sign Up
      </Link>
    </div>
  );
}
