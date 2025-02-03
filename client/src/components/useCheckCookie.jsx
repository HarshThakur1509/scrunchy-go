import { useState, useEffect } from "react";
import axios from "axios";

const useCheckCookie = () => {
  const [cookieExists, setCookieExists] = useState(false);
  const [loading, setLoading] = useState(true);

  const checkCookie = async () => {
    try {
      const response = await axios.get(
        "https://scrunchy.harshthakur.site/api/cookie",
        {
          withCredentials: true,
          timeout: 3000,
        }
      );
      setCookieExists(response.data.exists);
    } catch (error) {
      console.log(error);

      setCookieExists(false);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    checkCookie();
  }, []);

  return { cookieExists, loading, refreshCookie: checkCookie };
};

export default useCheckCookie;
