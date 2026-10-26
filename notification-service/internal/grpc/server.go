package grpc

import (
	"context"
	"fmt"

	"notification-service/internal/controllers"

	pb "semen_project/proto/gen/notificationpb"
)

type Server struct {
	pb.UnimplementedNotificationServiceServer
	handler *controllers.NotificationHandler
}

func NewServer(handler *controllers.NotificationHandler) *Server {
	return &Server{
		handler: handler,
	}
}

func (s *Server) GetNotifications(
	ctx context.Context,
	req *pb.GetNotificationsRequest,
) (*pb.GetNotificationsResponse, error) {
	notifications, err := s.handler.GetNotifications(
		ctx,
		req.GetUserId(),
		req.GetLimit(),
		req.GetOffset(),
	)
	if err != nil {
		return nil, fmt.Errorf("get notifications: %w", err)
	}

	response := &pb.GetNotificationsResponse{
		Notifications: make([]*pb.Notification, 0, len(notifications)),
	}

	for _, n := range notifications {
		response.Notifications = append(
			response.Notifications,
			&pb.Notification{
				Id:        n.ID,
				UserId:    n.UserID,
				SenderId:  n.SenderID,
				MessageId: n.MessageID,
				Type:      n.Type,
				IsRead:    n.IsRead,
				CreatedAt: n.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			},
		)
	}

	return response, nil
}
func (s *Server) GetUnreadCount(
	ctx context.Context,
	req *pb.GetUnreadCountRequest,
) (*pb.GetUnreadCountResponse, error) {
	count, err := s.handler.GetUnreadCount(ctx, req.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("get unread count: %w", err)
	}

	return &pb.GetUnreadCountResponse{
		Count: count,
	}, nil
}

func (s *Server) MarkAsRead(
	ctx context.Context,
	req *pb.MarkAsReadRequest,
) (*pb.MarkAsReadResponse, error) {
	success, err := s.handler.MarkAsRead(
		ctx,
		req.GetUserId(),
		req.GetNotificationId(),
	)
	if err != nil {
		return nil, fmt.Errorf("mark as read: %w", err)
	}

	return &pb.MarkAsReadResponse{
		Success: success,
	}, nil
}
