# goDBtest
Test DB connection oracle\postgres from application parameters

Application for Zabbix monitoring, test application DB connection parameter.
Connection template:
jdbc:oracle:thin:@dbo-server:1521/apptest1
jdbc:postgresql://db-server1/dbbig1?:5432/apptest1

Using:
goDBtest -user User -pass Password -constr jdbc:postgresql://localhost/postgres?:5432/public
