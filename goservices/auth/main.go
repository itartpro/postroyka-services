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

type server struct{}

type Instructions struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var service = "auth"

func result(status string, data string) string {
	return `{"name":"` + service + `","status":` + status + `,"data":` + data + `}`
}

func (*server) PassData(ctx context.Context, req *grpcc.DataRequest) (*grpcc.DataResponse, error) {

	var res grpcc.DataResponse
	res.Result = result("false", `"noop or error"`)

	instructions := req.GetData().GetInstructions()
	op := req.GetData().GetAction()

	if op == "register" {
		var user dbops.User
		err := json.Unmarshal([]byte(instructions), &user)
		if err != nil {
			return &res, err
		}

		user.Password = hashing.GeneratePassword(user.Password)
		user, err = dbops.TryRegister(user)
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

	if op == "login" {
		var in Instructions
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

	if op == "hash" {
		var in Instructions
		err := json.Unmarshal([]byte(instructions), &in)
		if err != nil {
			return &res, err
		}

		res.Result = hashing.GeneratePassword(in.Password)
		return &res, nil
	}

	if op == "validate" {
		var in Instructions
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
		var in Instructions
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
		var in Instructions
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
	if op == "update_service_choices" {
		ids := struct {
			LoginId    int32   `json:"login_id"`
			ServiceIds []int32 `json:"service_ids"`
		}{}
		err := json.Unmarshal([]byte(instructions), &ids)
		if err != nil {
			return &res, err
		}

		var choices []dbops.ServiceChoice
		for _, v := range ids.ServiceIds {
			choices = append(choices, dbops.ServiceChoice{
				LoginId:   ids.LoginId,
				ServiceId: v,
			})
		}

		err = dbops.UpdateServiceChoices(choices)
		if err != nil {
			return &res, err
		}

		res.Result = result("true", "update_service_choices")
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
	s := grpc.NewServer(grpc.Creds(ok))
	lis, err := net.Listen("tcp", ":50003")
	if err != nil {
		log.Fatal("Failed to listen ", err)
	}

	println("Hi, I'm a " + service + " microservice listening...")

	grpcc.RegisterCommunicationServiceServer(s, &server{})
	err = s.Serve(lis)
	if err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
