package models

import (
	"github.com/evil-router/isfired/database"
	"log"
)

type comment struct {
	Name     string
	Reason   string
	Time     string
	Location string
	Status   bool
}

type site struct {
	ID   int64
	Name string
	Site string
}

func GetSite(name string) (int64, error) {
	var id int64
	log.Printf("get site %v", name)
	db, err := database.GetDB()
	if err != nil {
		log.Print(err)
		return 0, err
	}
	err = db.QueryRow("select Site_ID  From site where Site = ?", name).Scan(&id)
	if err != nil {
		log.Print(err)
		return 0, err
	}


	return id, nil
}

func GetComment(site string, count int64, offset int64) ([]comment, error) {
	var c []comment
	db, err := database.GetDB()
	if err != nil {
		log.Print(err)
		return c, err
	}
	rows, err := db.Query(`select  site.Name, comment.message, comment.Status,comment.time,comment.location   from comment ,site
	where comment.FK_Site_ID = site.Site_ID  AND site.Site = ?
	Order by comment.PK_comment_ID desc
	limit ? OFFSET ?;`, site, count, offset)
	if err != nil {
		log.Print(err)
		return c, err
	}
	defer rows.Close()
	for rows.Next() {
		var n, m, t, l string
		var s bool
		err := rows.Scan(&n, &m, &s, &t, &l)
		if err != nil {
			log.Fatal(err)
		}
		c = append(c, comment{n, m, t, l, s})
	}

	return c, nil
}

func SetComment(site string, comment string,city string, status bool,) ( error) {
	db, err := database.GetDB()
	if err != nil {
		log.Print(err)
		return  err
	}
	id,err :=GetSite(site)
	if err != nil {
		log.Print(err)
		return  err
	}

	rows, err := db.Query("INSERT INTO `Fired`.`comment` (`FK_Site_ID`, `message`, `time`, `location`, `Status`)"+
		"VALUES (?, ?, DEFAULT, ? , ?)", id,comment,city ,status)
	if err != nil {
		log.Print(err)
		return  err
	}
	defer rows.Close()
	return  nil
}
