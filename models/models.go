package models

import (
	"github.com/evil-router/isfired/database"
	uuid2 "github.com/satori/go.uuid"
	"log"
	"golang.org/x/net/idna"
)

type comment struct {
	Name     string
	Reason   string
	Time     string
	Location string
	Status   bool
}

type site struct {
	ID   string
	Name string
	Site string
}

func GetSite(name string) (string, error) {
	var id string
	log.Printf("get site %v", name)
	db, err := database.GetDB()
	if err != nil {
		log.Print(err)
		return "", err
	}
	err = db.QueryRow("select site.Site_ID  From site where Site = ?", name).Scan(&id)
	if err != nil {
		log.Print(err)
		return "", err
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
		n, _= idna.ToUnicode(n)
		m, _ = idna.ToUnicode(m)
		c = append(c, comment{n, m, t, l, s})
	}

	return c, nil
}

func SetComment(site string, comment string, city string, status bool) error {
	db, err := database.GetDB()
	if err != nil {
		log.Print(err)
		return err
	}
	site,_ = idna.ToASCII(site)
	id, err := GetSite(site)
	if err != nil {
		log.Print(err)
		return err
	}
	comment,_= idna.ToASCII(comment)

	rows, err := db.Query("INSERT INTO `Fired`.`comment` (`FK_Site_ID`, `message`, `time`, `location`, `Status`)"+
		"VALUES (?, ?, DEFAULT, ? , ?)", id, comment, city, status)
	defer rows.Close()
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func GetActiveSites() ([]site, error) {
	var sites []site
	db, err := database.GetDB()
	if err != nil {
		log.Print(err)
		return sites, err
	}
	rows, err := db.Query("Select Site,Name,PK_ID from site")
	defer rows.Close()
	if err != nil {
		log.Print(err)
		return sites, err
	}
	for rows.Next() {
		var s, n, i string
		err := rows.Scan(&s, &n, &i)
		if err != nil {
			log.Fatal(err)
		}
		n,_ = idna.ToUnicode(n)
		s,_= idna.ToUnicode(s)
		sites = append(sites, site{i, n, s})
	}

	return sites, nil
}

func AddSite(site, name string) error {
	db, err := database.GetDB()
	if err != nil {
		log.Print(err)
		return err
	}
	uuid, _ := uuid2.NewV4()
	rows, err := db.Query("INSERT INTO `Fired`.`site` (`Site_ID`, `Name`, `Site`, `PK_ID`)"+
		"VALUES (?, ?, ? , DEFAULT)", uuid.String(), name, site)
	defer rows.Close()
	if err != nil {
		log.Printf("Site Add %v", err)
		return err
	}
	SetComment(site, "Welcome", "xn--no8h", false)

	return nil
}
