package dbops

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"go.mods/hashing"
)

type User struct {
	Id           int32     `json:"id"`
	Login        string    `json:"login"`
	Password     string    `json:"password"`
	Level        int16     `json:"level"`
	Refresh      []string  `json:"refresh"`
	Created      time.Time `json:"created"`
	Avatar       bool      `json:"avatar"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Rating       int16     `json:"rating"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	PaternalName string    `json:"paternal_name"`
	LastOnline   time.Time `json:"last_online"`
	About        string    `json:"about"`
	Balance      int32     `json:"balance"`
	TownId       int32     `json:"town_id"`
	Legal        int16     `json:"legal"`
}

type Country struct {
	Id   int16  `json:"id"`
	Name string `json:"name"`
}

type Region struct {
	Id        int32  `json:"id"`
	Name      string `json:"name"`
	CountryId int16  `json:"country_id"`
}

type Town struct {
	Id        int32  `json:"id"`
	Name      string `json:"name"`
	CountryId int16  `json:"country_id"`
	RegionId  int16  `json:"region_id"`
}

type ServiceChoice struct {
	Id		  int32 `json:"id"`
	LoginId   int32 `json:"login_id"`
	ServiceId int32 `json:"service_id"`
}

func TryLogin(login string, pwd string) (User, error) {

	var user User

	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return user, err
	}
	defer conn.Close()

	if err := pgxscan.Get(ctx, conn, &user, `SELECT * FROM logins WHERE login=$1 OR email=$1 OR phone=$1`, login); err != nil {
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
	if err != nil {
		return u, err
	}
	defer conn.Close()

	err = pgxscan.Get(ctx, conn, &u, `SELECT * FROM logins WHERE id = $1`, u.Id)
	if err != nil {
		return u, err
	}

	return u, nil
}

func TryRegister(u User) (User, error) {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return u, err
	}
	defer conn.Close()

	//check for users with matching email OR phone
	var dup User
	if u.Email != "" {
		_ = pgxscan.Get(ctx, conn, &dup, `SELECT * FROM logins WHERE email=$1`, u.Email)
		if dup.Id > 0 {
			err = errors.New("duplicate user")
			return u, err
		}
	}

	if u.Phone != "" {
		_ = pgxscan.Get(ctx, conn, &dup, `SELECT * FROM logins WHERE phone=$1`, u.Phone)
		if dup.Id > 0 {
			err = errors.New("duplicate user")
			return u, err
		}
	}

	row := conn.QueryRow(ctx, "INSERT INTO logins (password, created, email, phone, first_name, last_name, paternal_name, last_online, town_id, legal, level)"+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
		u.Password, u.Created, u.Email, u.Phone, u.FirstName, u.LastName, u.PaternalName, u.LastOnline, u.TownId, u.Legal)

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
	if err != nil {
		return cs, err
	}
	defer conn.Close()

	err = pgxscan.Select(ctx, conn, &cs, `SELECT * FROM countries`)
	if err != nil {
		return cs, err
	}

	return cs, nil
}

func ReadRegions(id int16) ([]Region, error) {
	var rs []Region

	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return rs, err
	}
	defer conn.Close()

	err = pgxscan.Select(ctx, conn, &rs, `SELECT * FROM regions WHERE country_id = $1`, id)
	if err != nil {
		return rs, err
	}

	return rs, nil
}

func ReadTowns(id int16) ([]Town, error) {
	var ts []Town

	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return ts, err
	}
	defer conn.Close()

	err = pgxscan.Select(ctx, conn, &ts, `SELECT * FROM towns WHERE region_id = $1`, id)
	if err != nil {
		return ts, err
	}

	return ts, nil
}

func NewCountry(c Country) (Country, error) {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return c, err
	}
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
	rows.Close()

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

func UpdateServiceChoices(choices []ServiceChoice) error {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {return err}
	defer conn.Close()

	//TODO update old ones, insert surplus new ones or delete excess old ones
	//var oldChoices []ServiceChoice
	//err = pgxscan.Select(ctx, conn, &oldChoices, `SELECT * FROM choices WHERE login_id = $1`, choices[0].LoginId)
	//if err != nil {return err}

	//INSERT MULTIPLE ROWS
	var inputRows [][]interface{}
	for _, v := range choices {
		inputRows = append(inputRows, []interface{}{
			v.LoginId,
			v.ServiceId,
		})
	}
	copyCount, err := conn.CopyFrom(ctx, pgx.Identifier{"choices"}, []string{"login_id", "service_id"}, pgx.CopyFromRows(inputRows))
	if err != nil {
		err = errors.New("Unexpected error for CopyFrom: "+err.Error())
		return err
	}
	if int(copyCount) != len(inputRows) {
		err = errors.New("Expected CopyFrom to return "+strconv.Itoa(len(inputRows))+" copied rows, but got  "+strconv.Itoa(int(copyCount)))
		return err
	}

	return nil
}
