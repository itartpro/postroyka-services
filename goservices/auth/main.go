package main

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"

	"go.mods/dbops"
	"go.mods/grpcc"
	"go.mods/hashing"
)

var service = "auth"

func result(status string, data string) string {
	return `{"name":"` + service + `","status":` + status + `,"data":` + data + `}`
}

type server struct{}

//implement PassData interface from grpcc
func (*server) PassData(ctx context.Context, req *grpcc.DataRequest) (*grpcc.DataResponse, error) {

	var res grpcc.DataResponse
	res.Result = result("false", `"noop or error"`)

	instructions := req.GetData().GetInstructions()
	op := req.GetData().GetAction()

	//whatever
	if op == "update-cell" {
		err := dbops.UpdateCell(instructions)
		if err != nil {return &res, err}
		res.Result = result("true", instructions)
		return &res, nil
	}

	//login stuff
	if op == "register" {
		var user dbops.User
		err := json.Unmarshal([]byte(instructions), &user)
		if err != nil {
			return &res, err
		}

		user.Password = hashing.GeneratePassword(user.Password)
		userString, err := dbops.TryRegister(user)
		if err != nil {
			return &res, err
		}

		res.Result = result("true", userString)
		return &res, nil
	}

	if op == "login" {
		var in dbops.User
		err := json.Unmarshal([]byte(instructions), &in)
		if err != nil {
			return &res, err
		}

		login, pwd, err := hashing.B64DecodeTryUser(in.Login, in.Password)
		if err != nil {
			return &res, err
		}

		user, err := dbops.TryLogin(login, pwd)
		if err != nil {
			return &res, err
		}
		user.Password = ""

		jm, err := json.Marshal(user)
		if err != nil {
			return &res, err
		}

		res.Result = result("true", string(jm))
		return &res, nil
	}

	if op == "get-profile" {
		var u dbops.User
		err := json.Unmarshal([]byte(instructions), &u)
		if err != nil {
			return &res, err
		}

		user, err := dbops.GetProfile(u)
		if err != nil {
			return &res, err
		}
		user.Password = ""

		jm, err := json.Marshal(user)
		if err != nil {
			return &res, err
		}

		res.Result = result("true", string(jm))
		return &res, nil
	}

	if op == "hash" {
		var in dbops.User
		err := json.Unmarshal([]byte(instructions), &in)
		if err != nil {
			return &res, err
		}

		res.Result = hashing.GeneratePassword(in.Password)
		return &res, nil
	}

	if op == "validate" {
		var in dbops.User
		err := json.Unmarshal([]byte(instructions), &in)
		if err != nil {
			return &res, err
		}

		res.Result = "true"
		err = hashing.ValidatePassword([]byte(in.Login), []byte(in.Password))
		if err != nil {
			return &res, err
		}

		return &res, nil
	}

	if op == "refresh" {
		var in dbops.User
		err := json.Unmarshal([]byte(instructions), &in)
		if err != nil {
			return &res, err
		}

		res.Result = "false"
		//in.Login user string id, in.Password cookie jti hash
		user, err := dbops.TryRefresh(in.Login, in.Password)
		if err != nil {
			return &res, err
		}

		jmuser, err := json.Marshal(user)
		if err != nil {
			return &res, err
		}

		res.Result = result("true", string(jmuser))
		return &res, nil
	}

	if op == "updateRef" {
		var in dbops.User
		err := json.Unmarshal([]byte(instructions), &in)
		if err != nil {
			return &res, err
		}

		//in.Login user string id, in.Password cookie jti hash
		err = dbops.UpdateRefresh(in.Login, in.Password)
		if err != nil {
			return &res, err
		}

		res.Result = result("true", "updateRef")
		return &res, nil
	}

	if op == "update-login" {
		var user dbops.User
		err := json.Unmarshal([]byte(instructions), &user)
		if err != nil {
			return &res, err
		}

		err = dbops.UpdateLogin(user)
		if err != nil {
			return &res, err
		}

		res.Result = result("true", `"updated successfully"`)
		return &res, nil
	}

	//countries
	if op == "read-countries" {
		countries, err := dbops.ReadCountries()
		if err != nil {
			return &res, err
		}

		jm, err := json.Marshal(countries)
		if err != nil {
			return &res, err
		}

		res.Result = result("true", string(jm))
		return &res, nil
	}

	if op == "read-regions" {
		var r dbops.Region
		err := json.Unmarshal([]byte(instructions), &r)
		if err != nil {
			return &res, err
		}

		regions, err := dbops.ReadRegions(r.CountryId)
		if err != nil {
			return &res, err
		}

		jm, err := json.Marshal(regions)
		if err != nil {
			return &res, err
		}

		res.Result = result("true", string(jm))
		return &res, nil
	}

	if op == "read-towns" {
		var t dbops.Town
		err := json.Unmarshal([]byte(instructions), &t)
		if err != nil {
			return &res, err
		}

		towns, err := dbops.ReadTowns(t.RegionId)
		if err != nil {
			return &res, err
		}

		jm, err := json.Marshal(towns)
		if err != nil {
			return &res, err
		}

		res.Result = result("true", string(jm))
		return &res, nil
	}

	if op == "new-country" {
		var country dbops.Country
		err := json.Unmarshal([]byte(instructions), &country)
		if err != nil {
			return &res, err
		}

		country, err = dbops.NewCountry(country)
		if err != nil {
			return &res, err
		}

		jm, err := json.Marshal(country)
		if err != nil {
			return &res, err
		}

		res.Result = result("true", string(jm))
		return &res, nil
	}

	if op == "get-profile-comments" {
		var u dbops.User
		err := json.Unmarshal([]byte(instructions), &u)
		if err != nil {
			return &res, err
		}

		comments, err := dbops.GetProfileComments(u.Id)
		if err != nil {
			return &res, err
		}

		jm, err := json.Marshal(comments)
		if err != nil {
			return &res, err
		}

		res.Result = result("true", string(jm))
		return &res, nil
	}

	//when a master updates their skills choices
	if op == "update-service-choices" {
		ids := struct {
			LoginId    int32   `json:"login_id"`
			ServiceIds []int32 `json:"service_ids"`
		}{}
		err := json.Unmarshal([]byte(instructions), &ids)
		if err != nil {
			return &res, err
		}

		var choices []dbops.Choice
		for _, v := range ids.ServiceIds {
			choices = append(choices, dbops.Choice{
				LoginId:   ids.LoginId,
				ServiceId: v,
			})
		}

		err = dbops.UpdateServiceChoices(choices)
		if err != nil {
			return &res, err
		}

		res.Result = result("true", `"update-service-choices"`)
		return &res, nil
	}

	if op == "update-service-prices" {
		var choices []dbops.Choice
		err := json.Unmarshal([]byte(instructions), &choices)
		if err != nil {
			return &res, err
		}

		err = dbops.UpdateServicePrices(choices)
		if err != nil {
			return &res, err
		}

		res.Result = result("true", `"update-service-prices"`)
		return &res, nil
	}

	if op == "masters-choices" {
		var u dbops.User
		err := json.Unmarshal([]byte(instructions), &u)
		if err != nil {
			return &res, err
		}

		choices, err := dbops.GetMastersChoices(u.Id)
		if err != nil {
			return &res, err
		}

		jm, err := json.Marshal(choices)
		if err != nil {
			return &res, err
		}

		res.Result = result("true", string(jm))
		return &res, nil
	}

	//portfolio stuff
	if op == "add-work" {
		err := dbops.AddWork(instructions)
		if err != nil {return &res, err}
		res.Result = result("true", `"added successfully"`)
		return &res, nil
	}

	if op == "update-work" {
		err := dbops.UpdateWork(instructions)
		if err != nil {return &res, err}
		res.Result = result("true", `"updated successfully"`)
		return &res, nil
	}

	if op == "get-portfolio" {
		str, err := dbops.GetPortfolio(instructions)
		if err != nil {return &res, err}
		res.Result = result("true", str)
		return &res, nil
	}

	//orders
	if op == "add-order" {
		str, err := dbops.AddOrder(instructions)
		if err != nil {return &res, err}
		res.Result = result("true", str)
		return &res, nil
	}

	if op == "get-orders" {
		str, err := dbops.GetOrders(instructions)
		if err != nil {return &res, err}
		res.Result = result("true", str)
		return &res, nil
	}

	return &res, nil
}

func main() {
	//init logging
	f, err := os.OpenFile(os.Getenv("GOSERVICES_LOG"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	ok, err := credentials.NewServerTLSFromFile(os.Getenv("SERVICEKEY_PEM"), os.Getenv("SERVICEKEY_KEY"))
	if err != nil {
		log.Fatalf("Failed to setup TLS:%v", err)
	}

	lis, err := net.Listen("tcp", ":50003")
	if err != nil {
		log.Fatal(service + "service failed to listen ", err)
	}

	println("Hi, I'm an " + service + " grpc comm. service listening...")

	s := grpc.NewServer(grpc.Creds(ok))
	grpcc.RegisterCommunicationServiceServer(s, &server{})
	err = s.Serve(lis)
	if err != nil {
		log.Fatal("Failed to serve grpc server " + service + ":", err)
	}
}
