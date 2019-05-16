package server

import (
	"context"
	"github.com/bootapp/srv-ui/proto/clients/dal-ui"
	"github.com/bootapp/rest-grpc-oauth2/auth"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type SrvInfoServiceServer struct {
	dalInfoClient dal_ui.InfoServiceClient
	dalInfoConn   *grpc.ClientConn
	auth          *auth.StatelessAuthenticator
}

func NewSrvInfoServiceServer(dalUiAddr string) *SrvInfoServiceServer {
	var err error
	s := &SrvInfoServiceServer{}
	s.dalInfoConn, err = grpc.Dial(dalUiAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect : %v", err)
	}

	s.dalInfoClient = dal_ui.NewInfoServiceClient(s.dalInfoConn)
	s.auth = auth.GetInstance()
	return s
}

func (s *SrvInfoServiceServer) close() {
	err := s.dalInfoConn.Close()
	if err != nil {
		glog.Error(err)
	}
}

func (s *SrvInfoServiceServer) SaveDicts(ctx context.Context, req *dal_ui.DictList) (*dal_ui.DictList, error) {
	glog.Info("GRpc Request Save Dicts")
	resp, err := s.dalInfoClient.SaveDict(ctx, req)
	if err != nil {
		glog.Error(err)
		return nil, status.Error(codes.Unavailable, err.Error())

	}
	return resp, nil
}

func (s *SrvInfoServiceServer) SaveArticles(ctx context.Context, req *dal_ui.ArticleList) (*dal_ui.ArticleList, error) {
	glog.Info("GRpc Request Save Articles")
	resp, err := s.dalInfoClient.SaveArticle(ctx, req)
	if err != nil {
		glog.Error(err)
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	return resp, nil
}

func (s *SrvInfoServiceServer) UpdateDictsStatusByDictName(ctx context.Context, req *dal_ui.PublishReq) (*dal_ui.DictList, error) {
	glog.Infof("GRpc Request Update Dicts Status By dictName : %s , status : %s .", req.DictName, req.Status)
	resp, err := s.dalInfoClient.UpdateDictStatusByDictName(ctx, req)
	if err != nil {
		glog.Error(err)
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	return resp, nil
}

func (s *SrvInfoServiceServer) BatchDeleteDictsById(ctx context.Context, req *dal_ui.BatchDictId) (*empty.Empty, error) {
	glog.Info("GRpc Request Batch Delete Dicts .")
	resp, err := s.dalInfoClient.BatchDeleteDictById(ctx, req)
	if err != nil {
		glog.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *SrvInfoServiceServer) BatchDeleteArticlesById(ctx context.Context, req *dal_ui.BatchArtId) (*empty.Empty, error) {
	glog.Info("GRpc Request Batch Delete Articles .")
	resp, err := s.dalInfoClient.BatchDeleteArticleById(ctx, req)
	if err != nil {
		glog.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *SrvInfoServiceServer) QueryDictsPage(ctx context.Context, req *dal_ui.DictPageReq) (*dal_ui.DictPageResp, error) {
	glog.Infof("GRpc Request Query Dicts by page : %s .", req.Page)
	resp, err := s.dalInfoClient.QueryDictPage(ctx, req)
	if err != nil {
		glog.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *SrvInfoServiceServer) QueryArticle(ctx context.Context, req *dal_ui.ArticleReq) (*dal_ui.Article, error) {
	glog.Info("GRpc Request Query Article .", req)
	resp, err := s.dalInfoClient.QueryArticle(ctx, req)
	if err != nil {
		glog.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *SrvInfoServiceServer) QueryMultiDictsByParent(ctx context.Context, req *dal_ui.MultiDictReq) (*dal_ui.MultiDictResp, error) {
	glog.Infof("GRpc Request Query Multi Dicts by parent : %s .", req.Parent)
	resp, err := s.dalInfoClient.QueryMultiDictByParent(ctx, req)
	if err != nil {
		glog.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *SrvInfoServiceServer) QueryMultiArticlesByDictName(ctx context.Context, req *dal_ui.MultiArticleReq) (*dal_ui.MultiArticleResp, error) {
	glog.Infof("GRpc Request Query Multi Articles by dictName : %s .", req.DictName)
	resp, err := s.dalInfoClient.QueryMultiArticleByDictName(ctx, req)
	if err != nil {
		glog.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
