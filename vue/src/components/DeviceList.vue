<template>
  <div ref="refToolbar" class="q-ma-md">
    <q-toolbar>
      <q-select
        v-model="typeSelect"
        :options="options"
        emit-value
        map-options
        label="Device type"
        style="width: 300px"
      ></q-select>
      <q-space />

      <q-btn label="Update" color="primary" @click="UpdateDevices"></q-btn>
    </q-toolbar>
  </div>

  <scroller ref="scrollArea" :expand="true" :bottom-margin="20" style="width: 99%">
    <div
      v-for="(item, index) in devices"
      :key="index"
      class="row q-gutter-md q-pa-xs q-ma-xs"
      style="background-color: #f6f6f6"
    >
      <device :deviceid="item.id"></device>
    </div>
  </scroller>
</template>

<script setup>
import { onMounted, ref, computed } from "vue";
import { mapFields } from "vuex-map-fields";
import { useStore } from "vuex";
import scroller from "./Scroller.vue";

import { Device } from "./device/index";

const store = useStore();

const typeSelect = ref(-1);

const options = [
  {
    label: "All",
    value: -1,
  },
  {
    label: "Remotes",
    value: 0,
  },
  {
    label: "Lights",
    value: 2,
  },
  {
    label: "Plugs",
    value: 3,
  },
  {
    label: "Blinds",
    value: 7,
  },
  {
    label: "Signal Repeaters",
    value: 6,
  },
];

var devices = computed(() => {
  return store.getters.devices(typeSelect.value);
});

onMounted(() => {
  UpdateDevices();
});

const test = () => {
  type.value = 3;
};

const isDimmable = (device) => {
  return device.type == "dimmer";
};

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

<style>
.Scroll {
  background-color: bisque;
  color: v-bind("scrollColor");
  height: v-bind("scrollHeight");
  border: 1px solid black;
}
</style>
