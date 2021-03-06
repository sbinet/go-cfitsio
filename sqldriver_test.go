package cfitsio

import (
	"database/sql"
	"testing"
)

func TestSqlDriver(t *testing.T) {
	db, err := sql.Open("fits", "testdata/file001.fits")
	if err != nil {
		t.Fatalf("error preparing sql abstraction for FITS: %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err == nil {
		t.Fatalf("expected ping to fail\n")
	}

	db, err = sql.Open("fits", "testdata/file001.fits[1]")
	if err != nil {
		t.Fatal("error opening fits file: %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("error pinging fits file: %v\n", err)
	}

	// var id int
	// var name string
	// rows, err := db.Query("select id, name from users where id = ?", 1)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	err := rows.Scan(&id, &name)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	fmt.Printf("id=%v, name=%q\n", id, name)
	// }
	// err = rows.Err()
	// if err != nil {
	// 	t.Fatal(err)
	// }

}
