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

  <div class="col-4" v-else></div>

  <div class="col-auto" v-if="hasToggle()">
    <device-state :deviceid="props.deviceid"></device-state>
  </div>
  <div class="col-auto" v-else></div>
</template>

<script setup>
import { DeviceId, DeviceName, DeviceDimmer, DeviceState } from ".";
import store from "../../store";

const props = defineProps({
  deviceid: {
    required: true,
  },
});

const isDimmable = () => {
  return store.state.websocket.devices[props.deviceid].type == 2;
};

const hasToggle = () => {
  let type = store.state.websocket.devices[props.deviceid].type;
  return type == 2 || type == 3;
};
</script>
