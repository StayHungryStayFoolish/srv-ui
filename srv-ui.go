package main

import (
	"context"
	"flag"
	"github.com/bootapp/srv-ui/server"
	"github.com/bootapp/rest-grpc-oauth2/auth"
	"github.com/bootapp/srv-core/oauth"
	pb "github.com/bootapp/srv-ui/proto/server"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	grpcEndpoint  = flag.String("grpc_endpoint", ":9093", "The endpoint of the ui grpc service")
	httpEndpoint  = flag.String("http_endpoint", ":8093", "The endpoint of the ui restful service")
	oauthEndpoint = flag.String("oauth_endpoint", ":9081", "The endpoint of the oauth server")
)

func main() {
	_ = flag.Set("alsologtostderr", "false")
	flag.Parse()
	defer glog.Flush()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(context.Background())
	authenticator := auth.GetInstance()
	oauthServer := oauth.NewPassOAuthServer()
	rpcSrv := server.GRpcServiceAddr{}

	server.ApolloConfig(ctx, false, &rpcSrv, &oauthServer, authenticator)

	go func() {
		glog.Info("oauth server listening ...")
		glog.Fatal(http.ListenAndServe(*oauthEndpoint, nil))
		glog.Info("oauth stopped")
	}()

	go func() {
		defer cancel()
		_ = gwRun(ctx, *httpEndpoint, *grpcEndpoint)
	}()
	_ = grpcRun(ctx, *grpcEndpoint, rpcSrv)
}

func grpcRun(ctx context.Context, grpcEndpoint string, addr server.GRpcServiceAddr) error {
	l, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	srvUiSrv := server.NewSrvInfoServiceServer(addr.DALUiSrv)
	pb.RegisterSrvInfoServiceServer(grpcServer, srvUiSrv)

	go func() {
		defer grpcServer.GracefulStop()
		<-ctx.Done()
		glog.Info("grpc server shutting down ...")
	}()
	glog.Info("grpc server running ...")
	return grpcServer.Serve(l)
}

func gwRun(ctx context.Context, httpEndpoint string, grpcEndpoint string) error {
	mux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(auth.GatewayResponseCookieAnnotator),
		runtime.WithMetadata(auth.GatewayRequestCookieParser))
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := pb.RegisterSrvInfoServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts); err != nil {
		glog.Fatal("Failed to start rest gateway : %v ", err)
		return err
	}

	srv := &http.Server{
		Addr:    httpEndpoint,
		Handler: mux,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		glog.Info("Rest gateway shutting down ...")
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			glog.Fatal("Failed to shutdown rest gateway %v", err)
		}
	}()
	glog.Info("RESTFull gateway running ...")
	return srv.ListenAndServe()
}
