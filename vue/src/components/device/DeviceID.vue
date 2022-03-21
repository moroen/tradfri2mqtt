<template>{{ id }}</template>

<script setup>
import { computed } from "vue";
import { useStore } from "vuex";

const store = useStore();

const props = defineProps({
  deviceid: {
    required: true,
  },
});

const id = computed({
  get() {
    return store.state.websocket.devices[props.deviceid].id;
  },
});
</script>
