Simple go program for smart anonymizing database
Requires list os colument in file tables.txt
Example file attached

To run start program put this line in terminal:
linux:
sh ./start.sh "username:password@protocol(address)/dbname"
windows:
go run main.go "username:password@protocol(address)/dbname"

Currently it supports only MySQL DB
