<template>
  <div ref="refToolbar">
    <q-toolbar>
      <q-space />
      <q-btn label="Update" color="primary" @click="UpdateDevices"></q-btn>
    </q-toolbar>
  </div>

  <scroller
    ref="scrollArea"
    :expand="true"
    :bottom-margin="20"
    style="width: 99%; border: 1px solid black"
  >
    <div v-for="(item, index) in devices" :key="index" class="row q-gutter-md">
      <div class="col-auto">
        <device-id :deviceid="item.id"></device-id>
      </div>
      <div class="col-3">
        <device-name :deviceid="item.id"></device-name>
      </div>
      <div class="col-auto">
        <device-state :deviceid="item.id"></device-state>
      </div>
    </div>
  </scroller>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { mapFields } from "vuex-map-fields";
import { useStore } from "vuex";
import scroller from "./Scroller.vue";

import { DeviceId, DeviceName, DeviceState } from "./device/index";

const store = useStore();

onMounted(() => {
  UpdateDevices();
});

const test = () => {};

const UpdateDevices = () => {
  console.log("UpdateDevices");
  store.dispatch(
    "sendWSCommand",
    JSON.stringify({
      class: "devices",
      command: "update",
    })
  );
};
</script>

<script>
export default {
  computed: {
    ...mapFields(["websocket.devices"]),
  },
  mounted() {},
};
</script>

<style>
.Scroll {
  background-color: bisque;
  color: v-bind("scrollColor");
  height: v-bind("scrollHeight");
  border: 1px solid black;
}
</style>
