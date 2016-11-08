package client

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	tiller "k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/version"
)

type TillerClient struct {
	address string
	context context.Context
}

func NewTillerClient(address string) *TillerClient {
	md := metadata.Pairs("x-helm-api-client", version.Version)
	ctx := metadata.NewContext(context.TODO(), md)
	return &TillerClient{address: address, context: ctx}
}

func (tc *TillerClient) execute(request func(tiller.ReleaseServiceClient)) error {
	conn, err := grpc.Dial(tc.address, grpc.WithInsecure())
	if err != nil {
		log.Debug("unable to dial tiller")
		return err
	}
	defer conn.Close()
	rsc := tiller.NewReleaseServiceClient(conn)
	request(rsc)
	return nil
}

func (tc *TillerClient) ListReleases(req *tiller.ListReleasesRequest) (res *tiller.ListReleasesResponse, err error) {
	log.Info(req)
	tc.execute(func(rsc tiller.ReleaseServiceClient) {
		lrc, err := rsc.ListReleases(tc.context, req)
		if err != nil {
			log.Debug("unable to list all releases")
			return
		}
		res, err = lrc.Recv()
	})
	return
}

func (tc *TillerClient) InstallRelease(req *tiller.InstallReleaseRequest) (res *tiller.InstallReleaseResponse, err error) {
	tc.execute(func(rsc tiller.ReleaseServiceClient) {
		res, err = rsc.InstallRelease(tc.context, req)
		if err != nil {
			log.Debug("unable to install release")
		}
	})
	return
}

func (tc *TillerClient) UninstallRelease(req *tiller.UninstallReleaseRequest) (res *tiller.UninstallReleaseResponse, err error) {
	tc.execute(func(rsc tiller.ReleaseServiceClient) {
		res, err = rsc.UninstallRelease(tc.context, req)
		if err != nil {
			log.Debug("unable to uninstall release")
		}
	})
	return
}

func (tc *TillerClient) GetReleaseContent(req *tiller.GetReleaseContentRequest) (res *tiller.GetReleaseContentResponse, err error) {
	tc.execute(func(rsc tiller.ReleaseServiceClient) {
		res, err = rsc.GetReleaseContent(tc.context, req)
		if err != nil {
			log.Debug("unable to get release content")
		}
	})
	return
}

func (tc *TillerClient) GetReleaseStatus(req *tiller.GetReleaseStatusRequest) (res *tiller.GetReleaseStatusResponse, err error) {
	tc.execute(func(rsc tiller.ReleaseServiceClient) {
		res, err = rsc.GetReleaseStatus(tc.context, req)
		if err != nil {
			log.Debug("unable to get release status")
		}
	})
	return
}
