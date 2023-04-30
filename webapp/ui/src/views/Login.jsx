import React, { useState } from "react";
import { AuthContext } from "../context/authContext";
import { useContext } from "react";
import { Link, useNavigate } from "react-router-dom";
import axios from "axios";
import ReCAPTCHA from "react-google-recaptcha";



const Login = () => {
    const [inputs, setInputs] = useState({
        username: "",
        password: "",
    });
    const navigate = useNavigate();
    const { login } = useContext(AuthContext);
    const [err, setError] = useState(null);
    const [captchaResponse, setCaptchaResponse] = useState(null);

    const handleLogin = async (e) => {
        e.preventDefault();
        if (captchaResponse) {
            try {
              await login(inputs);
              navigate("/");
            } catch (err) {
              setError(err.response.data);
            }
        }
    };

    const handleChange = (e) => {
        setInputs((prev) => ({ ...prev, [e.target.name]: e.target.value }));
    };

    const onCaptchaChange = (response) => {
        setCaptchaResponse(response);
      };

    return (
        <div className="authenticate">
            <h1>Welcome to our UPS server!</h1>
            <h1>Login Page</h1>
            <form>
                <input
                    required
                    type="text"
                    placeholder="username"
                    name="username"
                    onChange={handleChange}
                />
                <input
                    required
                    type="password"
                    placeholder="password"
                    name="password"
                    onChange={handleChange}
                />
                <ReCAPTCHA sitekey="6Led-MslAAAAALlQDKLhLg5VDPURffbIUjhtjk1f" onChange={onCaptchaChange} />
                <button onClick={handleLogin}>Login</button>
                {err && <p>{err}</p>}
                <span>
                    Register here if you don't have an account! <Link to="/register">Register</Link>
                </span>
            </form>
        </div>
    );
};

export default Login;
