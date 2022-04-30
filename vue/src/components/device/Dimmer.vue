<template>
  <div class="column">
    <div class="col self-end">{{ level }}</div>
    <div>
      <q-slider v-model="level" :min="0" :max="255" @change="saveValue" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from "vue";
import { useStore } from "vuex";
import sendCommand from "./common";

const store = useStore();

const props = defineProps({
  deviceid: {
    required: true,
  },
});

const saveValue = () => {
  const command = {
    class: "device",
    // command: "state",
    value: {
      deviceid: props.deviceid,
      payload: { brightness: store.state.websocket.devices[props.deviceid].dimmer },
    }
  };
  store.dispatch("sendWSCommand", JSON.stringify(command));
};

const level = computed({
  get() {
    return store.state.websocket.devices[props.deviceid].dimmer;
  },
  set(value) {
    var payload = {
      deviceid: props.deviceid,
      level: value,
    };
    store.commit("setDeviceDimmerLevel", payload);
  },
});

const standard = ref(0);
</script>
