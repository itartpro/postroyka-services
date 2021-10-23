package dbops

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"

	"go.mods/hashing"
)

type User struct {
	Id         int32     `json:"id"`
	Password   string    `json:"password"`
	Refresh    []string  `json:"refresh"`
	Created    time.Time `json:"created"`
	LastOnline time.Time `json:"last_online"`
	Rating     int16     `json:"rating"`
	//cant really change above stuff (except password)
	Login        string `json:"login"`
	Level        int16  `json:"level"`
	Avatar       bool   `json:"avatar"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PaternalName string `json:"paternal_name"`
	About        string `json:"about"`
	Balance      int32  `json:"balance"`
	TownId       int32  `json:"town_id"`
	RegionId     int16  `json:"region_id"`
	Legal        int16  `json:"legal"`
	Company      int16  `json:"company"`
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

type Choice struct {
	Id        int32 `json:"id"`
	LoginId   int32 `json:"login_id"`
	ServiceId int32 `json:"service_id"`
	Price     int32 `json:"price"`
	Parent    bool  `json:"parent"`
}

type Comment struct {
	Id          int32  `json:"id"`
	MasterId    int32  `json:"master_id"`
	ClientId    int32  `json:"client_id"`
	OrderId     int32  `json:"order_id"`
	ClientName  string `json:"client_name"`
	Politeness  int16  `json:"politeness"`
	Punctuality int16  `json:"punctuality"`
	Speed       int16  `json:"speed"`
	Balance     int16  `json:"balance"`
	Overall     int16  `json:"overall"`
	Text        string `json:"text"`
}

type cell struct {
	Id     int32  `json:"id"`
	Column string `json:"column"`
	Value  string `json:"value"`
	Table  string `json:"table"`
}

type PortfolioWork struct {
	Id          int32  `json:"id"`
	LoginId     int32  `json:"login_id"`
	ServiceId   int32  `json:"service_id"`
	OrderId     int32  `json:"order_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Volume      string `json:"volume"`
	Price       string `json:"price"`
}

type Order struct {
	Id          int32  `json:"id"`
	LoginId     int32  `json:"login_id"`
	ServiceId   int32  `json:"service_id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	RegionId  	int16  `json:"region_id"`
	TownId		int32  `json:"town_id"`
	Budget      int32  `json:"budget"`
	Created  time.Time `json:"created"`
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

func TryRegister(u User) (string, error) {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	//check for users with matching email OR phone
	var dup User
	if u.Email != "" {
		_ = pgxscan.Get(ctx, conn, &dup, `SELECT * FROM logins WHERE email = $1`, u.Email)
		if dup.Id > 0 {
			err = errors.New(u.Email + " is taken")
			return "", err
		}
	}

	if u.Phone != "" {
		_ = pgxscan.Get(ctx, conn, &dup, `SELECT * FROM logins WHERE phone = $1`, u.Phone)
		if dup.Id > 0 {
			err = errors.New(u.Phone + " is taken")
			return "", err
		}
	}

	row := conn.QueryRow(ctx, "INSERT INTO logins (password, created, email, phone, first_name, last_name, paternal_name, last_online, town_id, region_id, legal, level)"+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id",
		u.Password, u.Created, u.Email, u.Phone, u.FirstName, u.LastName, u.PaternalName, u.LastOnline, u.TownId, u.RegionId, u.Legal, u.Level)

	var id int32
	if err = row.Scan(&id); err != nil {
		return "", err
	}
	u.Id = id

	u.Password = ""
	jm, err := json.Marshal(u)

	return string(jm), nil
}

func UpdateLogin(u User) error {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer conn.Close()

	//check for users with matching email OR phone
	var dup User
	if u.Email != "" {
		_ = pgxscan.Get(ctx, conn, &dup, `SELECT * FROM logins WHERE email = $1 AND id != $2`, u.Email, u.Id)
		if dup.Id > 0 {
			err = errors.New(u.Email + " is taken")
			return err
		}
	}

	if u.Phone != "" {
		_ = pgxscan.Get(ctx, conn, &dup, `SELECT * FROM logins WHERE phone = $1 AND id != $2`, u.Phone, u.Id)
		if dup.Id > 0 {
			err = errors.New(u.Phone + " is taken")
			return err
		}
	}

	ct, err := conn.Exec(ctx, `UPDATE logins SET login = $1, level = $2, avatar = $3, email = $4, phone = $5, first_name = $6, last_name = $7, paternal_name = $8, about = $9, balance = $10, town_id = $11, region_id = $12, legal = $13, company = $14 WHERE id = $15`,
		u.Login, u.Level, u.Avatar, u.Email, u.Phone, u.FirstName, u.LastName, u.PaternalName, u.About, u.Balance, u.TownId, u.RegionId, u.Legal, u.Company, u.Id)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		err = errors.New(`"no rows updated"`)
		return err
	}

	return nil
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


func UpdateServiceChoices(news []Choice) error {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer conn.Close()

	//First get all old choices in case something needs to be deleted
	var old []Choice
	err = pgxscan.Select(ctx, conn, &old, `SELECT * FROM choices WHERE login_id = $1 AND parent = true`, news[0].LoginId)
	if err != nil {
		return err
	}

	if len(news) >= len(old) {
		for i := range news {
			if i < len(old) {
				_, err = conn.Exec(ctx, `UPDATE choices SET service_id = $1 WHERE id = $2`, news[i].ServiceId, old[i].Id)
				if err != nil {
					return err
				}
			} else {
				_, err = conn.Exec(ctx, `INSERT INTO choices (login_id, service_id, parent) VALUES ($1, $2, true)`, news[i].LoginId, news[i].ServiceId)
				if err != nil {
					return err
				}
			}
		}
	}

	if len(old) > len(news) {
		for i := range old {
			if i < len(news) {
				_, err = conn.Exec(ctx, `UPDATE choices SET service_id = $1 WHERE id = $2`, news[i].ServiceId, old[i].Id)
				if err != nil {
					return err
				}
			} else {
				_, err = conn.Exec(ctx, `DELETE FROM choices WHERE id = $1`, old[i].Id)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func UpdateServicePrices(news []Choice) error {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer conn.Close()

	//First get all old choices in case something needs to be deleted
	var old []Choice
	err = pgxscan.Select(ctx, conn, &old, `SELECT * FROM choices WHERE login_id = $1 AND parent = false`, news[0].LoginId)
	if err != nil {
		return err
	}

	if len(news) >= len(old) {
		for i := range news {
			if i < len(old) {
				_, err = conn.Exec(ctx, `UPDATE choices SET service_id = $1, price = $2 WHERE id = $3`, news[i].ServiceId, news[i].Price, old[i].Id)
				if err != nil {
					return err
				}
			} else {
				_, err = conn.Exec(ctx, `INSERT INTO choices (login_id, service_id, price, parent) VALUES ($1, $2, $3, false)`, news[i].LoginId, news[i].ServiceId, news[i].Price)
				if err != nil {
					return err
				}
			}
		}
	}

	if len(old) > len(news) {
		for i := range old {
			if i < len(news) {
				_, err = conn.Exec(ctx, `UPDATE choices SET service_id = $1, price = $2 WHERE id = $3`, news[i].ServiceId, news[i].Price, old[i].Id)
				if err != nil {
					return err
				}
			} else {
				_, err = conn.Exec(ctx, `DELETE FROM choices WHERE id = $1`, old[i].Id)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func GetMastersChoices(id int32) ([]Choice, error) {
	var cs []Choice

	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return cs, err
	}
	defer conn.Close()

	err = pgxscan.Select(ctx, conn, &cs, `SELECT * FROM choices WHERE login_id = $1`, id)
	if err != nil {
		return cs, err
	}

	return cs, nil
}


func UpdateCell(instructions string) error {
	var c cell

	err := json.Unmarshal([]byte(instructions), &c)
	if err != nil {
		return err
	}

	if c.Table != "logins" {
		err = errors.New("access denied")
		return err
	}

	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer conn.Close()

	ct, err := conn.Exec(ctx, `Update `+c.Table+` SET `+c.Column+` = $1 WHERE id = $2`, c.Value, c.Id)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		err = errors.New("no rows found")
		return err
	}

	return nil
}


func AddWork(instructions string) error {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer conn.Close()

	var w PortfolioWork

	err = json.Unmarshal([]byte(instructions), &w)
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, `INSERT INTO portfolio (login_id, order_id, name, service_id, description, volume, price) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		w.LoginId, w.OrderId, w.Name, w.ServiceId, w.Description, w.Volume, w.Price)
	if err != nil {
		return err
	}

	return nil
}

