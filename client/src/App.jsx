import "./static/css/App.css";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { createContext, useState } from "react";
import { Home } from "./pages/Home";
import { Cart } from "./pages/Cart";
import { Login } from "./pages/Login";
import { Register } from "./pages/Register";
import { Admin } from "./pages/Admin";
import { Nav } from "./components/Nav";
import { ForgotPassword } from "./components/ForgotPassword";
import { ResetPassword } from "./components/ResetPassword";
import { Checkout } from "./components/Checkout";

export const LoginContext = createContext();
const client = new QueryClient({
  // defaultOptions: {
  //   queries: {
  //     refetchOnWindowFocus: false,
  //     staleTime: 1000 * 60 * 5, // 5 minutes
  //   },
  // },
});

function App() {
  const [auth, setAuth] = useState(false);
  const [admin, setAdmin] = useState(false);
  const [userDetails, setUserDetails] = useState({});

  return (
    <div className="App">
      <QueryClientProvider client={client}>
        <LoginContext.Provider
          value={{
            auth,
            setAuth,
            admin,
            setAdmin,
            userDetails,
            setUserDetails,
          }}
        >
          <Router>
            <Nav />
            <Routes>
              <Route exact path="/" element={<Home />} />
              <Route exact path="/cart" element={<Cart />} />
              <Route exact path="/login" element={<Login />} />
              <Route exact path="/register" element={<Register />} />
              <Route exact path="/admin" element={<Admin />} />
              <Route
                exact
                path="/forgot-password"
                element={<ForgotPassword />}
              />
              <Route exact path="/reset-password" element={<ResetPassword />} />
              <Route exact path="/checkout" element={<Checkout />} />
              <Route exact path="*" element={<h1>PAGE NOT FOUND!!</h1>} />
            </Routes>
          </Router>
        </LoginContext.Provider>
      </QueryClientProvider>
    </div>
  );
}

export default App;
