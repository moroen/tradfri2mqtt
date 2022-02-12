<template>
  <div class="q-pa-md" style="max-width: 300px">About</div>

  <web-socket></web-socket>
</template>

<script setup>
import { mapFields } from "vuex-map-fields";
import { useStore } from "vuex";
import WebSocket from "../components/WebSocket.vue";

const store = useStore();
</script>

<script>
export default {
  computed: {
    ...mapFields(["tradfri.gateway", "tradfri.ident"]),
  },
};
</script>
