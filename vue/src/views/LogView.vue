<template>
  <div class="q-pa-md" style="max-width: 300px">About</div>
  <q-page class="column" :style-fn="styleFn">
    <log-view></log-view>
  </q-page>
</template>

<script setup>
import { mapFields } from "vuex-map-fields";
import { useStore } from "vuex";
import LogView from "../components/LogViewer.vue";

const store = useStore();

const loc = window.location.hostname;
const port = window.location.port;

const styleFn = function (offset, height) {
  let pageheight = height - offset;
  console.log("PageHeight: " + pageheight);
  return "height: " + pageheight + "px";
};
</script>

<script>
export default {
  computed: {
    ...mapFields(["tradfri.gateway", "tradfri.ident"]),
  },
};
</script>
