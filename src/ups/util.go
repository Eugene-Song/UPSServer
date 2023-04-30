package ups

import (
	"fmt"
	"os"
)

func (u *UPS) PrintUPS() error {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "UPS State:")

	fmt.Fprintf(file, "SimSpeed: %d\n", u.SimSpeed)

	u.PackageMutex.Lock()
	fmt.Fprintln(file, "Package:")
	for k, v := range u.Package {
		fmt.Fprintf(file, "  Key: %d, Value: %v\n", k, v)
	}
	u.PackageMutex.Unlock()

	fmt.Fprintln(file, "UnAckedPickup:")
	for k, v := range u.UnAckedPickup {
		fmt.Fprintf(file, "  Key: %d, Value: %v\n", k, v)
	}

	fmt.Fprintln(file, "UnAckedDeliver:")
	for k, v := range u.UnAckedDeliver {
		fmt.Fprintf(file, "  Key: %d, Value: %v\n", k, v)
	}

	u.TruckMutex.Lock()
	fmt.Fprintln(file, "Truck:")
	for k, v := range u.Truck {
		fmt.Fprintf(file, "  Key: %d, Value: %s\n", k, v)
	}
	u.TruckMutex.Unlock()

	u.MapTruckShipMutex.Lock()
	fmt.Fprintln(file, "MapTruckShip:")
	for k, v := range u.MapTruckShip {
		fmt.Fprintf(file, "  Key: %d, Value: %v\n", k, v)
	}
	u.MapTruckShipMutex.Unlock()

	return nil
}
