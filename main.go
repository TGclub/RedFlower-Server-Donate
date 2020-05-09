package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"log"
	"strconv"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func getLog(c *gin.Context) {
	fmt.Println("log")
	pgIndex := c.Query("pageIndex")
	pageIndex, err := strconv.Atoi(pgIndex)
	if err != nil {
		log.Println(err)
		c.JSON(400, "pageIndex is not a number")
		return
	}

	pgSize := c.Query("pageSize")
	pageSize, err := strconv.Atoi(pgSize)
	if err != nil {
		log.Println(err)
		c.JSON(400, "pageSize is not a number")
		return
	}
	fmt.Println(pageIndex)
	fmt.Println(pageSize)

	name := viper.GetString("mysql.name")
	pass := viper.GetString("mysql.pass")
	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")

	datasource := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8", name, pass, host, port)
	fmt.Println(datasource)
	db, err := sql.Open("mysql", datasource)
	sql := fmt.Sprintf("SELECT `donate_id`,`amount`,`name`,`phone` FROM `redflower`.`donate` ORDER BY `donate_id` DESC limit %d,%d", pageIndex*pageSize, pageSize)
	fmt.Println(sql)
	rows, err := db.Query(sql)

	maps := getMaps(rows)
	max_index := 0

	for k, v := range maps { //查询出来的数组
		if k > max_index {
			max_index = k
		}
		v["time"] = v["donate_id"][:10]
		delete(v, "donate_id")
	}
	var values []map[string]string = make([]map[string]string, max_index+1)
	for k, v := range maps {
		values[k] = v
	}

	sql = "select count(*) from `redflower`.`donate`"
	rows, err = db.Query(sql)
	maps = getMaps(rows)
	fmt.Println(maps[0])
	nums := 0
	for _, v := range maps[0] {
		nums, _ = strconv.Atoi(v)
	}
	fmt.Println(nums)

	page_num := nums/pageSize + 1
	fmt.Println(page_num)

	c.JSON(200, gin.H{"data": values, "total_num": nums, "page_num": page_num})
}
func main() {
	viper.SetConfigFile("./config.json")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("read config error", err)
	}

	port := viper.GetString("port")

	eng := gin.Default()
	eng.Use(CORSMiddleware())
	eng.GET("/donate", getLog)
	eng.Run(":" + port)

}
func getMaps(query *sql.Rows) map[int]map[string]string {
	column, _ := query.Columns()              //读出查询出的列字段名
	values := make([][]byte, len(column))     //values是每个列的值，这里获取到byte里
	scans := make([]interface{}, len(column)) //因为每次查询出来的列是不定长的，用len(column)定住当次查询的长度
	for i := range values {                   //让每一行数据都填充到[][]byte里面
		scans[i] = &values[i]
	}
	results := make(map[int]map[string]string) //最后得到的map
	i := 0
	for query.Next() { //循环，让游标往下移动
		if err := query.Scan(scans...); err != nil { //query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
			fmt.Println(err)
			return nil
		}
		row := make(map[string]string) //每行数据
		for k, v := range values {     //每行数据是放在values里面，现在把它挪到row里
			key := column[k]
			row[key] = string(v)
		}
		results[i] = row //装入结果集中
		i++
	}
	//for _, v := range results { //查询出来的数组
	//	fmt.Println(v)
	//}
	return results
}
