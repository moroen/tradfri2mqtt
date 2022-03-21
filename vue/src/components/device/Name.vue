<template>{{ name }}</template>

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
</script>