func UpdateWork(instructions string) error {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer conn.Close()

	var w PortfolioWork

	err = json.Unmarshal([]byte(instructions), &w)
	if err != nil {
		return err
	}

	ct, err := conn.Exec(ctx,
		`UPDATE portfolio SET order_id = $1, name = $2, service_id = $3, description = $4, volume = $5, price = $6 WHERE id = $7`,
		w.OrderId, w.Name, w.ServiceId, w.Description, w.Volume, w.Price, w.Id)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		err = errors.New(`"no rows updated"`)
		return err
	}

	return nil
}

func GetPortfolio(instructions string) (string, error) {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	var ins PortfolioWork
	err = json.Unmarshal([]byte(instructions), &ins)
	if err != nil {
		return "", err
	}

	var ws []PortfolioWork

	err = pgxscan.Select(ctx, conn, &ws, `SELECT * FROM portfolio WHERE login_id = $1`, ins.LoginId)
	if err != nil {
		return "", err
	}

	jm, err := json.Marshal(ws)
	if err != nil {
		return "", err
	}

	return string(jm), nil
}

func GetProfileComments(id int32) ([]Comment, error) {
	var cs []Comment

	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return cs, err
	}
	defer conn.Close()

	err = pgxscan.Select(ctx, conn, &cs, `SELECT * FROM comments WHERE master_id = $1`, id)
	if err != nil {
		return cs, err
	}

	return cs, nil
}


func AddOrder(instructions string) (string, error) {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	var o Order

	err = json.Unmarshal([]byte(instructions), &o)
	if err != nil {
		return "", err
	}

	sql := `INSERT INTO orders (login_id, service_id, name, title, description, region_id, town_id, budget, created) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	row := conn.QueryRow(ctx, sql, o.LoginId, o.ServiceId, o.Name, o.Title, o.Description, o.RegionId, o.TownId, o.Budget, o.Created)
	if err = row.Scan(&o.Id); err != nil {
		return "", err
	}

	jm, err := json.Marshal(o)

	return string(jm), nil
}
