package models

import (
	"github.com/evil-router/isfired/database"
	"log"
)

type comment struct {
	Name   string
	Reason string
	Time   string
	Location string
	Status bool

}

type site struct {
	ID   int64
	Name string
	Site string
}

func GetSite(name string) (site, error) {
	var s site
	db, err := database.GetDB()
	if err != nil {
		log.Print(err)
		return s, err
	}
	err = db.QueryRow("select PK_site_ID, Name,Site From site where PK_site_ID = ?", 1).Scan(&s.ID, &s.Name, &s.Site)
	if err != nil {
		log.Print(err)
		return s, err
	}

	return s, nil
}

func GetComment(site string, count int64, offset int64) ([]comment, error) {
	var c []comment
	db, err := database.GetDB()
	if err != nil {
		log.Print(err)
		return c, err
	}
	rows, err := db.Query(`select  site.Name, comment.message, comment.Status,comment.time,comment.location   from comment ,site
	where comment.FK_PK_site_ID = site.PK_site_ID  AND site.Site = ?
	limit ? OFFSET ?;`, site,count, offset)
	if err != nil {
		log.Print(err)
		return c, err
	}
	defer rows.Close()
	for rows.Next() {
		var n, m,t,l string
		var s bool
		err := rows.Scan(&n, &m, &s,&t,&l)
		if err != nil {
			log.Fatal(err)
		}
		c = append(c, comment{n, m,t,l,s})
	}

	return c, nil
}
