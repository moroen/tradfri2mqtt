<template>
  <div class="column">
    <div>{{ name }}</div>
    <div class="text-caption text-weight-thin">{{ model }}</div>
  </div>
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

const name = computed({
  get() {
    return store.state.websocket.devices[props.deviceid].name;
  },
});

const model = computed({
  get() {
    return store.state.websocket.devices[props.deviceid].model;
  },
});
</script>
