package grpc

import (
	notificationv1 "github.com/JrMarcco/jotice-api/api/notification/v1"
	nsvc "github.com/JrMarcco/jotice/internal/service/notification"
	tsvc "github.com/JrMarcco/jotice/internal/service/template"
)

type NotificationServer struct {
	notificationv1.UnimplementedNotificationServiceServer
	notificationv1.UnimplementedNotificationQueryServiceServer

	svc     nsvc.Service
	sendSvc nsvc.SendService
	txSvc   nsvc.TxService
	tplSvc  tsvc.TplService
}

func NewServer(
	notificationService nsvc.Service, sendService nsvc.SendService, txSvc nsvc.TxService,
	tplService tsvc.TplService,
) *NotificationServer {
	return &NotificationServer{
		svc:     notificationService,
		sendSvc: sendService,
		txSvc:   txSvc,
		tplSvc:  tplService,
	}
}
