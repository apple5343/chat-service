package project

import (
	"context"
	pb "project-service/pkg/api/project_v1"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) ListUserProjects(ctx context.Context, req *emptypb.Empty) (*pb.ListUserProjectsResponse, error) {
	var userID string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		userID = md["user_id"][0]
	}
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Authentication required. Please provide a valid token.")
	}

	projects, err := s.projectService.GetAllByUserID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("error while listing user projects")
		return nil, status.Error(codes.Internal, "Failed get user projects.")
	}

	response := make([]*pb.Project, len(projects))
	for i, project := range projects {
		response[i] = &pb.Project{
			Id:                           project.ID,
			Title:                        project.Title,
			Description:                  project.Description,
			Users:                        project.Users,
			AdminId:                      project.AdminID,
			NotificationSubscribersTgIds: project.NotificationSubscribersTGIDS,
			CreatedAt:                    timestamppb.New(project.CreatedAt),
		}
	}

	return &pb.ListUserProjectsResponse{
		Project: response,
	}, nil
}
