import React, { useContext } from "react";
import { Link } from "react-router-dom";
import { AuthContext } from "../context/authContext";
import Logo from "../img/ups.png";

const Navbar = () => {
    const { currentUser, logout } = useContext(AuthContext);

    return (
        <div className="navbar">
            <div className="container">
                <div className="logo">
                    <Link to="/">
                        <img src={Logo} alt="" />
                    </Link>
                </div>
                <div className="links">
                    <span>{currentUser?.username}</span>
                    {currentUser ? (
                        <span onClick={logout}>Logout</span>
                    ) : (
                        <Link className="link" to="/login">
                            Login
                        </Link>
                    )}
                </div>
            </div>
        </div>
    );
};

export default Navbar;
