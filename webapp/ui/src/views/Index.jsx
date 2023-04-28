import React from "react";
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import axios from "axios";
import { AuthContext } from "../context/authContext";
import { useContext } from "react";

const Index = () => {
    const [inputs, setInputs] = useState({
        trackID: "",
      });
      const [err, setError] = useState(null);
    
      const navigate = useNavigate();

      const { currentUser } = useContext(AuthContext);
    
      const handleChange = (e) => {
        setInputs((prev) => ({ ...prev, [e.target.name]: e.target.value }));
      };
    
      const handleSubmit = async (e) => {
        e.preventDefault();
        try {
          navigate(`/package/${inputs.trackID}`);
        } catch (err) {
          setError(err.response.data);
        }
      };

      const handleALLSubmit = async (e) => {
        e.preventDefault();
        try {
          navigate(`/package/all`);
        } catch (err) {
          setError(err.response.data);
        }
      };
    
      return (
        <div className="auth">
          <h1>Track you package here!</h1>
          <form>
            <input
              required
              type="text"
              placeholder="trackID"
              name="trackID"
              onChange={handleChange}
            />
            <button onClick={handleSubmit}>Track your package</button>
            {err && <p>{err}</p>}
          </form>
          
            {currentUser && (
                <button onClick={handleALLSubmit}>All your package</button>
            )}
          
          
        </div>
    );
}

export default Index;