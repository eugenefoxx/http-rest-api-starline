package sqlstore

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

// MB52SAPStockRepository ...
type MB52SAPStockRepository struct {
	store *Store
}

// MB52 ...
type MB52 struct {
	Material    string `json:"Material"`
	Plant       string `json:"Plant"`
	Warehouse   string `json:"Warehouse"`
	Description string `json:"Description"`
	FreeStock   string `json:"FreeStock"`
	Inspection  string `json:"Inspection"`
	Block       string `json:"Block"`
	Lot         string `json:"Lot"`
	SPP         string `json:"SPP"`
}

// ImportDate ...
func (r *MB52SAPStockRepository) ImportDate() {
	creatfloatvalue()

	_, err := r.store.db.Exec(`TRUNCATE TABLE sap_stock_import`)
	if err != nil {
		panic(err)
	}
	println("truncate table sap_stock_import")
	cmd := "psql"
	file := "/home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/MB52_out.csv"
	args := []string{"-U", "postgres", "-d", "starline", "-c", fmt.Sprintf(`\copy sap_stock_import from '%s' delimiter '|' csv header;`, file)}
	// вариант с чтением шапки и подачи через командную строку
	//args := []string{"-U", "postgres", "-d", "starline", "-c", fmt.Sprintf(`\copy id_sap_import from '%s' delimiter ',' csv header;`, os.Args[1])}
	v, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		panic(string(v))
	}
	println("Import date to sap_stock_import Ok")
}

func creatfloatvalue() {
	csvFile, err := os.Open("/home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/export.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.LazyQuotes = true
	reader.Comma = '|'
	if err != nil {
		//	fmt.Println("Ошибка", err)
		fmt.Println(err)
		os.Exit(1)
	}
	defer csvFile.Close()
	var dataMB52 []MB52
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		dataMB52 = append(dataMB52, MB52{
			Material:    line[0],
			Plant:       line[1],
			Warehouse:   line[2],
			Description: line[3],
			FreeStock:   line[4],
			Inspection:  line[5],
			Block:       line[6],
			Lot:         line[7],
			SPP:         line[8],
		})
	}
	dataMB52JSON, _ := json.Marshal(dataMB52)

	err = json.Unmarshal([]byte(dataMB52JSON), &dataMB52)
	if err != nil {
		fmt.Println(err)
	}

	csvdatafile, err := os.Create("/home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/MB52_out.csv")

	if err != nil {
		fmt.Println(err)
	}
	defer csvdatafile.Close()

	writer := csv.NewWriter(csvdatafile)
	writer.Comma = '|'
	// add title
	//	writer.Write([]string{"Material", "Plant", "Warehouse", "Description", "FreeStock", "Inspection", "Block", "Lot", "SPP"})
	for _, worker := range dataMB52 {
		var record []string

		record = append(record, worker.Material)
		record = append(record, worker.Plant)
		record = append(record, worker.Warehouse)
		record = append(record, worker.Description)
		rFreeStock := strings.NewReplacer(",", ".")
		record = append(record, rFreeStock.Replace(worker.FreeStock))
		rInspection := strings.NewReplacer(",", ".")
		record = append(record, rInspection.Replace(worker.Inspection))
		rBlock := strings.NewReplacer(",", ".")
		record = append(record, rBlock.Replace(worker.Block))
		record = append(record, worker.Lot)
		record = append(record, worker.SPP)

		writer.Write(record)
	}

	// remember to flush!
	writer.Flush()
}
