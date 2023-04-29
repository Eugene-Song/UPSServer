package ups

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func (u *UPS) updatePackageTable(packageMetaData *PackageMetaData) {
	db := u.DB

	packageID := packageMetaData.PackageId
	status := packageMetaData.Status
	currentX := packageMetaData.currX
	currentY := packageMetaData.currY
	username := packageMetaData.username
	destinationX := packageMetaData.DestX
	destinationY := packageMetaData.DestY

	var query string
	if username != "" {
		query = `
		INSERT INTO package (packageID, status, currentX, currentY, destinationX, destinationY, username, date)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW()) AS new_values
		ON DUPLICATE KEY UPDATE
			status = new_values.status,
			currentX = new_values.currentX,
			currentY = new_values.currentY,
			destinationX = new_values.destinationX,
			destinationY = new_values.destinationY,
			username = new_values.username;
	`
	} else {
		query = `
		INSERT INTO package (packageID, status, currentX, currentY, destinationX, destinationY, username, date)
		VALUES (?, ?, ?, ?, ?, ?, NULL, NOW()) AS new_values
		ON DUPLICATE KEY UPDATE
			status = new_values.status,
			currentX = new_values.currentX,
			currentY = new_values.currentY,
			destinationX = new_values.destinationX,
			destinationY = new_values.destinationY,
	`
	}

	result, err := db.Exec(query, packageID, status, currentX, currentY, destinationX, destinationY, username)

	if err != nil {
		log.Fatal(err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully inserted or updated row with ID %d. Rows affected: %d\n", packageID, affectedRows)

}
