package main

// TODO: защита от частого включения-выключения должна быть реализована непостредственно
// в схеме управления электродвигателем

import (
	"strconv"
	"context"
	"errors"
	"time"
	
	ventService_RPC "ventService/proto"
)

// VentServiceImpl structure used to store channel for stop signal
type VentServiceImpl struct {
	State            bool
	LastTurnOnTime   time.Time
	LastWorkDuration time.Duration
}

var ventServiceState = VentServiceImpl{State: false}

// TurnOff will turn off channel vent regardless of sensor values
func (state *VentServiceImpl) TurnOff(context.Context, *ventService_RPC.TurnOffMsg, *ventService_RPC.TurnOffResponse) error {
	return state.off()
}

func (state *VentServiceImpl) off() error {
	if state.State != false {
		state.State = false
		state.LastWorkDuration = time.Since(state.LastTurnOnTime)
	}

	if !mqttClient.IsConnectionOpen() {
		return errors.New("MQTT connection lost")
	}

	//fmt.Println("Turn vent OFF")
	if token := mqttClient.Publish("ROOT/Actuators/BathVentActuator/State", 0, false, strconv.FormatBool(state.State)); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	if token := mqttClient.Publish("ROOT/Actuators/BathVentActuator/LastWorkDuration", 0, false, state.LastWorkDuration.String()); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// TurnOn will turn on channel vent regardless of sensor values
func (state *VentServiceImpl) TurnOn(context.Context, *ventService_RPC.TurnOnMsg, *ventService_RPC.TurnOnResponse) error {
	return state.on()
}

func (state *VentServiceImpl) on() error {
	if state.State != true {
		state.State = true
		state.LastTurnOnTime = time.Now()
	}

	if !mqttClient.IsConnectionOpen() {
		return errors.New("MQTT connection lost")
	}

	//fmt.Println("Turn vent ON")
	if token := mqttClient.Publish("ROOT/Actuators/BathVentActuator/State", 0, false, strconv.FormatBool(state.State)); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	if token := mqttClient.Publish("ROOT/Actuators/BathVentActuator/LastTurnOnTime", 0, false, state.LastTurnOnTime.String()); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

