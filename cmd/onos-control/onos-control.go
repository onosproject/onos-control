// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

/*
Package onos-control is the main entry point to the ONOS control subsystem.

It connects to devices through a Southbound P4Runtime interface and
in turn provides a P4Runtime interface northbound for other systems to connect to.

Arguments

-caPath <the location of a CA certificate>

-keyPath <the location of a client private key>

-certPath <the location of a client certificate>


See ../../docs/run.md for how to run the application.
*/
package main

import (
	"flag"
	"github.com/onosproject/onos-control/pkg/manager"
	"github.com/onosproject/onos-control/pkg/northbound"
	"github.com/onosproject/onos-control/pkg/northbound/admin"
	"github.com/onosproject/onos-control/pkg/northbound/diags"
	log "k8s.io/klog"
)

// The main entry point
func main() {
	caPath := flag.String("caPath", "", "path to CA certificate")
	keyPath := flag.String("keyPath", "", "path to client private key")
	certPath := flag.String("certPath", "", "path to client certificate")

	//lines 93-109 are implemented according to
	// https://github.com/kubernetes/klog/blob/master/examples/coexist_glog/coexist_glog.go
	// because of libraries importing glog. With glog import we can't call log.InitFlags(nil) as per klog readme
	// thus the alsologtostderr is not set properly and we issue multiple logs.
	// Calling log.InitFlags(nil) throws panic with error `flag redefined: log_dir`
	err := flag.Set("alsologtostderr", "true")
	if err != nil {
		log.Error("Cant' avoid double Error logging ", err)
	}
	flag.Parse()

	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	log.InitFlags(klogFlags)

	// Sync the glog and klog flags.
	flag.CommandLine.VisitAll(func(f1 *flag.Flag) {
		f2 := klogFlags.Lookup(f1.Name)
		if f2 != nil {
			value := f1.Value.String()
			_ = f2.Value.Set(value)
		}
	})
	log.Info("Starting onos-control")

	mgr, err := manager.LoadManager()
	if err != nil {
		log.Fatal("Unable to load onos-control ", err)
	} else {
		mgr.Run()
		err = startServer(*caPath, *keyPath, *certPath)
		if err != nil {
			log.Fatal("Unable to start onos-control ", err)
		}
	}
}

// Creates gRPC server and registers various services; then serves.
func startServer(caPath string, keyPath string, certPath string) error {
	s := northbound.NewServer(northbound.NewServerConfig(caPath, keyPath, certPath))
	s.AddService(admin.Service{})
	s.AddService(diags.Service{})

	return s.Serve(func(started string) {
		log.Info("Started NBI on ", started)
	})
}
