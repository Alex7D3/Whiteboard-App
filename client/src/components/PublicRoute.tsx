import { Navigate, Outlet } from "react-router";
import { useAuth } from "../context/auth";
import { Spinner } from "./Spinner";

export default function PublicRoute() {
  const { isAuthenticated, isLoading } = useAuth();
  return (
    isLoading
      ? <Spinner/>
      : isAuthenticated
      ? <Navigate to="/" replace />
      : <Outlet />
  );
}
