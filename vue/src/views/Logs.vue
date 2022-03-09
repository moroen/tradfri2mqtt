<template>
  <q-page padding style="padding-top: 40px">
    <q-tab-panels v-model="tab" animated>
      <q-tab-panel name="log">
        <log-view></log-view>
      </q-tab-panel>

      <q-tab-panel name="wsLog">
        <ws-view></ws-view>
      </q-tab-panel>
    </q-tab-panels>
    <q-page-sticky expand position="top">
      <q-toolbar class="bg-white text-white">
        <q-tabs
          v-model="tab"
          dense
          class="text-grey"
          active-color="primary"
          indicator-color="primary"
          align="left"
          narrow-indicator
        >
          <q-tab name="log" label="Log" />
          <q-tab name="wsLog" label="Websocket" />
        </q-tabs>
      </q-toolbar>
    </q-page-sticky>
  </q-page>
</template>

<script setup>
import { mapFields } from "vuex-map-fields";
import { useStore } from "vuex";
import { ref } from "vue";

import logView from "../components/LogTradfri.vue";
import wsView from "../components/LogWebsocket.vue";

const store = useStore();

const tab = ref("log");

const loc = window.location.hostname;
const port = window.location.port;
</script>

<script>
export default {
  computed: {
    ...mapFields(["tradfri.gateway", "tradfri.ident"]),
  },
};
</script>
