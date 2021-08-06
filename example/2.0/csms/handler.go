package main

import (
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/availability"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/sirupsen/logrus"
)

type ConnectorState struct {
	status         availability.ConnectorStatus
	firmwareStatus firmware.FirmwareStatus
	connectors     map[int]*ConnectorInfo // No assumptions about the # of connectors
	transactions   map[int]*TransactionInfo
}

type TransactionInfo struct {
	id          int
	startTime   *types.DateTime
	endTime     *types.DateTime
	startMeter  int
	endMeter    int
	connectorId int
	idTag       string
}

type ConnectorInfo struct {
	status             core.ConnectorStatus
	currentTransaction int
}

type ChargingStationState struct {
	availabilityStatus availability.ChangeAvailabilityStatus
	//	diagnosticsStatus firmware.DiagnosticsStatus
	firmwareStatus firmware.FirmwareStatus
	connectors     map[int]*ConnectorInfo // No assumptions about the # of connectors
	transactions   map[int]*TransactionInfo
	//	errorCode         core.ChargePointErrorCode
}

type CSMSHandler struct {
	chargingStations map[string]*ChargingStationState
}

//------------- Core Profile Callbacks -----------

func (handler *CSMSHandler) OnHeartbeat(chargePointId string, request *availability.HeartbeatRequest) (confirmation *availability.HeartbeatResponse, err error) {
	logDefault(chargePointId, request.GetFeatureName()).Infof("heartbeat handled")
	return availability.NewHeartbeatResponse(types.NewDateTime(time.Now())), nil
}

func logDefault(chargePointId string, feature string) *logrus.Entry {
	return log.WithFields(logrus.Fields{"client": chargePointId, "message": feature})
}
