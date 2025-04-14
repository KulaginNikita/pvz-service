package pvzgrpc

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/KulaginNikita/pvz-service/pkg/pvz_v1/pvz"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
	pb.UnimplementedPVZServiceServer
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{pool: pool}
}

func (h *Handler) GetPVZList(ctx context.Context, req *pb.GetPVZListRequest) (*pb.GetPVZListResponse, error) {
	rows, err := h.pool.Query(ctx, `SELECT id, registration_date, city FROM pvz`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*pb.PVZ
	for rows.Next() {
		var id string
		var regDate time.Time
		var city string

		if err := rows.Scan(&id, &regDate, &city); err != nil {
			return nil, err
		}

		result = append(result, &pb.PVZ{
			Id:               id,
			RegistrationDate: timestamppb.New(regDate),
			City:             city,
		})
	}

	return &pb.GetPVZListResponse{Pvzs: result}, nil
}
