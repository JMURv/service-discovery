package mapper

import (
	pb "github.com/JMURv/service-discovery/api/pb"
	md "github.com/JMURv/service-discovery/pkg/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ListItemToProto(u []*md.Item) []*pb.ItemMsg {
	res := make([]*pb.ItemMsg, len(u))
	for i, v := range u {
		res[i] = ItemToProto(v)
	}
	return res
}

func ItemToProto(req *md.Item) *pb.ItemMsg {
	item := &pb.ItemMsg{
		Id:              req.ID.String(),
		Article:         req.Article,
		Title:           req.Title,
		Description:     req.Description,
		Price:           float32(req.Price),
		QuantityInStock: uint32(req.QuantityInStock),
		Src:             req.Src,
		Alt:             req.Alt,
		InStock:         req.InStock,
		IsHit:           req.IsHit,
		IsRec:           req.IsRec,
		ParentItemId:    req.ParentItemID.String(),
		Media:           ListItemMediaToProto(req.Media),
		Attributes:      ListItemAttributesToProto(req.Attributes),
		CreatedAt: &timestamppb.Timestamp{
			Seconds: req.CreatedAt.Unix(),
			Nanos:   int32(req.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: req.UpdatedAt.Unix(),
			Nanos:   int32(req.UpdatedAt.Nanosecond()),
		},
	}

	if len(req.RelatedProducts) > 0 {
		res := make([]*pb.RelatedProduct, len(req.RelatedProducts))
		for i, v := range req.RelatedProducts {
			res[i] = RelatedProductsToProto(&v)
		}
		item.RelatedProducts = res
	}

	if len(req.Categories) > 0 {
		res := make([]*pb.CategoryMsg, len(req.Categories))
		for i, v := range req.Categories {
			res[i] = CategoryToProto(v)
		}
		item.Categories = res
	}

	if len(req.Variants) > 0 {
		res := make([]*pb.ItemMsg, len(req.Variants))
		for i, v := range req.Variants {
			res[i] = ItemToProto(&v)
		}
		item.Variants = res
	}

	return item
}

func ListRelatedProductsToProto(u []*md.RelatedProduct) []*pb.RelatedProduct {
	res := make([]*pb.RelatedProduct, len(u))
	for i, v := range u {
		res[i] = RelatedProductsToProto(v)
	}
	return res
}

func RelatedProductsToProto(req *md.RelatedProduct) *pb.RelatedProduct {
	return &pb.RelatedProduct{
		Id:            req.ID,
		ItemId:        req.ItemID.String(),
		RelatedItemId: req.RelatedItemID.String(),
		RelatedItem:   ItemToProto(&req.RelatedItem),
		CreatedAt: &timestamppb.Timestamp{
			Seconds: req.CreatedAt.Unix(),
			Nanos:   int32(req.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: req.UpdatedAt.Unix(),
			Nanos:   int32(req.UpdatedAt.Nanosecond()),
		},
	}
}

func ListItemAttributesToProto(req []md.ItemAttribute) []*pb.ItemAttribute {
	res := make([]*pb.ItemAttribute, len(req))
	for i, v := range req {
		res[i] = &pb.ItemAttribute{
			Id:     v.ID,
			Name:   v.Name,
			Value:  v.Value,
			ItemId: v.ItemID.String(),
			CreatedAt: &timestamppb.Timestamp{
				Seconds: v.CreatedAt.Unix(),
				Nanos:   int32(v.CreatedAt.Nanosecond()),
			},
			UpdatedAt: &timestamppb.Timestamp{
				Seconds: v.UpdatedAt.Unix(),
				Nanos:   int32(v.UpdatedAt.Nanosecond()),
			},
		}
	}
	return res
}

func ListItemMediaToProto(req []md.ItemMedia) []*pb.ItemMedia {
	res := make([]*pb.ItemMedia, len(req))
	for i, v := range req {
		res[i] = &pb.ItemMedia{
			Id:     v.ID,
			Src:    v.Src,
			Alt:    v.Alt,
			ItemId: v.ItemID.String(),
			CreatedAt: &timestamppb.Timestamp{
				Seconds: v.CreatedAt.Unix(),
				Nanos:   int32(v.CreatedAt.Nanosecond()),
			},
			UpdatedAt: &timestamppb.Timestamp{
				Seconds: v.UpdatedAt.Unix(),
				Nanos:   int32(v.UpdatedAt.Nanosecond()),
			},
		}
	}
	return res
}

func ItemFromProto(req *pb.ItemMsg) *md.Item {
	modelItem := &md.Item{
		Article:         req.Article,
		Title:           req.Title,
		Description:     req.Description,
		Price:           float64(req.Price),
		QuantityInStock: int(req.QuantityInStock),
		Src:             req.Src,
		Alt:             req.Alt,
		InStock:         req.InStock,
		IsHit:           req.IsHit,
		IsRec:           req.IsRec,
		CreatedAt:       req.CreatedAt.AsTime(),
		UpdatedAt:       req.UpdatedAt.AsTime(),
	}

	uid, err := uuid.Parse(req.Id)
	if err != nil {
		zap.L().Debug("failed to parse user id")
	} else {
		modelItem.ID = uid
	}

	parentItemID, err := uuid.Parse(req.ParentItemId)
	if err != nil {
		zap.L().Debug("failed to parse parent item id")
	} else {
		modelItem.ParentItemID = &parentItemID
	}

	if len(req.Categories) > 0 {
		res := make([]*md.Category, len(req.Categories))
		for i, v := range req.Categories {
			res[i] = CategoryFromProto(v)
		}
		modelItem.Categories = res
	}

	if len(req.Media) > 0 {
		res := make([]md.ItemMedia, len(req.Media))
		for i, v := range req.Media {
			iMedia := md.ItemMedia{
				ID:        v.Id,
				Src:       v.Src,
				Alt:       v.Alt,
				CreatedAt: v.CreatedAt.AsTime(),
				UpdatedAt: v.UpdatedAt.AsTime(),
			}

			uid, err := uuid.Parse(v.ItemId)
			if err != nil {
				zap.L().Debug("failed to parse user id")
			} else {
				iMedia.ItemID = uid
			}

			res[i] = iMedia
		}
		modelItem.Media = res
	}

	if len(req.Attributes) > 0 {
		res := make([]md.ItemAttribute, len(req.Attributes))
		for i, v := range req.Attributes {
			res[i] = md.ItemAttribute{
				ID:        v.Id,
				Name:      v.Name,
				Value:     v.Value,
				CreatedAt: v.CreatedAt.AsTime(),
				UpdatedAt: v.UpdatedAt.AsTime(),
			}
		}
		modelItem.Attributes = res
	}

	if len(req.RelatedProducts) > 0 {
		res := make([]md.RelatedProduct, len(req.RelatedProducts))
		for i, v := range req.RelatedProducts {
			pr := md.RelatedProduct{
				ID:          v.Id,
				RelatedItem: *ItemFromProto(v.RelatedItem),
				CreatedAt:   v.CreatedAt.AsTime(),
				UpdatedAt:   v.UpdatedAt.AsTime(),
			}

			itemUid, err := uuid.Parse(v.ItemId)
			if err != nil {
				zap.L().Debug("failed to parse user id")
			} else {
				pr.ItemID = itemUid
			}

			relatedItemUid, err := uuid.Parse(v.RelatedItemId)
			if err != nil {
				zap.L().Debug("failed to parse user id")
			} else {
				pr.RelatedItemID = relatedItemUid
			}

			res[i] = pr
		}
		modelItem.RelatedProducts = res
	}

	return modelItem
}
