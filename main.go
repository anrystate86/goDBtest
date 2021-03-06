package main

import (
	"database/sql"
	"flag"
	"fmt"
	"regexp"
	"strings"

	_ "github.com/godror/godror"
	_ "github.com/lib/pq"
)

func testConnstrinng(tconnstr string) (string, string, string) {
	reOra := regexp.MustCompile(`^jdbc:oracle:thin:@[\w\-\.]+:\d+\/[\w\d\_]+`)
	rePg := regexp.MustCompile(`^jdbc:postgresql:\/\/[\w\-\.]+\/[\w\d\_]+\?:[\w\d\_]+`)
	if reOra.Match([]byte(tconnstr)) {
		return "godror", strings.TrimPrefix(tconnstr, "jdbc:oracle:thin:@"), ""
	} else if rePg.Match([]byte(tconnstr)) {
		tr1 := strings.TrimPrefix(tconnstr, "jdbc:postgresql://")
		serverStr := strings.Split(tr1, "/")[0]
		portStr := strings.Split(strings.Split(tr1, "/")[1], ":")[1]
		dbStr := strings.Split(strings.Split(tr1, "/")[1], "?")[0]
		shemaStr := strings.Split(tr1, "/")[2]
		return "postgres", fmt.Sprintf("host=%s port=%s dbname=%s sslmode=disable search_path=%s", serverStr, portStr, dbStr, shemaStr), dbStr
	} else {
		return "error", "Error in parsing connection string", ""
	}
}

func checkDB(ctypeDB, cconStr, cuser, cpass, cdb string) string {
	type responce struct {
		tresp string
	}
	var cquerry string
	var cconnString string
	//querry
	if ctypeDB == "godror" {
		cquerry = "select 1 from dual;"
		cconnString = fmt.Sprintf("%s/%s%s", cuser, cpass, cconStr)
	} else if ctypeDB == "postgres" {
		cquerry = fmt.Sprintf("SELECT 1 as test FROM pg_database WHERE datname='%s';", cdb)
		cconnString = fmt.Sprintf("user=%s password=%s %s", cuser, cpass, cconStr)
	} else {
		return "error"
	}

	db, err := sql.Open(ctypeDB, cconnString)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer db.Close()

	var test responce
	err = db.QueryRow(cquerry).Scan(&test.tresp)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if test.tresp == "1" {
		return "OK"
	}
	fmt.Println(err)
	return ""
}

func main() {
	// Получение параметров запуска утилиты
	//----------------------------------------------------------
	userPtr := flag.String("user", "usernamevar", "username")
	passPtr := flag.String("pass", "passwordvar", "password")
	constrPtr := flag.String("constr", "connection string", "connection string for DB")
	flag.Parse()
	// Проверка добавления параметров запуска утилиты
	if (*userPtr != "usernamevar") && (*passPtr != "passwordvar") && (*constrPtr != "connection string") {
		// Проверка строки подключения
		typeDB, connectionString, dbName := testConnstrinng(*constrPtr)
		if typeDB != "error" {
			fmt.Println(checkDB(typeDB, connectionString, *userPtr, *passPtr, dbName))
		} else {
			fmt.Println(connectionString)
		}
	} else {
		fmt.Println("wrong usage, must testConnDB -user USERNAME -pass PASSWORD -constr CONNECTION_STRING_FOR_DB")
	}

}
