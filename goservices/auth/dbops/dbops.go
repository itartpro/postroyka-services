package dbops

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"

	"go.mods/hashing"
)

type User struct {
	Id		  int32     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	FullName  string     `json:"full_name"`
	Gender    string    `json:"gender"`
	Birthdate time.Time `json:"birthdate"`
	Created   time.Time	`json:"created"`
	CountryId int16     `json:"country_id"`
	RegionId  int16		`json:"region_id"`
	TownId    int32     `json:"town_id"`
	Marital   string	`json:"marital"`
	Phone	  string	`json:"phone"`
	Site 	  string	`json:"site"`
	Level     int16		`json:"level"`
	Refresh   []string	`json:"refresh"`
	Avatar	  bool		`json:"avatar"`
	Weight    int16 	`json:"weight"`
	Height    int16		`json:"height"`
}

type Country struct {
	Id int16    `json:"id"`
	Name string `json:"name"`
}

type Region struct {
	Id 		  int32  `json:"id"`
	Name 	  string `json:"name"`
	CountryId int16  `json:"country_id"`
}

type Town struct {
	Id 		  int32  `json:"id"`
	Name 	  string `json:"name"`
	CountryId int16  `json:"country_id"`
	RegionId  int16  `json:"region_id"`
}

func TryLogin(login string, pwd string) (User, error) {

	var user User

	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return user, err
	}
	defer conn.Close()

	if err := pgxscan.Get(ctx, conn, &user, `SELECT * FROM logins WHERE email=$1 OR phone=$1`, login); err != nil {
		return user, err
	}

	if err := hashing.ValidatePassword([]byte(pwd), []byte(user.Password)); err != nil {
		return user, err
	}

	return user, nil
}

func GetProfile(u User) (User, error) {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil { return u, err }

	err = pgxscan.Get(ctx, conn, &u, `SELECT * FROM logins WHERE id = $1`, u.Id)
	if err != nil { return u, err }

	return u, nil
}

func TryRegister(u User) (User, error)  {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))

	//check for users with matching email OR phone
	var login string
	var dup User
	if u.Email != "" {
		login += u.Email
	} else {
		login += u.Phone
	}
	_ = pgxscan.Get(ctx, conn, &dup, `SELECT * FROM logins WHERE email=$1 OR phone=$1`, login)
	if dup.Id > 0 {
		err = errors.New("duplicate user")
		return u, err
	}

	row := conn.QueryRow(ctx, "INSERT INTO logins (email, phone, password, full_name, gender, birthdate, created, country_id)"+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		u.Email, u.Phone, u.Password, u.FullName, u.Gender, u.Birthdate, u.Created, u.CountryId)

	var id int32
	if err = row.Scan(&id); err != nil {
		return u, err
	}
	u.Id = id

	return u, nil
}

func ReadCountries() ([]Country, error) {
	var cs []Country

	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil { return cs, err }
	defer conn.Close()

	err = pgxscan.Select(ctx, conn, &cs, `SELECT * FROM countries`)
	if err != nil { return cs, err }

	return cs, nil
}

func ReadRegions(id int16) ([]Region, error) {
	var rs []Region

	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil { return rs, err }
	defer conn.Close()

	err = pgxscan.Select(ctx, conn, &rs, `SELECT * FROM regions WHERE country_id = $1`, id)
	if err != nil { return rs, err }

	return rs, nil
}

func ReadTowns(id int16) ([]Town, error) {
	var ts []Town

	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil { return ts, err }
	defer conn.Close()

	err = pgxscan.Select(ctx, conn, &ts, `SELECT * FROM towns WHERE region_id = $1`, id)
	if err != nil { return ts, err }

	return ts, nil
}

func NewCountry(c Country) (Country, error)  {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil { return c, err }
	defer conn.Close()

	row := conn.QueryRow(ctx, "INSERT INTO countries (name) VALUES ($1) RETURNING id", c.Name)

	var id int16
	if err = row.Scan(&id); err != nil {
		return c, err
	}
	c.Id = id

	return c, nil
}

func UpdateRefresh(id string, hash string) error {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer conn.Close()

	var user User

	rows, err := conn.Query(ctx, `SELECT * FROM logins WHERE id=$1`, id)
	if err != nil {
		return err
	}

	if err := pgxscan.ScanOne(&user, rows); err != nil {
		return err
	}

	s := user.Refresh
	//limit to 5 refreshes/devices, if its more than 4, remove first one
	if len(s) > 4 {
		s = s[1:]
	}
	user.Refresh = append(s, hash)

	ct, err := conn.Exec(ctx, `UPDATE logins SET refresh = $1 WHERE id = $2`, user.Refresh, user.Id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() != 1 {
		print("No row found to update")
	}

	return nil
}

func TryRefresh(id string, hash string) (User, error) {

	var user User

	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return user, err
	}
	defer conn.Close()

	rows, err := conn.Query(ctx, `SELECT * FROM logins WHERE id=$1`, id)
	if err != nil {
		return user, err
	}
	defer rows.Close()

	if err := pgxscan.ScanOne(&user, rows); err != nil {
		return user, err
	}

	for _, v := range user.Refresh {
		if v == hash {
			return user, nil
		}
	}
	err = errors.New("refresh token not found in database")
	return user, err
}
