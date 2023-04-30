package ups

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func (u *UPS) UpdatePackageTable(packageMetaData *PackageMetaData) {
	log.Printf("Enter UpdatePackageTable function")
	db := u.DB

	packageID := packageMetaData.PackageId
	status := packageMetaData.Status
	currentX := packageMetaData.currX
	currentY := packageMetaData.currY
	username := packageMetaData.username
	destinationX := packageMetaData.DestX
	destinationY := packageMetaData.DestY
	item := packageMetaData.itemDetails

	var query string
	if username != "" {
		query = `
		INSERT INTO package (packageID, status, currentX, currentY, destinationX, destinationY, username, date, item)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), ?) AS new_values
		ON DUPLICATE KEY UPDATE
			status = new_values.status,
			currentX = new_values.currentX,
			currentY = new_values.currentY,
			destinationX = new_values.destinationX,
			destinationY = new_values.destinationY,
			username = new_values.username,
			item = new_values.item;
	`
	} else {
		query = `
		INSERT INTO package (packageID, status, currentX, currentY, destinationX, destinationY, username, date, item)
		VALUES (?, ?, ?, ?, ?, ?, NULL, NOW(), ?) AS new_values
		ON DUPLICATE KEY UPDATE
			status = new_values.status,
			currentX = new_values.currentX,
			currentY = new_values.currentY,
			destinationX = new_values.destinationX,
			destinationY = new_values.destinationY,
		    item = new_values.item;
	`
	}
	var result sql.Result
	var err error
	if username != "" {
		result, err = db.Exec(query, packageID, status, currentX, currentY, destinationX, destinationY, username, item)
	} else {
		result, err = db.Exec(query, packageID, status, currentX, currentY, destinationX, destinationY, item)
	}

	if err != nil {
		log.Fatal(err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully inserted or updated row with ID %d. Rows affected: %d\n", packageID, affectedRows)

}
