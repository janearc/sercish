package main

import (
	"github.com/janearc/sux/sux"
	"github.com/sirupsen/logrus"
)

// simple state storage, but not S3 that's something different

func (s *Service) Boot() {
	// 1. ask sux what our blob looks like
	// 2. if we have a blob, load it
	// 3. if we don't have a blob, create one
	// 4. store that blob in the service struct

	// function getState (token string, secret string) *sux.Sux {
	// function setState(s sux.Sux, state string) error {

	// the prototype for OneBlob is:
	// func (s sux.Sux) OneBlob(token string, secret string) (json.RawMessage, error) {

	sss := s.sux

	// we create a new session here, but we're not going to want to do this unless
	// we say hey i got a sid here do you have a blob in that big ol database of yours
	// and sux is like nah i ain't know what that is
	sid, err := sux.NewSession()

	if err != nil {
		logrus.Fatalf("Failed to create session: %v", err)
	}

	// let's just make some shit up here
	oe := sss.Open(sid)

	if oe == nil {
		logrus.WithError(oe).Errorf("unknown state session [%s]", sid)
	}

	// let's check our work and make sure that sss knows we're live
	if sss.Live() != true {
		// this means that sux has to record the fact that we have state
		// 1. there's gotta be a json blob that we *can* dump in the backend
		// 2. there's gotta be a valid connection to the storage backend whatever that is
		// 3. we might care whether there are uncommitted changes to the blob
		//    but then we'd need a sum or something and that sounds messy.
		logrus.Errorf("state session [%s] opened successfully, but sss reports dead")
	}
}
