package main

import (
	// 导如 protoc 自动生成的包
	pb "./proto/consignment"
	"context"
	"log"
	"github.com/micro/go-micro"
	vesselPb "../vessel-service/proto/vessel"
)

const (
	PORT = ":50051"
)

//
// 仓库接口
//
type IRepository interface {
	Create(consignment *pb.Consignment) (*pb.Consignment, error) // 存放新货物
}

//
// 我们存放多批货物的仓库，实现了 IRepository 接口
//
type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.consignments = append(repo.consignments, consignment)
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

//
// 定义微服务
//
type service struct {
	repo Repository
	// consignment-service 作为客户端调用 vessel-service 的函数
	vesselClient vesselPb.VesselServiceClient
}

//
// service 实现 consignment.pb.go 中的 ShippingServiceServer 接口
// 使 service 作为 gRPC 的服务端
//
// 托运新的货物
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) ( error) {
	// 检查是否有合适的货轮
	vReq := &vesselPb.Specification{
		Capacity: int32(len(req.Containers)),
		MaxWeight: req.Weight,
	}

	vResp , err := s.vesselClient.FindAvailable(context.Background(), vReq)
	if err != nil {
		return  err
	}

	// 货物被承运
	log.Printf("found vessel :%s\n", vResp.Vessel.Name)
	req.VesselId = vResp.Vessel.Id

	// 接收承运的货物
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	resp.Created = true
	resp.Consignment = consignment
	//resp = &pb.Response{Created: true, Consignment: consignment}
	return nil
}

func (s *service) GetConsignments(crx context.Context ,req *pb.GetRequest,resp *pb.Response) (error){
	allConsignments := s.repo.GetAll()
	resp.Consignments = allConsignments
	//resp = &pb.Response{
	//	Consignments: allConsignments,
	//}
	return nil
}

func main() {
	server := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	//解析命令
	server.Init()
	repo := Repository{}

	// 作为 vessel-service 的客户端
	vClient := vesselPb.NewVesselServiceClient("go.micro.srv.vessel", server.Client())
	// 向 rRPC 服务器注册微服务
	// 此时会把我们自己实现的微服务 service 与协议中的 ShippingServiceServer 绑定
	pb.RegisterShippingServiceHandler(server.Server(), &service{repo, vClient})

	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}