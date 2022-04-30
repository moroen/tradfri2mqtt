<template>
  <q-toggle v-model="state" />
</template>

<script setup>
import { computed } from "vue";
import { useStore } from "vuex";

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
      class: "device",
      // command: "state",
      value: {
        deviceid: props.deviceid,
        payload: { State: value },
      }
    };
    store.dispatch("sendWSCommand", JSON.stringify(command));
  },
});
</script>
