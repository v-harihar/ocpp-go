package main

import (
	"os"
	"strconv"
	"time"

	ocpp2 "github.com/lorenzodonini/ocpp-go/ocpp2.0.1"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/sirupsen/logrus"
)

const (
	defaultListenPort          = 8887
	defaultHeartbeatInterval   = 600
	envVarServerPort           = "SERVER_LISTEN_PORT"
	envVarTls                  = "TLS_ENABLED"
	envVarCaCertificate        = "CA_CERTIFICATE_PATH"
	envVarServerCertificate    = "SERVER_CERTIFICATE_PATH"
	envVarServerCertificateKey = "SERVER_CERTIFICATE_KEY_PATH"
)

var log *logrus.Logger
var csms ocpp2.CSMS

func setupCSMS() ocpp2.CSMS {
	return ocpp2.NewCSMS(nil, nil)
}

func exampleRoutine(chargingStationID string, handler *CSMSHandler) {
	time.S

}

func main() {
	// Load config from ENV
	var listenPort = defaultListenPort
	port, _ := os.LookupEnv(envVarServerPort)
	if p, err := strconv.Atoi(port); err == nil {
		listenPort = p
	} else {
		log.Printf("no valid %v environment variable found, using default port", envVarServerPort)
	}

	csms = setupCSMS()

	handler := &CSMSHandler{chargingStations: map[string]*ChargingStationState{}}

	csms.SetAvailabilityHandler(handler)

	csms.SetNewChargingStationHandler(func(chargingStation ocpp2.ChargingStationConnection) {
		handler.chargingStations[chargingStation.ID()] = &ConnectorState{connectors: map[int]*ConnectorInfo{}, transactions: map[int]*TransactionInfo{}}
		log.WithField("client", chargingStation.ID()).Info("new charge point connected")
		go exampleRoutine(chargingStation.ID(), handler)
	})

	csms.SetChargingStationDisconnectedHandle(func(chargeStation ocpp2.ChargingStationConnection) {
		log.WithField("client", chargeStation.ID()).Info("charge point disconnected")
		delete(handler.chargingStations, chargeStation.ID())
	})

	ocppj.SetLogger(log)
	// Run csms
	log.Infof("starting csms on port %v", listenPort)
	csms.Start(listenPort, "/{ws}")
	log.Info("stopped csms")
}
