import { useContext, useState, useEffect } from "react";
import logo from "../static/images/logo.png";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCartShopping } from "@fortawesome/free-solid-svg-icons";
import { LoginContext } from "../App";
import useCheckCookie from "./useCheckCookie";
import { useNavigate, Link } from "react-router-dom";
import axios from "axios";

export const Nav = () => {
  const { admin, setAdmin } = useContext(LoginContext);
  const { setUserDetails } = useContext(LoginContext);
  const navigate = useNavigate();
  const { cookieExists, refreshCookie } = useCheckCookie();
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  const onLogout = async () => {
    if (isLoggingOut) return;
    try {
      setIsLoggingOut(true);
      await axios.get("https://scrunchy.harshthakur.site/api/auth/logout", {
        withCredentials: true,
        timeout: 5000, // Ensure request doesn't hang
      });
      setAdmin(false);
      setUserDetails({});
    } catch (err) {
      console.log(err);
    } finally {
      // Always perform client-side cleanup
      await refreshCookie();
      navigate("/login");
      setIsLoggingOut(false);
    }
  };

  const fetchAdminStatus = async () => {
    if (cookieExists) {
      try {
        await axios.get("https://scrunchy.harshthakur.site/api/admin/isadmin", {
          withCredentials: true,
        });
        setAdmin(true);
      } catch (err) {
        console.log(err);
      }
    } else {
      setAdmin(false);
    }
  };

  const fetchUserDetails = async () => {
    if (cookieExists) {
      try {
        const res = await axios.get(
          "https://scrunchy.harshthakur.site/api/users/validate",
          {
            withCredentials: true,
          }
        );

        setUserDetails(res.data);
      } catch (err) {
        console.log(err);
      }
    } else {
      setUserDetails({});
      console.log("Cannot fetch user details");
    }
  };

  useEffect(() => {
    fetchAdminStatus();
    fetchUserDetails();
  }, [cookieExists]);

  return (
    <header className="Nav">
      <nav>
        <Link to="/">
          <span id="logo">
            <img src={logo} alt="Logo" />
          </span>
        </Link>
      </nav>
      <div className="nav-info">
        <Link to="/cart">
          <FontAwesomeIcon icon={faCartShopping} />
        </Link>
        {admin && (
          <Link to="/admin">
            <span>Admin</span>
          </Link>
        )}
        {!cookieExists ? (
          <>
            <Link to="/login" className="nav-link" aria-label="Login">
              Login
            </Link>
            <Link to="/register" className="nav-link" aria-label="Register">
              Register
            </Link>
          </>
        ) : (
          <button
            onClick={onLogout}
            className="logout-button"
            disabled={isLoggingOut}
            aria-label={isLoggingOut ? "Logging out..." : "Logout"}
          >
            {isLoggingOut ? "Logging out..." : "Logout"}
          </button>
        )}
      </div>
    </header>
  );
};
