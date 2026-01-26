import { Navigate, Outlet, useLocation } from "react-router";
import { useAuth } from "../context/auth";
import { Spinner } from "./Spinner";

export default function PrivateRoute() {
  const { isAuthenticated, isLoading } = useAuth();
  const location = useLocation();
  return (
    isLoading
      ? <Spinner/>
      : isAuthenticated
      ? <Outlet />
      : <Navigate to="/login" state={{ from: location.pathname }} replace />
  );
}
