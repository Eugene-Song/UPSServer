package ups

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func (u *UPS) updatePackageTable(packageMetaData *PackageMetaData) {

	db := u.DB
	query := `
		INSERT INTO package (packageID, status, currentX, currentY, destinationX, destinationY, username, date)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
			status = VALUES(status),
			currentX = VALUES(currentX),
			currentY = VALUES(currentY),
			destinationX = VALUES(destinationX),
			destinationY = VALUES(destinationY),
			username = VALUES(username),
			date = NOW()
	`

	packageID := packageMetaData.packageId
	status := packageMetaData.status
	currentX := packageMetaData.currX
	currentY := packageMetaData.currY
	username := packageMetaData.username
	destinationX := packageMetaData.destX
	destinationY := packageMetaData.destY

	result, err := db.Exec(query, packageID, status, currentX, currentY, destinationX, destinationY, username)
	if err != nil {
		log.Fatal(err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully inserted or updated row with ID %d. Rows affected: %d\n", packageID, affectedRows)

}
