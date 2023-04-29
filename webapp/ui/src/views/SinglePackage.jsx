import React, { useState, useEffect } from "react";
import axios from "axios";
import { useNavigate, useParams, Link } from "react-router-dom";
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import { Icon } from 'leaflet';
import Logo from "../img/marker.png";
import HomeLogo from "../img/home.png"
import { useContext } from "react";
import { AuthContext } from "../context/authContext";

const customIcon1 = new Icon({
  iconUrl: Logo,
  iconSize: [40, 40],
  iconAnchor: [12.5, 41],
  popupAnchor: [0, -41],
});

const customIcon2 = new Icon({
  iconUrl: HomeLogo,
  iconSize: [50, 50],
  iconAnchor: [12.5, 41],
  popupAnchor: [0, -41],
});

const SinglePackage = () => {
  const [packageDetails, setPackageDetails] = useState(null);
  const [position, setPosition] = useState([51.505, -0.09]);
  const [home, setHome] = useState([0, 0]);
  const { id } = useParams();
  const navigate = useNavigate();
  const { currentUser } = useContext(AuthContext);
  const [inputs, setInputs] = useState({
    X: 0,
    Y: 0,
  });
  const [err, setError] = useState(null);
  

  const fetchPackageDetails = async () => {
    try {
      console.log("trackID", id);
      const response = await axios.get(`/package/${id}`);
      setPackageDetails(response.data);
      console.log("response.data", response.data);
      setPosition([51.505+response.data.currentY/1000, -0.09+response.data.currentX/1000]);
      setHome([51.505+response.data.destinationY/1000, -0.09+response.data.destinationX/1000]);

    } catch (error) {
      console.error("Error fetching package details:", error);
      navigate("/package");
    }
  };

  useEffect(() => {
    fetchPackageDetails();
  }, [id, navigate]);



  useEffect(() => {
      fetchPackageDetails();
      // Set up an interval to fetch the package location every 5 seconds (5000 ms)
      const interval = setInterval(fetchPackageDetails, 5000);
      // Cleanup the interval when the component is unmounted
      return () => clearInterval(interval);
  }, []);

  const handleChange = (e) => {
    setInputs((prev) => ({ ...prev, [e.target.name]: parseInt(e.target.value, 10) }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await axios.put(`/package/${packageDetails.shipID}`, inputs);
    } catch (err) {
      setError(err.response.data);
    }
  };

  return (
    <div>
      <div>
        {packageDetails ? (
          <div className="Single">
            <h1>Package Details</h1>
            <p>Track ID: {packageDetails.packageID}</p>
            <p>Status: {packageDetails.status}</p>
            {packageDetails.username && currentUser.username === packageDetails.username && (
            <div className="edit">
              <h1>Change your delivery address here!</h1>
              <form>
                <input
                  required
                  type="number"
                  placeholder="New address x"
                  name="X"
                  onChange={handleChange}
                />
                <input
                  required
                  type="number"
                  placeholder="New address y"
                  name="Y"
                  onChange={handleChange}
                />
                <button onClick={handleSubmit}>Update your address</button>
                {err && <p>{err}</p>}
              </form>
            </div>
            )}
          </div>
          
        ) : (
          <p className="loading">Loading package details...</p>
        )}
      </div>
      <div className="App">
            <MapContainer center={position} zoom={9} style={{ height: '100vh', width: '100%' }}>
                <TileLayer
                    attribution='&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                />
                <Marker position={position} icon={customIcon1}>
                    <Popup>Package Location</Popup>
                </Marker>
                <Marker position={home} icon={customIcon2}>
                    <Popup>Destination Location</Popup>
                </Marker>
            </MapContainer>
        </div>
    </div>
    
  );
};

export default SinglePackage;
