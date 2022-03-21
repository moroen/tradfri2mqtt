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

const uri = `/api/devices/${props.deviceid}/set/state`;

const state = computed({
  get() {
    return store.state.websocket.devices[props.deviceid].state;
  },
  set(value) {
    const stateObject = {
      state: value,
    };

    const requestOptions = {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(stateObject),
    };

    fetch(uri, requestOptions)
      .then((response) => {
        if (response.status == 200) {
          return response.json();
        } else {
          throw new Error("Unable to set state");
        }
      })
      .then((data) => {
        // console.log("Got data");
        // console.log(data);
      })
      .catch((err) => {
        console.error("Unable to save setting\n" + err);
        showError("Unable to save settings");
      });

    console.log(value);
  },
});
</script>
