import React, { useState, useEffect } from "react";
import axios from "axios";
import { useNavigate, useParams } from "react-router-dom";

const SinglePackage = () => {
  const [packageDetails, setPackageDetails] = useState(null);
  const { id } = useParams();
  const navigate = useNavigate();

  useEffect(() => {
    const fetchPackageDetails = async () => {
      try {
        console.log("trackID", id);
        const response = await axios.get(`/package/${id}`);
        setPackageDetails(response.data);
      } catch (error) {
        console.error("Error fetching package details:", error);
        navigate("/package");
      }
    };

    fetchPackageDetails();
  }, [id, navigate]);

  return (
    <div>
      {packageDetails ? (
        <div className="Single">
          <h1>Package Details</h1>
          <p>Track ID: {packageDetails.packageID}</p>
          <p>Status: {packageDetails.status}</p>
        </div>
      ) : (
        <p className="loading">Loading package details...</p>
      )}
    </div>
  );
};

export default SinglePackage;
