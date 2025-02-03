import { useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import axios from "axios";

export const ResetPassword = () => {
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const [searchParams] = useSearchParams(); // Extract query parameters
  const navigate = useNavigate(); // Hook to navigate to other routes
  const token = searchParams.get("token"); // Get the token from the URL

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      // Send the token and new password to the backend
      const response = await axios.post(
        "https://scrunchy.harshthakur.site/api/users/reset",
        {
          token,
          password,
        }
      );
      setMessage(response.data); // Display success message

      // Redirect to the login page after 3 seconds
      setTimeout(() => {
        navigate("/login");
      }, 2000);
    } catch (error) {
      setMessage(error.response?.data || "An error occurred"); // Display error message
    }
  };

  return (
    <div>
      <h2>Reset Password</h2>
      {token ? (
        <form onSubmit={handleSubmit}>
          <input
            type="password"
            placeholder="Enter new password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          <button type="submit">Reset Password</button>
        </form>
      ) : (
        <p>Invalid or missing token</p>
      )}
      {message && <p>{message}</p>}
    </div>
  );
};
