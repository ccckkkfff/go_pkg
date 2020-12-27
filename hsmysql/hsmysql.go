package hsmysql

import (
	"bytes"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	jsoniter "github.com/json-iterator/go"
	"main/pkg"
	"reflect"
	"time"
)

//-------------------------------------------------------
// DelCols
// ------------------------------------------------------
//Return Value: err
//Params:
//		1)table string:			table name
//		2)pk string: 			primary key
//delete column from table 10 times/per
func DelCols(db *sql.DB, table,pk string)error{
	//delSql := fmt.Sprintf("delete from %s order by %s limit 10",table,pk)
	//_,err := db.Exec(delSql)
	//return err
	return nil
}


//-------------------------------------------------------
// InsertCols
// ------------------------------------------------------
//Return Value: err
//Params:
//		1)db *sql.DB:			database
//		2)sql string: 			insert prefix
//		3)data interface{}      the values of insert
//insert columns
func InsertCols(db *sql.DB, sql string,data interface{})(error,string){
	if db == nil{
		return errors.New("db err!"),""
	}

	ReflectValue := reflect.ValueOf(data)
	if ReflectValue.Kind() == reflect.Ptr{
		ReflectValue = ReflectValue.Elem()
	}

	if ReflectValue.Kind() != reflect.Slice{
		return errors.New("InsertCols err cause reflect type not slice!"+ReflectValue.Kind().String()),""
	}

	slen := ReflectValue.Len()
	if slen == 0{
		return errors.New("InsertCols err cause without values!"),""
	}

	if ReflectValue.Index(0).Kind() != reflect.Struct{
		//fmt.Println(ReflectValue.Index(0).NumField())
		return errors.New("InsertCols err cause reflect type not Struct!"),""
	}

	buffer := new(bytes.Buffer)
	buffer.Write([]byte(sql))

	ilen := ReflectValue.Index(0).NumField()
	for i:=0;i<slen;i++{
		if i == 0{
			buffer.WriteString("(")
		}else{
			buffer.WriteString(",(")
		}

		for j:=0;j<ilen;j++{
			if j != 0{
				buffer.WriteString(",")
			}

			v := ReflectValue.Index(i).Field(j)
			switch k := v.Kind();k{
			case reflect.String:
				buffer.WriteString("\""+v.String()+"\"")
			case reflect.Struct:
				t := v.Interface().(time.Time)
				buffer.WriteString("\""+t.Format("2006-01-02 15:04:05")+"\"")
				//return errors.New("InsertCols err cause have not reflect type!"),""
			default:
				return errors.New("InsertCols err cause "+ v.Type().String()+" cannot understand!"),""
			}
		}
		buffer.WriteString(")")
	}

	_,err := db.Exec(buffer.String())
	return err,buffer.String()
}



//-------------------------------------------------------
// GetDBVal
// ------------------------------------------------------
//Return Value: string
//Params:
//		1)db *sql.DB:			database
//		2)sql string: 			select sql
//insert columns
func GetDBVal(db *sql.DB,sql string)(string){
	rows,err := db.Query(sql)
	if err != nil{
		return ""
	}
	defer rows.Close()

	var v = make([]*pkg.RedisCache,0,10)
	for {
		var t = new(pkg.RedisCache)
		rows.Scan(&t.User)
		if t.User != ""{
			v = append(v,t)
		}

		if !rows.Next(){
			break
		}
	}
	cache,_ := jsoniter.Marshal(&v)

	return string(cache)
}