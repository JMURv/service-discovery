package mapper

import (
	pb "github.com/JMURv/service-discovery/api/pb"
	md "github.com/JMURv/service-discovery/pkg/model"
)

func ListSvcToProto(u []md.Service) []*pb.ServiceMsg {
	res := make([]*pb.ServiceMsg, len(u))
	for i := 0; i < len(u); i++ {
		res[i] = SvcToProto(&u[i])
	}
	return res
}

func SvcToProto(req *md.Service) *pb.ServiceMsg {
	return &pb.ServiceMsg{
		Name:     req.Name,
		Address:  req.Address,
		IsActive: req.IsActive,
	}
}
