<template>
  <q-toggle v-model="state" />
</template>

<script setup>
import { computed } from "vue";
import { useStore } from "vuex";
import sendCommand from "./common";

const store = useStore();

const props = defineProps({
  deviceid: {
    required: true,
  },
});

const state = computed({
  get() {
    return store.state.websocket.devices[props.deviceid].state;
  },
  set(value) {
    const command = {
      command: "state",
      state: value,
    };

    sendCommand(props.deviceid, command);
  },
});
</script>
