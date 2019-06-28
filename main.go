package main

import (
	_ "github.com/koebeltw/LineSlotCreator/ls/create"
	"bufio"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

//func init(){
//	SQLroot := "dcslot0913:5tgbnhy6@tcp(210.242.208.216:6360)/dc_mobile?charset=utf8"
//	x, err := xorm.NewEngine("mysql", SQLroot)
//
//	_, err = x.Query("SET @Id=0;")
//	if err != nil {
//		log.Fatalf("Fail to sync database: %v\n", err)
//	}
//
//	_, err = x.Query("SET @TempToken:='';")
//	if err != nil {
//		log.Fatalf("Fail to sync database: %v\n", err)
//	}
//
//
//	a, _ := x.Query(
//		`SELECT *
//FROM
//(
//	SELECT ROUND(MAX(OrderId) * RAND()) AS RandomId, A.Token
//	FROM(
//		SELECT a.*, @TempToken, @Id := CASE WHEN (@TempToken <> Token) THEN 0 ELSE (@Id + 1)  END OrderId, @TempToken:=Token
//		FROM role_push_new a
//		WHERE Token <> ""
//		ORDER BY Token
//	)A
//	GROUP BY Token
//)A
//LEFT JOIN
//(
//	SELECT a.*, @TempToken, @Id := CASE WHEN (@TempToken <> Token) THEN 0 ELSE (@Id + 1)  END OrderId, @TempToken:=Token
//	FROM role_push_new a
//	WHERE Token <> ""
//	ORDER BY Token
//)B ON A.RandomId = B.OrderId AND A.Token = B.Token;`)
//	if err != nil {
//		log.Fatalf("Fail to sync database: %v\n", err)
//	}
//	fmt.Println(a)
//}

func main() {
	for {
		r := bufio.NewReader(os.Stdin)
		b, _, _ := r.ReadLine()
		log.Println(string(b))
		switch {
		case string(b) == "send":
			//Client.Send(1, 1, []byte{1})
		}
	}

	//defer close(resultCh)
	//defer close(tempCh)
	//defer Engine.Close()
	//
	//fullRunWheel()
	//
	//resultWG.Wait()
	//tempWG.Wait()
	//saveResult(results[:saveResultCount])
	//saveTemp_Win(temps_Win[:saveTemp_WinCount])
	//saveTemp_NoWin(temps_NoWin[:saveTemp_NoWinCount])
	//
	//o := make([]*OddsCount, 0, len(OddsArray))
	//for _, Value := range OddsArray {
	//	o = append(o, Value)
	//}
	//
	//if _, err := Engine.Insert(&o); err != nil {
	//	log.Fatalf("Fail to sync database: %v\n", err)
	//}
	//
	//tempLineOddsCount := make([]LineOddsCount, 0)
	//for i := 0; i < len(LineOddsCountArray); i++ {
	//	for j := 0; j < len(LineOddsCountArray[i]); j++ {
	//		for k := 0; k < len(LineOddsCountArray[i][j]); k++ {
	//			tempLineOddsCount = append(tempLineOddsCount, LineOddsCountArray[i][j][k])
	//		}
	//	}
	//}
	//if _, err := Engine.Insert(&tempLineOddsCount); err != nil {
	//	log.Fatalf("Fail to sync database: %v\n", err)
	//}
}
