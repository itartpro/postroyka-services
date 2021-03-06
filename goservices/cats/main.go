package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"go.mods/grpcc"
)

type cat struct {
	Id          int32     `json:"id"`
	ParentId    int32     `json:"parent_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Keywords    string    `json:"keywords"`
	Author      string    `json:"author"`
	H1          string    `json:"h1"`
	Text        string    `json:"text"`
	Image       string    `json:"image"`
	SortOrder   int32     `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	Extra       string    `json:"extra"`
}

type catSummary struct {
	Id          int32     `json:"id"`
	ParentId    int32     `json:"parent_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Image       string    `json:"image"`
	SortOrder   int32     `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	Extra       string    `json:"extra"`
}

type cell struct {
	Id     int32  `json:"id"`
	Column string `json:"column"`
	Value  string `json:"value"`
}

var service = "cats"

func result(status string, data string) string {
	return `{"name":"` + service + `","status":` + status + `,"data":` + data + `}`
}

type server struct{}

func (*server) PassData(ctx context.Context, req *grpcc.DataRequest) (*grpcc.DataResponse, error) {

	var res grpcc.DataResponse
	res.Result = result("false", `"noop or error"`)

	instructions := req.GetData().GetInstructions()

	var c cat
	if err := json.Unmarshal([]byte(instructions), &c); err != nil {
		res.Result = result("false", service+" couldn't unmarshal instructions "+err.Error())
		return &res, err
	}

	op := req.GetData().GetAction()

	ctx = context.Background()
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return &res, err
	}
	defer conn.Close()

	summary := "id, parent_id, name, slug, sort_order, image, created_at, extra"

	if op == "create" {
		row := conn.QueryRow(ctx, "INSERT INTO cats (parent_id, name, slug, title, description, keywords, author, h1, text, image, sort_order, created_at, extra)"+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id",
			c.ParentId, c.Name, c.Slug, c.Title, c.Description, c.Keywords, c.Author, c.H1, c.Text, c.Image, c.SortOrder, c.CreatedAt, c.Extra)

		var id int32
		if err = row.Scan(&id); err != nil {
			return &res, err
		}
		c.Id = id
		c.SortOrder = id

		//check for duplicate slugs
		var dup cat
		_ = pgxscan.Get(ctx, conn, &dup, `SELECT * FROM cats WHERE slug=$1`, c.Slug)
		if dup.Id > 0 {
			c.Slug += "-" + strconv.Itoa(int(c.Id))
			if _, err = conn.Exec(ctx, `UPDATE cats SET sort_order = $1, slug = $2 WHERE id = $1`, id, c.Slug); err != nil {
				c.SortOrder = 0
			}
		} else {
			if _, err = conn.Exec(ctx, `UPDATE cats SET sort_order = $1 WHERE id = $1`, id); err != nil {
				c.SortOrder = 0
			}
		}

		b, err := json.Marshal(c)
		if err != nil {
			res.Result = result("false", `"insert success, marshal fail"`)
			return &res, nil
		}

		res.Result = result("true", string(b))
	}

	if op == "delete" {
		var c catSummary
		if err := json.Unmarshal([]byte(instructions), &c); err != nil {
			return &res, err
		}

		var cats []*catSummary
		if err = pgxscan.Select(ctx, conn, &cats, `SELECT id, parent_id FROM cats WHERE parent_id=$1`, c.Id); err != nil {
			return &res, err
		}

		if len(cats) > 0 {
			res.Result = result("false", `"delete children"`)
		} else {
			_, err = conn.Exec(ctx, "DELETE FROM cats WHERE id=$1", c.Id)
			if err != nil {
				return &res, err
			}

			if err = os.RemoveAll(os.Getenv("UPLOADS_DIR") + "cats/" + strconv.Itoa(int(c.Id))); err != nil {
				return &res, err
			}

			res.Result = result("true", `"cat deleted:`+strconv.Itoa(int(c.Id))+`"`)
		}

		return &res, nil
	}

	if op == "read" {
		if c.Id != 0 {
			if err := pgxscan.Get(ctx, conn, &c, `SELECT * FROM cats WHERE id=$1`, c.Id); err != nil {
				return &res, err
			}
		} else {
			if err := pgxscan.Get(ctx, conn, &c, `SELECT * FROM cats WHERE slug=$1`, c.Slug); err != nil {
				return &res, err
			}
		}

		b, err := json.Marshal(c)
		if err != nil {
			res.Result = result("false", `"read success, marshal fail"`)
			return &res, nil
		}

		res.Result = result("true", string(b))
	}

	if op == "read_all" {
		var cats []*catSummary
		if err = pgxscan.Select(ctx, conn, &cats, `SELECT `+summary+` FROM cats ORDER BY sort_order ASC`); err != nil {
			return &res, err
		}

		if len(cats) < 1 {
			res.Result = result("false", `"no rows found"`)
			return &res, nil
		}

		b, err := json.Marshal(cats)
		if err != nil {
			return &res, err
		}

		res.Result = result("true", string(b))
	}

	if op == "read-where-in" {
		whereIn := struct {
			Column    string   `json:"column"`
			Values 	  []string `json:"values"`
		}{}
		err := json.Unmarshal([]byte(instructions), &whereIn)
		if err != nil {
			return &res, err
		}

		var str string
		for _, v := range whereIn.Values {
			str += v + `,`
		}
		str = str[:len(str)-1] // remove last ","

		var cats []*catSummary
		sqlStr := `SELECT `+summary+` FROM cats WHERE `+whereIn.Column+` IN (`+str+`) ORDER BY sort_order ASC`
		err = pgxscan.Select(ctx, conn, &cats, sqlStr)
		if err != nil {
			return &res, err
		}

		if len(cats) < 1 {
			res.Result = result("false", `"no rows found"`)
			return &res, nil
		}

		b, err := json.Marshal(cats)
		if err != nil {
			return &res, err
		}

		res.Result = result("true", string(b))
	}

	if op == "update" {
		var c cat
		if err := json.Unmarshal([]byte(instructions), &c); err != nil {
			return &res, err
		}

		//check for duplicate slugs
		var dup cat
		_ = pgxscan.Get(ctx, conn, &dup, `SELECT * FROM cats WHERE slug=$1`, c.Slug)
		if dup.Id > 0 {
			c.Slug += "-" + strconv.Itoa(int(c.Id))
		}

		ct, err := conn.Exec(ctx, `UPDATE cats SET parent_id = $1, name = $2, slug = $3, title = $4, description = $5, keywords = $6, author = $7, h1 = $8, text = $9, image = $10, sort_order = $11, created_at = $12, extra = $13 WHERE id = $14`,
			c.ParentId, c.Name, c.Slug, c.Title, c.Description, c.Keywords, c.Author, c.H1, c.Text, c.Image, c.SortOrder, c.CreatedAt, c.Extra, c.Id)
		if err != nil {
			return &res, err
		}

		if ct.RowsAffected() == 0 {
			res.Result = result("false", `"no rows updated"`)
			return &res, nil
		}

		res.Result = result("true", `"updated successfully"`)
		return &res, nil
	}

	if op == "update-cell" {
		var c cell
		if err := json.Unmarshal([]byte(instructions), &c); err != nil {
			return &res, err
		}

		ct, err := conn.Exec(ctx, `Update cats SET `+c.Column+` = $1 WHERE id = $2`, c.Value, c.Id)
		if err != nil {
			log.Println(err)
			return &res, err
		}

		if ct.RowsAffected() == 0 {
			res.Result = result("false", `"no rows found"`)
			return &res, nil
		}

		res.Result = result("true", `"updated successfully"`)
		return &res, nil
	}

	return &res, nil
}

func main() {
	ok, err := credentials.NewServerTLSFromFile(os.Getenv("SERVICEKEY_PEM"), os.Getenv("SERVICEKEY_KEY"))
	if err != nil {
		log.Fatalf("Failed to setup TLS:%v", err)
	}

	lis, err := net.Listen("tcp", ":50004")
	if err != nil {
		log.Fatal(service + "service failed to listen ", err)
	}

	log.Println("Hi, I'm a " + service + " grpc comm. service listening...")

	s := grpc.NewServer(grpc.Creds(ok))
	grpcc.RegisterCommunicationServiceServer(s, &server{})
	err = s.Serve(lis)
	if err != nil {
		log.Fatal("Failed to serve grpc server " + service + ":", err)
	}
}