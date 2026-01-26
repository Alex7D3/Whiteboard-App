import { useAuth } from "../context/auth";

export default function RoomView() {
  const { user } = useAuth();
  return (
    <div>hello {user?.username}</div>
  );
}
