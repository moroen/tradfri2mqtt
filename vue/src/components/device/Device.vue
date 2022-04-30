<template>
  <div class="col-auto">
    <device-id :deviceid="props.deviceid"></device-id>
  </div>
  <div class="col-3">
    <device-name :deviceid="props.deviceid"></device-name>
  </div>
  <div class="col-4" v-if="isDimmable()">
    <device-dimmer :deviceid="props.deviceid"></device-dimmer>
  </div>

  <div class="col-4" v-else-if="isBlind()">
    <device-position :deviceid="props.deviceid"></device-position>
  </div>

  <div class="col-4" v-else></div>

  <div class="col-1" v-if="hasToggle()">
    <device-state :deviceid="props.deviceid"></device-state>
  </div>
  <div class="col-1" v-else></div>

  <div>
    <device-battery class="col-1 items-start text-grey-" :deviceid="props.deviceid"></device-battery>
  </div>

</template>

<script setup>
import { DeviceId, DeviceName, DeviceDimmer, DeviceState, DevicePosition, DeviceBattery } from ".";
import store from "../../store";

const props = defineProps({
  deviceid: {
    required: true,
  },
});

const isDimmable = () => {
  return store.state.websocket.devices[props.deviceid].type == 2;
};

const isBlind = () => {
  return store.state.websocket.devices[props.deviceid].type == 7;
};

const hasToggle = () => {
  let type = store.state.websocket.devices[props.deviceid].type;
  return type == 2 || type == 3;
};
</script>
