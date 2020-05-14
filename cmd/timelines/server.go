package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"github.com/soheilhy/cmux"
	"github.com/zhughes3/kruonis/cmd/timelines/models"
	"google.golang.org/grpc"
)

type server struct {
	name       string
	lis        net.Listener
	grpcServer *grpc.Server
	httpServer *http.Server
	db         *db
	jwtKey     []byte
}

func NewServer(cfg *serverConfig, l net.Listener, d *sql.DB) *server {
	return &server{
		name: cfg.httpHost + ":" + cfg.httpPort,
		lis:  l,
		db: &db{
			db: d,
		},
		jwtKey: []byte(cfg.jwtKey),
	}
}

func (s *server) Start() error {
	var err error

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// tcpMuxer
	tcpMux := cmux.New(s.lis)

	// connection dispatcher rules
	grpcL := tcpMux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	httpL := tcpMux.Match(cmux.HTTP1Fast())

	//init gprc
	s.grpcServer, err = s.withGRPC(ctx)
	if err != nil {
		log.Fatal("Unable to init gRPC server")
		return err
	}

	s.httpServer, err = prepareHTTP(ctx, s.name)
	if err != nil {
		log.Fatalln("unable to init http server")
		return err
	}

	// start servers
	go func() {
		if err := s.grpcServer.Serve(grpcL); err != nil {
			log.Fatalln("unable to start grpc server")
		}
	}()

	go func() {
		if err := s.httpServer.Serve(httpL); err != nil {
			log.Fatalln("unable to start http server")
		}
	}()

	return tcpMux.Serve()
}

func (s *server) withGRPC(ctx context.Context) (*grpc.Server, error) {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(s.authUnaryServerInterceptor))
	models.RegisterTimelineServiceServer(grpcServer, s)
	return grpcServer, nil
}

func prepareHTTP(ctx context.Context, name string) (*http.Server, error) {
	router := http.NewServeMux()

	gw, err := prepareGateway(name, ctx)
	if err != nil {
		log.Fatalln("unable to init grpc gateway")
		return nil, err
	}
	router.Handle("/", gw)

	handler := cors.AllowAll().Handler(router)

	return &http.Server{
		Addr:         name,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}, nil
}

// setup the grpc-gateway
func prepareGateway(target string, ctx context.Context) (http.Handler, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
		return nil, err
	}

	gwMux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(httpResponseModifier),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}),
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),
	)

	err = models.RegisterTimelineServiceHandler(ctx, gwMux, conn)
	if err != nil {
		return nil, err
	}

	return gwMux, nil
}

func httpResponseModifier(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	// set cookies
	if tkn := md.HeaderMD.Get("access"); len(tkn) > 0 {
		cookie := http.Cookie{
			Name:     "access",
			Value:    tkn[0],
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		delete(md.HeaderMD, "access")
		delete(w.Header(), "Grpc-Metadata-Access")
	}
	if tkn := md.HeaderMD.Get("refresh"); len(tkn) > 0 {
		cookie := http.Cookie{
			Name:     "refresh",
			Value:    tkn[0],
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		delete(md.HeaderMD, "refresh")
		delete(w.Header(), "Grpc-Metadata-Refresh")
	}

	// set http status code
	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		w.WriteHeader(code)
		delete(md.HeaderMD, "x-http-code")
		delete(w.Header(), "Grpc-Metadata-X-Http-Code")
	}

	delete(w.Header(), "Grpc-Metadata-Content-Type")

	return nil
}
