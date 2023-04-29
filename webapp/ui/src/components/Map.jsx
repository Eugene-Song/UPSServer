import React, { useEffect, useState } from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';

const Map = () => {
    const [position, setPosition] = useState(null);

    const fetchPackageLocation = () => {
        fetch('/api/package-location')
            .then((res) => res.json())
            .then((data) => {
                setPosition([data.y, data.x]);
            });
    };

    useEffect(() => {
        fetchPackageLocation();

        // Set up an interval to fetch the package location every 5 seconds (5000 ms)
        const interval = setInterval(fetchPackageLocation, 5000);

        // Cleanup the interval when the component is unmounted
        return () => clearInterval(interval);
    }, []);

    if (!position) {
        return <div>Loading...</div>;
    }

    return (
        <div className="App">
            <MapContainer center={position} zoom={13} style={{ height: '100vh', width: '100%' }}>
                <TileLayer
                    attribution='&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                />
                <Marker position={position}>
                    <Popup>Package Location</Popup>
                </Marker>
            </MapContainer>
        </div>
    );
}

export default Map;







