package sqlstore

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

// PanacimStockRepository ...
type PanacimStockRepository struct {
	store *Store
}

// ImportDate ...
func (r *PanacimStockRepository) ImportDate() {

	rows := panaread()
	//	time.Sleep(4 * time.Second)
	writeChanges(rows)
	//	time.Sleep(4 * time.Second)
	_, err := r.store.db.Exec(`TRUNCATE TABLE panacim_stock`)
	if err != nil {
		panic(err)
	}
	println("truncate table panacim_stock")
	cmd := "psql"
	file := "/home/webserver/http-rest-api-starline/import_date/outputPanaCIM.csv"
	args := []string{"-U", "postgres", "-d", "starline", "-c", fmt.Sprintf(`\copy panacim_stock from '%s' delimiter ',' csv header;`, file)}
	v, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		panic(string(v))
	}
	println("Import date to panacim_stock Ok")

	//	return nil
}

func panaread() [][]string {
	f1, err := os.Open("/home/webserver/http-rest-api-starline/import_date/MaterialManagementReport.csv") // Копия.csv  google_report_source.tsv
	if err != nil {
		panic(err)
	}
	defer f1.Close()

	linesPana, err := readCsv(f1) //
	if err != nil {
		panic(err)
	}

	return linesPana
}

func writeChanges(rows [][]string) {
	f, err := os.Create("/home/webserver/http-rest-api-starline/import_date/outputPanaCIM.csv")
	if err != nil {
		log.Fatal(err)
	}
	err = csv.NewWriter(f).WriteAll(rows)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func readCsv(rs io.ReadSeeker) ([][]string, error) {
	// ReadCsv(filename string) ([][]string, error)

	// Open CSV file
	//	f, err := os.Open(filename)
	//	reader := csv.NewReader(bufio.NewReader(f))
	//	reader.LazyQuotes = true
	//	reader.Comma = ';'

	// Skip first row(line)
	row1, err := bufio.NewReader(rs).ReadSlice('\n')
	if err != nil {
		return nil, err
	}

	_, err = rs.Seek(int64(len(row1)), io.SeekStart)
	if err != nil {
		return nil, err
	}

	//	if err != nil {
	//		return [][]string{}, err
	//	}
	//	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(rs).ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return lines, nil

}
