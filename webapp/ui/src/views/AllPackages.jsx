import React, { useState, useEffect } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { useContext } from "react";
import { AuthContext } from "../context/authContext";

const AllPackages = () => {
    const [allPackages, setAllPackages] = useState(null);
    const navigate = useNavigate();
    const { currentUser } = useContext(AuthContext);

    useEffect(() => {
        const fetchAllPackages = async () => {
            try {
                const response = await axios.get("/package/all");
                setAllPackages(response.data);
            } catch (error) {
                console.error("Error fetching all package details:", error);
                navigate("/");
            }
        };

        fetchAllPackages();
    }, [navigate]);

    return (
        <div className="AllPackages">
            {allPackages ? (
                <ul className="package-list">
                    {allPackages.map((pkg) => (
                        <li key={pkg.trackID}>
                            <h2>{pkg.trackID}</h2>
                            <p><span className="label">PackageID:</span> {pkg.packageID}</p>
                            <p><span className="label">Status:</span> {pkg.status}</p>
                        </li>
                    ))}
                </ul>
            ) : (
                <p className="loading">Loading package details...</p>
            )}
        </div>
    );
};

export default AllPackages;
