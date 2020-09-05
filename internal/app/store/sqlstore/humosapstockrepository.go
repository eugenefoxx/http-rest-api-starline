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
	"strconv"
	"strings"
)

// HUMO ...
type HUMO struct {
	ID        string `json:"id"`
	Package   string `json:"package"`
	Plant     string `json:"plant"`
	Warehouse string `json:"warehouse"`
	Material  string `json:"material"`
	Quantity  string `json:"quantity"`
	Pcs       string `json:"pcs"`
	Lot       string `json:"lot"`
	Status1   string `json:"status1"`
	Status2   string `json:"status2"`
	SPP       string `json:"spp"`
}

// HUMOSAPStockRepository ...
type HUMOSAPStockRepository struct {
	store *Store
}

// ImportDate ...
func (r *HUMOSAPStockRepository) ImportDate() {
	humounzip()
	//unzipCmd := exec.Command("bash", "-c", "find /home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/out.csv* -exec gunzip -c {} > /home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/out.csv")
	//_, errunzipCmd := unzipCmd.Output()
	//if errunzipCmd != nil {
	//	panic(errunzipCmd)
	//}

	// humodeletezero()
	humodeletezero1Cmd := exec.Command("bash", "-c", "sed 's/00000000000//g' /home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/out.csv > /home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/output1.csv")
	_, errhumodeletezero1Cmd := humodeletezero1Cmd.Output()
	if errhumodeletezero1Cmd != nil {
		panic(errhumodeletezero1Cmd)
	}
	humodeletezero2Cmd := exec.Command("bash", "-c", "sed 's/0000000000//g' /home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/output1.csv > /home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/output2.csv")
	_, errhumodeletezero2Cmd := humodeletezero2Cmd.Output()
	if errhumodeletezero1Cmd != nil {
		panic(errhumodeletezero2Cmd)
	}
	//	time.Sleep(2 * time.Second)
	replacement()

	// удаляем пустые строки
	lsCmd := exec.Command("bash", "-c", "sed '/^$/d' /home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/good_out_id_sap.csv > /home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/2good_out_id_sap.csv")
	_, errCmd := lsCmd.Output()
	if errCmd != nil {
		panic(errCmd)
	}

	//	clearzero()
	_, err := r.store.db.Exec(`TRUNCATE TABLE id_sap_import`)
	if err != nil {
		panic(err)
	}
	println("truncate table id_sap_import")
	cmd := "psql"
	file := "/home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/2good_out_id_sap.csv"
	args := []string{"-U", "postgres", "-d", "starline", "-c", fmt.Sprintf(`\copy id_sap_import from '%s' delimiter ';';`, file)}
	// вариант с чтением шапки и подачи через командную строку
	//args := []string{"-U", "postgres", "-d", "starline", "-c", fmt.Sprintf(`\copy id_sap_import from '%s' delimiter ',' csv header;`, os.Args[1])}
	v, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		panic(string(v))
	}
	println("Import date to id_sap_import Ok")

}

func humounzip() {
	// разархивация файла humo
	cmdbash := exec.Command("/bin/sh", "-c", "/home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/unzip.sh")
	//	cmd := exec.Command("Z:/Program Files/rentgen.exe")
	//	log.Printf("Running command and waiting for it to finish...")
	_, err := cmdbash.Output()
	if err != nil {
		log.Fatal(err)
	}
}

func humodeletezero() {
	cmdbash := exec.Command("/bin/sh", "-c", "/home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/humodeletezero.sh")
	//	cmd := exec.Command("Z:/Program Files/rentgen.exe")
	//	log.Printf("Running command and waiting for it to finish...")
	_, err := cmdbash.Output()
	if err != nil {
		log.Fatal(err)
	}
}

func replacement() {
	badCSVfile, err := os.Open("/home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/output2.csv")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer badCSVfile.Close()

	reader := bufio.NewReader(badCSVfile)
	scanner := bufio.NewScanner(reader)

	count()
	ss := count()
	fmt.Println(ss)
	//	lineCount := 0

	//	for scanner.Scan() {
	//		lineCount++
	//	}
	//	fmt.Println("number of lines:", lineCount-3)

	goodCSVfile, err := os.Create("/home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/good_out_id_sap.csv")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Writing to good_out_id_sap.csv")

	defer goodCSVfile.Close()

	for scanner.Scan() {
		//	lineCount++

		//	fmt.Println("linescount -", c)
		//	str := " rows selected."
		//	c := lineCount
		str := fmt.Sprintf("%s rows selected.", ss)
		// clean extra comma
		cleaned := strings.TrimSuffix(scanner.Text(), str)
		cleaned = strings.TrimSuffix(cleaned, str)

		//	c := lineCount + "rows selected."
		//	select := c + "rows selected."
		// replace double or triple comma
		cleaned = strings.Replace(cleaned, str, " ", -1)
		//	cleaned = strings.Replace(cleaned, strconv.Itoa(lineCount-3), " ", -1)
		//fmt.Println(cleaned)
		//	fmt.Println("fff-", str)
		_, err := io.WriteString(goodCSVfile, cleaned+"\n")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	//	fmt.Println("number of lines:", lineCount-3)
	//	fmt.Println(str)
	fmt.Println("Кол-во записей в HUMO", ss)

}

func count() string {
	file, _ := os.Open("/home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/output2.csv")
	fileScanner := bufio.NewScanner(file)
	defer file.Close()
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	fmt.Println("number of lines:", lineCount)
	c := strconv.Itoa(lineCount - 3)
	return c
}

func clearzero() {
	csvFile, err := os.Open("/home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/good_out_id_sap.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.LazyQuotes = true
	reader.Comma = ';'
	if err != nil {
		//	fmt.Println("Ошибка", err)
		fmt.Println(err)
		os.Exit(1)
	}
	defer csvFile.Close()
	var humo []HUMO
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		humo = append(humo, HUMO{
			ID:        line[0],
			Package:   line[1],
			Plant:     line[2],
			Warehouse: line[3],
			Material:  line[4],
			Quantity:  line[5],
			Pcs:       line[6],
			Lot:       line[7],
			Status1:   line[8],
			Status2:   line[9],
			SPP:       line[10],
		})
	}
	humoJSON, _ := json.Marshal(humo)

	err = json.Unmarshal([]byte(humoJSON), &humo)
	if err != nil {
		fmt.Println(err)
	}

	csvdatafile, err := os.Create("/home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/clear_good_out_id_sap.csv")

	if err != nil {
		fmt.Println(err)
	}
	defer csvdatafile.Close()
	writer := csv.NewWriter(csvdatafile)
	writer.Comma = ';'
	for _, worker := range humo {
		var record []string

		rID := strings.NewReplacer("0000000000", "")
		//	rIDD := strconv.Itoa(rID)
		record = append(record, rID.Replace(worker.ID))
		rPackage := strings.NewReplacer("00000000000", "")
		record = append(record, rPackage.Replace(worker.Package))
		record = append(record, worker.Plant)
		record = append(record, worker.Warehouse)
		rMaterial := strings.NewReplacer("00000000000", "")
		record = append(record, rMaterial.Replace(worker.Material))
		record = append(record, worker.Quantity)
		record = append(record, worker.Pcs)
		rLot := strings.NewReplacer("0000000000", "0")
		record = append(record, rLot.Replace(worker.Lot))
		record = append(record, worker.Status1)
		record = append(record, worker.Status2)
		record = append(record, worker.SPP)

		writer.Write(record)
	}

	// remember to flush!
	writer.Flush()
}
