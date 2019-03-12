package main

import(
	pb "./proto/vessel"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"github.com/micro/go-micro"
	"log"
)
type Repository interface {
	FindAvailable(*pb.Specification)(* pb.Vessel , error)
}

type VesselRepository struct{
	vessles []*pb.Vessel
}

//实现接口
func (repo *VesselRepository)FindAvailable(spec *pb.Specification)(*pb.Vessel,error){
	//选择最近一条容量、载重都符合的货轮
	for _, v := range repo.vessles{
		if v.Capacity >= spec.Capacity && v.MaxWeight >= spec.MaxWeight{
			return v , nil
		}
	}
	return nil, errors.New("no vessele can be use")
}

// 定义活产服务
type service struct {
	repo Repository
}

func(s *service) FindAvailable(ctx context.Context , spec *pb.Specification, resp * pb.Response)error{
	//调用内部方法查找
	v, err := s.repo.FindAvailable(spec)
	if err != nil {
		return err
	}
	resp.Vessel =v
	return nil
}
func main() {
	//停留在港口的货船， 先写死
	vessles := []*pb.Vessel{
		{Id: "vessel001", Name:"Boaty McBoatface" , MaxWeight: 200000, Capacity: 500},
	}
	repo := &VesselRepository{vessles}
	server := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)
	server.Init()

	// 将实现服务端的API 注册到服务端
	pb.RegisterVesselServiceHandler(server.Server() , &service{repo})

	if err := server.Run() ; err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
