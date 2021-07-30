Simple go program for smart anonymizing database

Requires list of columns in file tables.txt

Example file attached

To run program put this line in terminal:

linux:
sh ./start.sh "username:password@protocol(address)/dbname"

windows:
go run main.go "username:password@protocol(address)/dbname"

Or simply use included Dockerfile

Currently it supports only MySQL DB

-------
Changing from updates on DB to generating SQL file

Comment lines from 101 to 110

add code for appengind updateQuery string to file:

f, err := os.OpenFile("output.sql", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	
if err != nil {

	log.Println(err)
	
}

defer f.Close()

if _, err := f.WriteString(updateQuery + ";\n"); err != nil {

	log.Println(err)
	
}


WARNING: THIS PIECE OF CODE WAS NOT TESTED - PLEASE TEST BEFORE USING IN REAL LIFE SCENARIO!!!
