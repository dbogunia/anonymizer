Simple go program for smart anonymizing database

Requires list of columns in file tables.txt

Example file attached

To run program put this line in terminal:

linux:
sh ./start.sh "username:password@protocol(address)/dbname"

windows:
go run main.go "username:password@protocol(address)/dbname"

Currently it supports only MySQL DB
