import { useForm } from "react-hook-form";
import * as yup from "yup";
import { yupResolver } from "@hookform/resolvers/yup";
import { useContext } from "react";
import { LoginContext } from "../App";
import { Navigate, Link } from "react-router-dom";
import axios from "axios";

export const Login = () => {
  const { auth, setAuth } = useContext(LoginContext);

  const onSubmit = async (formdata) => {
    const email = formdata.email;
    const password = formdata.password;

    try {
      await axios.post(
        "http://localhost:3000/users/login",
        { email, password },
        { withCredentials: true }
      );
      localStorage.setItem("login", "true");
      setAuth(true);
    } catch (err) {
      console.log(err);
    }
  };

  const schema = yup.object().shape({
    email: yup.string().required("email required"),
    password: yup.string().min(4).max(20).required(),
  });

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    resolver: yupResolver(schema),
  });

  if (auth) {
    return <Navigate to="/" />;
  }

  return (
    <div className="Login">
      <h1>Login</h1>
      <form className="Form" onSubmit={handleSubmit(onSubmit)}>
        <input type="text" placeholder="email..." {...register("email")} />
        <p>{errors.email?.message}</p>
        <input
          type="password"
          placeholder="Password..."
          {...register("password")}
        />

        <button className="btn3" type="submit">
          Submit
        </button>
      </form>
      <Link to="/forgot-password">Forgot Password</Link>
    </div>
  );
};
